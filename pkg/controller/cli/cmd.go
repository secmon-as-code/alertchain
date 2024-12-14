package cli

import (
	"context"

	"github.com/m-mizutani/alertchain/pkg/controller/cli/config"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/urfave/cli/v2"
)

type CLI struct {
	app *cli.App
}

func New() *CLI {
	var (
		loggerConfig config.Logger
	)
	flags := []cli.Flag{}
	flags = append(flags, loggerConfig.Flags()...)

	defers := []func(){}

	app := &cli.App{
		Name:  "alertchain",
		Flags: flags,
		Before: func(ctx *cli.Context) error {
			closer, err := loggerConfig.Configure()
			if err != nil {
				return err
			}
			defers = append(defers, closer)

			return nil
		},
		After: func(ctx *cli.Context) error {
			for _, f := range defers {
				f()
			}
			return nil
		},

		Commands: []*cli.Command{
			cmdServe(),
			cmdRun(),
			cmdPlay(),
			cmdNew(),
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Show version",
				Action: func(c *cli.Context) error {
					println("alertchain", types.AppVersion)
					return nil
				},
			},
		},
	}
	return &CLI{app: app}
}

func (x *CLI) Run(ctx context.Context, argv []string) error {
	if err := x.app.RunContext(ctx, argv); err != nil {
		utils.Logger().Error("cli failed", utils.ErrLog(err))
		return err
	}

	return nil
}
