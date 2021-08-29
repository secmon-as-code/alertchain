package alertchain

import (
	"bytes"
	"context"
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
	Executable(attr *Attribute) bool
	Execute(ctx context.Context, attr *Attribute) error
}

type ActionEntry struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Action Action `json:"-"`
}

type ActionLog struct {
	ent.ActionLog
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
			utils.HandleError(err)
		}
	}()

	return NewAlert(newAlert, clients.DB), nil
}

func saveAlert(ctx context.Context, dbClient db.Interface, recv *Alert) (*ent.Alert, error) {
	if err := recv.Validate(); err != nil {
		return nil, goerr.Wrap(err)
	}

	created, err := dbClient.NewAlert(ctx)
	if err != nil {
		return nil, err
	}

	if err := dbClient.UpdateAlert(ctx, created.ID, &recv.Alert); err != nil {
		return nil, err
	}

	attrs := make([]*ent.Attribute, len(recv.Attributes))
	for i, attr := range recv.Attributes {
		attrs[i] = &attr.Attribute
	}
	if err := dbClient.AddAttributes(ctx, created.ID, attrs); err != nil {
		return nil, err
	}

	newAlert, err := dbClient.GetAlert(ctx, created.ID)
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
	client db.Interface
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
	var execLog ent.ExecLog
	var logW logWriter
	ctx = utils.InjectLogWriter(ctx, &logW)

	handlerError := func(err error) {
		wrapped := goerr.Wrap(err).With("task.Name", input.task.Name())
		utils.HandleError(wrapped)
		input.errCh <- wrapped

		if taskLog != nil {
			utils.CopyErrorToExecLog(err, &execLog)
		}
	}

	defer func() {
		if taskLog != nil {
			execLog.Log = logW.String()
			execLog.Timestamp = time.Now().UTC().UnixNano()

			if err := input.client.AppendTaskLog(ctx, taskLog.ID, &execLog); err != nil {
				input.errCh <- err
			}
		}
	}()

	alert := NewAlert(input.alert, input.client)

	taskLog, err := input.client.NewTaskLog(ctx, input.alert.ID, input.task.Name(), input.stage)
	if err != nil {
		handlerError(err)
		return
	}

	if err := input.client.AppendTaskLog(ctx, taskLog.ID, &ent.ExecLog{
		Status: types.ExecStart,
	}); err != nil {
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

// LogWriter returns io.Writer to output log message. Log message output via io.Writer will be stored into TaskLog and displayed in Web UI.
func LogWriter(ctx context.Context) io.Writer {
	return utils.LogWriter(ctx)
}
