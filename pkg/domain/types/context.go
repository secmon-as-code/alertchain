package types

import "context"

type Context struct {
	context.Context
	stack int
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

func (x *Context) Stack() int { return x.stack }
