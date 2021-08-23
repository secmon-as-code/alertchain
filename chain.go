package alertchain

import (
	"context"
	"time"
)

type Chain struct {
	Stages   []*Stage
	Optional []Task
	Sources  []Source
}

type Stage struct {
	Timeout   time.Duration
	ExitOnErr bool
	Tasks     []Task
}

func (x *Chain) NewStage() *Stage {
	stage := &Stage{}
	x.Stages = append(x.Stages, stage)
	return stage
}

func (x *Stage) AddTask(task Task) {
	x.Tasks = append(x.Tasks, task)
}

type Task interface {
	Name() string
	Description() string
	Execute(ctx context.Context, alert *Alert) error
	Optionable(alert *Alert) bool
}

type Source interface {
	Name() string
	Run(alertCh chan *Alert) error
}
