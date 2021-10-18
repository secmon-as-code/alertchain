package types

import (
	"context"
	"time"
)

type Context struct {
	base context.Context
}

func WrapContext(ctx context.Context) *Context {
	return &Context{base: ctx}
}

func NewContext() *Context {
	return &Context{base: context.Background()}
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

func (x *Context) SetTimeout(timeout time.Duration) context.CancelFunc {
	ctx, cancel := context.WithTimeout(x.base, timeout)
	x.base = ctx
	return cancel
}
