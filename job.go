package alertchain

import (
	"sync"
	"time"

	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

type Jobs []*Job
type Job struct {
	Timeout   time.Duration
	ExitOnErr bool
	Tasks     []Task
}

type Task interface {
	Name() string
	Execute(ctx *types.Context, alert *Alert) error
}

func (x Jobs) Execute(ctx *types.Context, client db.Interface, alertID types.AlertID) error {
	for idx, job := range x {
		ctx.Logger().With("job", job).With("step", idx).Trace("Starting Job")
		if err := job.Execute(ctx, client, alertID); err != nil {
			return err
		}
		ctx.Logger().With("job", job).Trace("Exiting Job")
	}
	return nil
}

func (x *Job) AddTask(task Task) {
	x.Tasks = append(x.Tasks, task)
}

func (x *Job) Execute(ctx *types.Context, client db.Interface, alertID types.AlertID) error {
	if len(x.Tasks) == 0 {
		return nil
	}

	cancel := ctx.SetTimeout(x.Timeout)
	if cancel != nil {
		defer cancel()
	}

	base, err := client.GetAlert(ctx, alertID)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(x.Tasks))

	for _, task := range x.Tasks {
		wg.Add(1)
		go executeTask(ctx, &executeTaskInput{
			task:   task,
			wg:     &wg,
			alert:  newAlert(base),
			client: client,
			errCh:  errCh,
		})
	}
	wg.Wait()
	ctx.Logger().With("job", x).Trace("Completed Job")

	close(errCh)
	for err := range errCh {
		ctx.Logger().Err(err).With("job", x).Error("failed job")
		if err != nil && x.ExitOnErr {
			return err
		}
	}

	return nil
}

type executeTaskInput struct {
	task   Task
	wg     *sync.WaitGroup
	alert  *Alert
	client db.Interface
	errCh  chan error
}

func executeTask(ctx *types.Context, input *executeTaskInput) {
	defer input.wg.Done()

	handlerError := func(err error) {
		input.errCh <- goerr.Wrap(err).With("task.Name", input.task.Name())
	}

	if err := input.task.Execute(ctx, input.alert); err != nil {
		handlerError(err)
		return
	}

	if err := input.alert.commit(ctx, input.client, input.alert.base.ID); err != nil {
		handlerError(err)
		return
	}
}
