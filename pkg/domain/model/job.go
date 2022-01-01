package model

import (
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/types"
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
