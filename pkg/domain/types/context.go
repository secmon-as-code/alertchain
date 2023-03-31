package types

import "context"

type Context struct {
	context.Context
}

type CtxOption func(ctx Context)

func NewContext(options ...CtxOption) *Context {
	ctx := &Context{
		Context: context.Background(),
	}

	return ctx
}

func WithBase(base context.Context) CtxOption {
	return func(ctx Context) {
		ctx.Context = base
	}
}
