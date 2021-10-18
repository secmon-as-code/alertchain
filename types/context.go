package types

import (
	"context"
	"time"
)

type Context struct {
	base    context.Context
	timeout context.Context
}

func NewContextWith(ctx context.Context) *Context {
	return &Context{
		base:    ctx,
		timeout: ctx,
	}
}

func NewContext() *Context {
	return NewContextWith(context.Background())
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
