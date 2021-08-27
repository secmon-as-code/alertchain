package alertchain

import (
	"context"
	"io"
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

type ctxKey string

const ctxLogOutput ctxKey = "logOutput"

func InjectLogOutput(ctx context.Context, w io.Writer) context.Context {
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
