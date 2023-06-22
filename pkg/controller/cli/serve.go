package cli

import (
	"github.com/getsentry/sentry-go"
	"github.com/m-mizutani/alertchain/pkg/chain"
	"github.com/m-mizutani/alertchain/pkg/controller/server"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

func cmdServe(cfg *model.Config) *cli.Command {
	var (
		addr          string
		disableAction bool
		enableSentry  bool
	)
	return &cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Flags: append([]cli.Flag{
			&cli.StringFlag{
				Name:        "addr",
				Usage:       "Bind address",
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
			&cli.BoolFlag{
				Name:        "enable-sentry",
				Usage:       "Enable sentry logging, you need to set SENTRY_DSN environment variable",
				EnvVars:     []string{"ALERTCHAIN_ENABLE_SENTRY"},
				Destination: &enableSentry,
			},
		}, cfg.Flags()...),

		Action: func(ctx *cli.Context) error {
			utils.Logger().Info("starting alertchain with serve mode",
				slog.String("addr", addr),
				slog.Bool("disable-action", disableAction),
				slog.Bool("enable-sentry", enableSentry),
				slog.Any("config", cfg),
			)

			var options []chain.Option
			if disableAction {
				options = append(options, chain.WithDisableAction())
			}

			if enableSentry {
				if err := sentry.Init(sentry.ClientOptions{}); err != nil {
					return goerr.Wrap(err, "Failed to initialize sentry")
				}
			}

			chain, err := buildChain(*cfg, options...)
			if err != nil {
				return err
			}

			utils.Logger().Info("starting alertchain with serve mode", slog.String("addr", addr))
			if err := server.New(chain.HandleAlert).Run(addr); err != nil {
				sentry.CaptureException(err)
				return err
			}

			return nil
		},
	}
}
