package types

import (
	"context"
	"time"

	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/zlog"
)

type Context struct {
	logger *zlog.LogEntity
	parent context.Context
}

type ContextOption func(ctx *Context)

func NewContext(options ...ContextOption) *Context {
	parent := context.Background()
	ctx := &Context{
		logger: utils.Logger.Log(),
		parent: parent,
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

func WithCtx(parent context.Context) ContextOption {
	return func(ctx *Context) {
		ctx.parent = parent
	}
}

// context.Context
func (x *Context) Deadline() (deadline time.Time, ok bool) {
	return x.parent.Deadline()
}
func (x *Context) Done() <-chan struct{} {
	return x.parent.Done()
}
func (x *Context) Err() error {
	return x.parent.Err()
}
func (x *Context) Value(key interface{}) interface{} {
	return x.parent.Value(key)
}

func (x *Context) WithTimeout(timeout time.Duration) (*Context, context.CancelFunc) {
	ctxTimeout, cancel := context.WithTimeout(x.parent, timeout)

	return &Context{
		parent: ctxTimeout,
		logger: x.logger,
	}, cancel
}

func (x *Context) Log() *zlog.LogEntity {
	return x.logger
}
