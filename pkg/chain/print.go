package chain

import (
	"io"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"golang.org/x/exp/slog"
)

type printHook struct {
	ctx *model.Context
}

func newPrintHook(ctx *model.Context) io.Writer {
	return &printHook{ctx: ctx}
}

func (x *printHook) Write(p []byte) (int, error) {
	x.ctx.Logger().Info("rego print message", slog.String("rego.msg", string(p)))
	return len(p), nil
}
