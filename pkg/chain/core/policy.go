package core

import (
	"errors"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/goerr"
	"golang.org/x/exp/slog"
)

func (x *Core) QueryAlertPolicy(ctx *model.Context, schema types.Schema, in, out any) error {
	if x.alertPolicy == nil {
		return nil
	}

	options := []policy.QueryOption{
		policy.WithPackageSuffix(string(schema)),
	}
	if x.enablePrint {
		options = append(options, policy.WithRegoPrint(makeRegoPrint(ctx)))
	}

	if err := x.alertPolicy.Query(ctx, in, out, options...); err != nil && !errors.Is(err, types.ErrNoPolicyResult) {
		return goerr.Wrap(err, "failed to evaluate alert policy").With("request", in)
	}
	ctx.Logger().Info("queried action policy", slog.Any("in", in), slog.Any("out", out))

	return nil
}

func (x *Core) QueryActionPolicy(ctx *model.Context, in, out any) error {
	if x.actionPolicy == nil {
		return nil
	}

	var options []policy.QueryOption
	if x.enablePrint {
		options = append(options, policy.WithRegoPrint(makeRegoPrint(ctx)))
	}

	if err := x.actionPolicy.Query(ctx, in, out, options...); err != nil && !errors.Is(err, types.ErrNoPolicyResult) {
		return goerr.Wrap(err, "failed to evaluate action policy").With("request", in)
	}
	ctx.Logger().Info("queried action policy", slog.Any("in", in), slog.Any("out", out))

	return nil
}

func makeRegoPrint(ctx *model.Context) policy.RegoPrint {
	return func(file string, row int, msg string) error {
		ctx.Logger().Info("rego print",
			slog.String("file", file),
			slog.Int("row", row),
			slog.String("msg", msg),
		)
		return nil
	}
}
