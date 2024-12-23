package core

import (
	"context"
	"errors"
	"fmt"

	"log/slog"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/infra/policy"
)

func (x *Core) QueryAlertPolicy(ctx context.Context, schema types.Schema, in, out any) error {
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
		return types.AsPolicyErr(goerr.Wrap(err, "failed to evaluate alert policy").With("request", in))
	}
	ctxutil.Logger(ctx).Info("queried action policy", slog.Any("in", in), slog.Any("out", out))

	return nil
}

func (x *Core) QueryActionPolicy(ctx context.Context, in, out any) error {
	if x.actionPolicy == nil {
		return nil
	}

	var options []policy.QueryOption
	if x.enablePrint {
		options = append(options, policy.WithRegoPrint(makeRegoPrint(ctx)))
	}

	if err := x.actionPolicy.Query(ctx, in, out, options...); err != nil && !errors.Is(err, types.ErrNoPolicyResult) {
		return types.AsPolicyErr(goerr.Wrap(err, "failed to evaluate action policy").With("request", in))
	}
	ctxutil.Logger(ctx).Info("queried action policy", slog.Any("in", in), slog.Any("out", out))

	return nil
}

func makeRegoPrint(ctx context.Context) policy.RegoPrint {
	return func(file string, row int, msg string) error {
		if ctxutil.IsCLI(ctx) {
			fmt.Printf("	%s:%d: %s\n", file, row, msg)
		} else {
			ctxutil.Logger(ctx).Info("rego print",
				slog.String("file", file),
				slog.Int("row", row),
				slog.String("msg", msg),
			)
		}
		return nil
	}
}
