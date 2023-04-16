package model

import (
	"context"

	"github.com/m-mizutani/alertchain/pkg/utils"
	"golang.org/x/exp/slog"
)

type Context struct {
	context.Context
	stack int
	alert Alert
}

type CtxOption func(ctx *Context)

func NewContext(options ...CtxOption) *Context {
	ctx := &Context{
		Context: context.Background(),
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
		ctx.alert = alert
	}
}

func WithStackIncrement() CtxOption {
	return func(ctx *Context) {
		ctx.stack++
	}
}

func (x *Context) New(options ...CtxOption) *Context {
	ctx := x

	for _, opt := range options {
		opt(ctx)
	}

	return ctx
}

func (x *Context) Stack() int   { return x.stack }
func (x *Context) Alert() Alert { return x.alert }

func (x *Context) Logger() *slog.Logger {
	return utils.Logger().With(
		slog.Any("alert_id", x.alert.ID),
	)
}
