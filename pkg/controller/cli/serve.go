package cli

import (
	"log/slog"

	"github.com/getsentry/sentry-go"
	"github.com/m-mizutani/alertchain/pkg/chain/core"
	"github.com/m-mizutani/alertchain/pkg/controller/cli/config"
	"github.com/m-mizutani/alertchain/pkg/controller/server"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/urfave/cli/v2"
)

func cmdServe() *cli.Command {
	var (
		addr          string
		disableAction bool

		dbCfg     config.Database
		policyCfg config.Policy
		sentryCfg config.Sentry
	)

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "addr",
			Usage:       "Bind address",
			Aliases:     []string{"a"},
			EnvVars:     []string{"ALERTCHAIN_ADDR"},
			Value:       "127.0.0.1:8080",
			Destination: &addr,
		},
		&cli.BoolFlag{
			Name:        "disable-action",
			Usage:       "Disable action execution (for debug or dry-run)",
			EnvVars:     []string{"ALERTCHAIN_DISABLE_ACTION"},
			Value:       false,
			Destination: &disableAction,
		},
	}
	flags = append(flags, dbCfg.Flags()...)
	flags = append(flags, policyCfg.Flags()...)
	flags = append(flags, sentryCfg.Flags()...)

	return &cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Flags:   flags,

		Action: func(c *cli.Context) error {
			utils.Logger().Info("starting alertchain with serve mode",
				slog.String("addr", addr),
				slog.Bool("disable-action", disableAction),
				slog.Any("database", dbCfg),
				slog.Any("sentry", sentryCfg),
			)

			var options []core.Option
			if disableAction {
				options = append(options, core.WithDisableAction())
			}

			ctx := model.NewContext(model.WithBase(c.Context))

			dbClient, dbCloser, err := dbCfg.New(ctx)
			if err != nil {
				return err
			}
			defer dbCloser()
			options = append(options, core.WithDatabase(dbClient))

			sentryCloser, err := sentryCfg.Configure()
			if err != nil {
				return err
			}
			defer sentryCloser()

			chain, err := buildChain(&policyCfg, options...)
			if err != nil {
				return err
			}

			authz, err := policyCfg.Load("authz")
			if err != nil {
				return err
			}

			utils.Logger().Info("starting alertchain with serve mode", slog.String("addr", addr))
			if err := server.New(chain.HandleAlert, authz).Run(addr); err != nil {
				sentry.CaptureException(err)
				return err
			}

			return nil
		},
	}
}
