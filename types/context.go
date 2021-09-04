package types

import (
	"context"
	"io"
	"sync"
	"time"
)

type Context struct {
	base   context.Context
	writer io.Writer
	wg     *sync.WaitGroup
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

func (x *Context) copy() *Context {
	ctx := *x
	return &ctx
}

// Original functions
func (x *Context) Writer() io.Writer { return x.writer }
func (x *Context) InjectWriter(w io.Writer) *Context {
	ctx := x.copy()
	ctx.writer = w
	return ctx
}

func (x *Context) WaitGroup() *sync.WaitGroup { return x.wg }
func (x *Context) InjectWaitGroup(wg *sync.WaitGroup) *Context {
	ctx := x.copy()
	ctx.wg = wg
	return ctx
}
