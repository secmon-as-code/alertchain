package alertchain

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/usecase"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/alertchain/types"
)

var logger = utils.Logger

type Chain struct {
	Jobs    Jobs
	Sources []Source
	Actions Actions
}

type Job struct {
	Timeout   time.Duration
	ExitOnErr bool
	Tasks     []Task
}

type Task interface {
	Name() string
	Execute(ctx *types.Context, alert *Alert) error
}

type Source interface {
	Name() string
	Run(alertCh chan *Alert) error
}

type Action interface {
	Name() string
	Executable(attr *Attribute) bool
	Execute(ctx *types.Context, attr *Attribute) error
}

type Jobs []*Job
type Actions []Action

func (x Jobs) Convert() []*usecase.Job {
	resp := make([]*usecase.Job, len(x))
	for i, j := range x {
		job := &usecase.Job{
			Timeout:   j.Timeout,
			ExitOnErr: j.ExitOnErr,
			Tasks:     make([]*usecase.Task, len(j.Tasks)),
		}

		for i, task := range j.Tasks {
			job.Tasks[i] = &usecase.Task{
				Name: task.Name(),
				Execute: func(ctx *types.Context, alert *ent.Alert) (*usecase.ChangeRequest, error) {
					newAlert := NewAlert(alert)
					if err := task.Execute(ctx, newAlert); err != nil {
						return nil, err
					}
					return &newAlert.ChangeRequest, nil
				},
			}
		}

		resp[i] = job
	}
	return resp
}

func (x Actions) Convert() []*usecase.Action {
	resp := make([]*usecase.Action, len(x))
	for i, action := range x {
		resp[i] = &usecase.Action{
			ID:   uuid.New().String(),
			Name: action.Name(),
			Executable: func(attr *ent.Attribute) bool {
				return action.Executable(newAttribute(attr))
			},
			Execute: func(ctx *types.Context, attr *ent.Attribute) error {
				return action.Execute(ctx, newAttribute(attr))
			},
		}
	}
	return resp
}

// TestInvokeTasks runs
func (x *Chain) TestInvokeTasks(t *testing.T, recv *Alert) (*Alert, error) {
	var wg sync.WaitGroup
	ctx := types.NewContext().InjectWaitGroup(&wg)

	clients := &infra.Clients{
		DB: db.NewDBMock(t),
	}

	alert, err := x.InvokeTasks(ctx, recv, clients)
	if err != nil {
		return nil, err
	}
	wg.Wait()

	created, err := clients.DB.GetAlert(types.NewContext(), alert.id)
	if err != nil {
		return nil, err
	}

	return NewAlert(created), nil
}

func (x *Chain) InvokeTasks(ctx *types.Context, recv *Alert, clients *infra.Clients) (*Alert, error) {
	uc := usecase.New(*clients, x.Jobs.Convert(), x.Actions.Convert())
	created, err := uc.HandleAlert(ctx, &recv.Alert, recv.Attributes.toEnt())
	if err != nil {
		return nil, err
	}

	return NewAlert(created), nil
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

func runSource(src Source, ch chan *Alert) {
	for {
		if err := src.Run(ch); err != nil {
			utils.HandleError(err)
		}
		time.Sleep(time.Second * 3)
	}
}

func (x *Chain) ActivateSources() chan *Alert {
	alertCh := make(chan *Alert, 256)
	for _, src := range x.Sources {
		go runSource(src, alertCh)
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
