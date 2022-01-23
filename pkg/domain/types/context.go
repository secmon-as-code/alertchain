package types

import (
	"context"
	"time"

	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/zlog"
)

type Context struct {
	logger *zlog.LogEntity
	base   context.Context
}

type ContextOption func(ctx *Context)

func NewContext(options ...ContextOption) *Context {
	base := context.Background()
	ctx := &Context{
		logger: utils.Logger.Log(),
		base:   base,
	}
	for _, opt := range options {
		opt(ctx)
	}
	return ctx
}

func WithLogger(logger *zlog.LogEntity) ContextOption {
	return func(ctx *Context) {
		ctx.logger = logger
	}
}

func WithCtx(ctx context.Context) ContextOption {
	return func(ctx *Context) {
		ctx.base = ctx
	}
}

// context.Context
func (x *Context) Deadline() (deadline time.Time, ok bool) {
	return x.base.Deadline()
}
func (x *Context) Done() <-chan struct{} {
	return x.base.Done()
}
func (x *Context) Err() error {
	return x.base.Err()
}
func (x *Context) Value(key interface{}) interface{} {
	return x.base.Value(key)
}

func (x *Context) WithTimeout(timeout time.Duration) (*Context, context.CancelFunc) {
	ctxTimeout, cancel := context.WithTimeout(x.base, timeout)

	return &Context{
		base:   ctxTimeout,
		logger: x.logger,
	}, cancel
}

func (x *Context) Log() *zlog.LogEntity {
	return x.logger
}
