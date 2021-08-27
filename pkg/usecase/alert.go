package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

func (x *usecase) GetAlerts(ctx context.Context) ([]*ent.Alert, error) {
	return x.clients.DB.GetAlerts(ctx)
}

func (x *usecase) GetAlert(ctx context.Context, id types.AlertID) (*ent.Alert, error) {
	return x.clients.DB.GetAlert(ctx, id)
}

type ctxKey string

const (
	ctxKeyWaitGroup ctxKey = "WaitGroup"
)

func getWaitGroupFromCtx(ctx context.Context) *sync.WaitGroup {
	obj := ctx.Value(ctxKeyWaitGroup)
	if obj == nil {
		return nil
	}
	wg, ok := obj.(*sync.WaitGroup)
	if !ok {
		return nil
	}
	return wg
}

func ContextWithWaitGroup(ctx context.Context) (context.Context, *sync.WaitGroup) {
	wg := new(sync.WaitGroup)
	resp := context.WithValue(ctx, ctxKeyWaitGroup, wg)
	return resp, wg
}

func (x *usecase) RecvAlert(ctx context.Context, recvAlert *alertchain.Alert) (*alertchain.Alert, error) {
	newAlert, err := saveAlert(ctx, x.clients.DB, recvAlert)
	if err != nil {
		return nil, err
	}

	wg := getWaitGroupFromCtx(ctx)
	if wg != nil {
		wg.Add(1)
	}

	go func() {
		if wg != nil {
			defer wg.Done()
		}

		if err := executeChain(ctx, x.chain, newAlert.ID, x.clients); err != nil {
			utils.OutputError(logger, err)
		}
	}()

	return alertchain.NewAlert(newAlert, x.clients.DB), nil
}

func saveAlert(ctx context.Context, client infra.DBClient, recv *alertchain.Alert) (*ent.Alert, error) {
	if err := validateAlert(recv); err != nil {
		return nil, goerr.Wrap(err)
	}

	created, err := client.NewAlert(ctx)
	if err != nil {
		return nil, err
	}

	if err := client.UpdateAlert(ctx, created.ID, &recv.Alert); err != nil {
		return nil, err
	}

	attrs := make([]*ent.Attribute, len(recv.Attributes))
	for i, attr := range recv.Attributes {
		attrs[i] = &attr.Attribute
	}
	if err := client.AddAttributes(ctx, created.ID, attrs); err != nil {
		return nil, err
	}

	newAlert, err := client.GetAlert(ctx, created.ID)
	if err != nil {
		return nil, err
	}

	return newAlert, nil
}

func executeChain(ctx context.Context, chain *alertchain.Chain, alertID types.AlertID, clients infra.Clients) error {
	for idx, stage := range chain.Stages {
		if len(stage.Tasks) == 0 {
			continue
		}
		if stage.Timeout > 0 {
			newCtx, cancel := context.WithTimeout(ctx, stage.Timeout)
			defer cancel()
			ctx = newCtx
		}

		alert, err := clients.DB.GetAlert(ctx, alertID)
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		errCh := make(chan error, len(stage.Tasks))

		for _, task := range stage.Tasks {
			wg.Add(1)
			go executeTask(ctx, &executeTaskInput{
				stage:  int64(idx),
				task:   task,
				wg:     &wg,
				alert:  alert,
				client: clients.DB,
				errCh:  errCh,
			})
		}
		wg.Wait()

		close(errCh)
		for err := range errCh {
			if err != nil && stage.ExitOnErr {
				return err
			}
		}
	}
	return nil
}

type executeTaskInput struct {
	stage  int64
	task   alertchain.Task
	wg     *sync.WaitGroup
	alert  *ent.Alert
	client infra.DBClient
	errCh  chan error
}

type logWriter struct {
	bytes.Buffer
}

func (x *logWriter) Write(p []byte) (n int, err error) {
	if n, err := x.Buffer.Write(p); err != nil {
		return n, err
	}
	return os.Stdout.Write(p)
}

func executeTask(ctx context.Context, input *executeTaskInput) {
	defer input.wg.Done()
	var taskLog *ent.TaskLog
	var logW logWriter
	ctx = alertchain.InjectLogOutput(ctx, &logW)

	handlerError := func(err error) {
		wrapped := goerr.Wrap(err).With("task.Name", input.task.Name())
		utils.OutputError(logger, wrapped)
		input.errCh <- wrapped

		if taskLog != nil {
			taskLog.Errmsg = err.Error()
			var goErr *goerr.Error
			if errors.As(err, &goErr) {
				for k, v := range goErr.Values() {
					taskLog.ErrValues = append(taskLog.ErrValues, fmt.Sprintf("%s=%v", k, v))
				}
				for _, st := range goErr.StackTrace() {
					taskLog.StackTrace = append(taskLog.StackTrace, fmt.Sprintf("%v", st))
				}
			}
			taskLog.Status = types.TaskFailure
		}
	}

	defer func() {
		taskLog.Log = logW.String()
		taskLog.ExitedAt = time.Now().UTC().UnixNano()
		if err := input.client.UpdateTaskLog(ctx, taskLog); err != nil {
			input.errCh <- err
		}
	}()

	alert := alertchain.NewAlert(input.alert, input.client)

	now := time.Now().UTC().UnixNano()
	taskLog, err := input.client.NewTaskLog(ctx, input.alert.ID, input.task.Name(), now, input.stage, false)

	if err != nil {
		handlerError(err)
		return
	}

	if err := input.task.Execute(ctx, alert); err != nil {
		handlerError(err)
		return
	}

	if err := alert.Commit(ctx); err != nil {
		handlerError(err)
		return
	}

	taskLog.Status = types.TaskSucceeded
}

func validateAlert(alert *alertchain.Alert) error {
	if alert.Title == "" {
		return goerr.Wrap(types.ErrInvalidInput, "'title' field is required")
	}
	if alert.Detector == "" {
		return goerr.Wrap(types.ErrInvalidInput, "'detector' field is required")
	}

	for _, attr := range alert.Attributes {
		if attr.Key == "" {
			return goerr.Wrap(types.ErrInvalidInput, "'key' field is required").With("attr", attr)
		}
		if attr.Value == "" {
			return goerr.Wrap(types.ErrInvalidInput, "'value' field is required").With("attr", attr)
		}

		if err := attr.Type.IsValid(); err != nil {
			return goerr.Wrap(err).With("attr", attr)
		}

		for _, s := range attr.Context {
			ctx := types.AttrContext(s)
			if err := ctx.IsValid(); err != nil {
				return goerr.Wrap(err).With("attr", attr)
			}
		}
	}

	return nil
}
