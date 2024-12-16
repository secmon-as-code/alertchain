package model

import (
	"context"
	"time"

	"log/slog"

	"github.com/secmon-lab/alertchain/pkg/utils"
)

type Clock func() time.Time

type Context struct {
	context.Context
	stack  int
	alert  *Alert
	dryRun bool
	clock  Clock
	onCLI  bool
}

type CtxOption func(ctx *Context)

func NewContext(options ...CtxOption) *Context {
	ctx := &Context{
		Context: context.Background(),
		clock:   time.Now,
	}

	for _, opt := range options {
		opt(ctx)
	}

	return ctx
}

func WithBase(base context.Context) CtxOption {
	return func(ctx *Context) {
		ctx.Context = base
	}
}

func WithAlert(alert Alert) CtxOption {
	return func(ctx *Context) {
		ctx.alert = &alert
	}
}

func WithStackIncrement() CtxOption {
	return func(ctx *Context) {
		ctx.stack++
	}
}

func WithDryRunMode() CtxOption {
	return func(ctx *Context) {
		ctx.dryRun = true
	}
}

func WithClock(clock Clock) CtxOption {
	return func(ctx *Context) {
		ctx.clock = clock
	}
}

func WithCLI() CtxOption {
	return func(ctx *Context) {
		ctx.onCLI = true
	}
}

func (x *Context) New(options ...CtxOption) *Context {
	ctx := x

	for _, opt := range options {
		opt(ctx)
	}

	return ctx
}

func (x *Context) Stack() int     { return x.stack }
func (x *Context) Alert() Alert   { return *x.alert }
func (x *Context) DryRun() bool   { return x.dryRun }
func (x *Context) Now() time.Time { return x.clock() }
func (x *Context) OnCLI() bool    { return x.onCLI }

func (x *Context) Logger() *slog.Logger {
	logger := utils.Logger()
	if x.alert != nil {
		logger = logger.With(slog.Any("alert_id", x.alert.ID))
	}
	return logger
}
