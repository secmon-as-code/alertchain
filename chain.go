package alertchain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

var logger = utils.Logger

type Chain struct {
	Jobs    []*Job
	Sources []Source
	Actions []Action
}

type Job struct {
	Timeout   time.Duration
	ExitOnErr bool
	Tasks     []Task
}

type Task interface {
	Name() string
	Execute(ctx context.Context, alert *Alert) error
}

type Source interface {
	Name() string
	Run(alertCh chan *Alert) error
}

type Action interface {
	Name() string
	Actionable(attr *Attribute) bool
	Act(ctx context.Context, attr *Attribute) error
}

// TestInvokeTasks runs
func (x *Chain) TestInvokeTasks(t *testing.T, recv *Alert) (*Alert, error) {
	ctx, wg := setWaitGroupToCtx(context.Background())

	clients := &infra.Clients{
		DB: db.NewDBMock(t),
	}

	alert, err := x.InvokeTasks(ctx, recv, clients)
	if err != nil {
		return nil, err
	}
	wg.Wait()

	created, err := clients.DB.GetAlert(context.Background(), alert.id)
	if err != nil {
		return nil, err
	}

	return NewAlert(created, nil), nil
}

func (x *Chain) LookupTask(taskType interface{}) Task {
	actual := reflect.TypeOf(taskType)
	for _, job := range x.Jobs {
		for _, task := range job.Tasks {
			if reflect.TypeOf(task) == actual {
				return task
			}
		}
	}

	return nil
}

type ctxKey string

const (
	ctxKeyWaitGroup ctxKey = "WaitGroup"
	ctxLogOutput    ctxKey = "logOutput"
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

func setWaitGroupToCtx(ctx context.Context) (context.Context, *sync.WaitGroup) {
	wg := new(sync.WaitGroup)
	resp := context.WithValue(ctx, ctxKeyWaitGroup, wg)
	return resp, wg
}

func (x *Chain) InvokeTasks(ctx context.Context, recv *Alert, clients *infra.Clients) (*Alert, error) {
	newAlert, err := saveAlert(ctx, clients.DB, recv)
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

		if err := x.executeJobs(ctx, newAlert.ID, clients); err != nil {
			utils.OutputError(logger, err)
		}
	}()

	return NewAlert(newAlert, clients.DB), nil
}

func saveAlert(ctx context.Context, client infra.DBClient, recv *Alert) (*ent.Alert, error) {
	if err := recv.Validate(); err != nil {
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

func (x *Chain) executeJobs(ctx context.Context, alertID types.AlertID, clients *infra.Clients) error {
	for idx, job := range x.Jobs {
		if len(job.Tasks) == 0 {
			continue
		}
		if job.Timeout > 0 {
			newCtx, cancel := context.WithTimeout(ctx, job.Timeout)
			defer cancel()
			ctx = newCtx
		}

		alert, err := clients.DB.GetAlert(ctx, alertID)
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		errCh := make(chan error, len(job.Tasks))

		for _, task := range job.Tasks {
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
			if err != nil && job.ExitOnErr {
				return err
			}
		}
	}
	return nil
}

type executeTaskInput struct {
	stage  int64
	task   Task
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
	ctx = injectLogOutput(ctx, &logW)

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

	alert := NewAlert(input.alert, input.client)

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

func (x *Chain) ActivateSources() chan *Alert {
	alertCh := make(chan *Alert, 256)
	for _, src := range x.Sources {
		go src.Run(alertCh)
	}
	return alertCh
}

func (x *Chain) NewJob() *Job {
	job := &Job{}
	x.Jobs = append(x.Jobs, job)
	return job
}

func (x *Job) AddTask(task Task) {
	x.Tasks = append(x.Tasks, task)
}

func injectLogOutput(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, ctxLogOutput, w)
}

// LogOutput returns io.Writer to output log message. Log message output via io.Writer will be stored into TaskLog and displayed in Web UI.
func LogOutput(ctx context.Context) io.Writer {
	value := ctx.Value(ctxLogOutput)
	if value == nil {
		panic("logOutput is not set in context")
	}
	w, ok := value.(io.Writer)
	if !ok {
		panic("logOutput is not io.Writer")
	}
	return w
}
