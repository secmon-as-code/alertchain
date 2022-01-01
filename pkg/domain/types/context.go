package types

import (
	"context"
	"time"

	"github.com/m-mizutani/zlog"
)

type Context struct {
	logger  *zlog.Logger
	base    context.Context
	timeout context.Context
}

func NewContextWith(ctx context.Context, logger *zlog.Logger) *Context {
	return &Context{
		base:    ctx,
		timeout: ctx,
	}
}

func NewContext(logger *zlog.Logger) *Context {
	return NewContextWith(context.Background(), logger)
}

// context.Context
func (x *Context) Deadline() (deadline time.Time, ok bool) {
	return x.timeout.Deadline()
}
func (x *Context) Done() <-chan struct{} {
	return x.timeout.Done()
}
func (x *Context) Err() error {
	return x.timeout.Err()
}
func (x *Context) Value(key interface{}) interface{} {
	return x.base.Value(key)
}

func (x *Context) SetTimeout(timeout time.Duration) context.CancelFunc {
	if timeout == 0 {
		x.timeout = x.base
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	x.timeout = ctx
	return cancel
}

func (x *Context) SetLogger(logger *zlog.Logger) {
	x.logger = logger
}

func (x *Context) Logger() *zlog.Logger {
	if x.logger == nil {
		x.logger = zlog.New()
	}
	return x.logger
}
