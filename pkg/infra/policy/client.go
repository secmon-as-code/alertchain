package policy

import "context"

type Client interface {
	Eval(ctx context.Context, in interface{}, out interface{}) error
}
