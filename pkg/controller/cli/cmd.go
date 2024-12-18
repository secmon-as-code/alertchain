package cli

import (
	"context"

	"github.com/secmon-lab/alertchain/pkg/controller/cli/config"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/logging"
	"github.com/urfave/cli/v3"
)

type CLI struct {
	app *cli.Command
}

func New() *CLI {
	var (
		loggerConfig config.Logger
	)
	flags := []cli.Flag{}
	flags = append(flags, loggerConfig.Flags()...)

	defers := []func(){}

	app := &cli.Command{
		Name:  "alertchain",
		Flags: flags,
		Before: func(ctx context.Context, _ *cli.Command) (context.Context, error) {
			closer, err := loggerConfig.Configure()
			if err != nil {
				return nil, err
			}
			defers = append(defers, closer)

			return ctx, nil
		},
		After: func(ctx context.Context, _ *cli.Command) error {
			for _, f := range defers {
				f()
			}
			return nil
		},

		Commands: []*cli.Command{
			cmdServe(),
			cmdRun(),
			cmdPlay(),
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Show version",
				Action: func(context.Context, *cli.Command) error {
					println("alertchain", types.AppVersion)
					return nil
				},
			},
		},
	}
	return &CLI{app: app}
}

func (x *CLI) Run(ctx context.Context, argv []string) error {
	if err := x.app.Run(ctx, argv); err != nil {
		ctxutil.Logger(ctx).Error("cli failed", logging.ErrAttr(err))
		return err
	}

	return nil
}
