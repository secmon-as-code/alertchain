package policy

import "github.com/open-policy-agent/opa/v1/topdown/print"

type regoPrintHook struct {
	callback RegoPrint
}

func (x *regoPrintHook) Print(ctx print.Context, msg string) error {
	return x.callback(ctx.Location.File, ctx.Location.Row, msg)
}
