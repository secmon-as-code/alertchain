package cli

import (
	"log/slog"

	"github.com/m-mizutani/alertchain/pkg/chain/core"
	"github.com/m-mizutani/alertchain/pkg/controller/cli/config"
	"github.com/m-mizutani/alertchain/pkg/controller/graphql"
	"github.com/m-mizutani/alertchain/pkg/controller/server"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/service"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/urfave/cli/v2"
)

func cmdServe() *cli.Command {
	var (
		addr          string
		disableAction bool
		playground    bool
		graphQL       bool

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
		&cli.BoolFlag{
			Name:        "graphql",
			Usage:       "Enable GraphQL",
			EnvVars:     []string{"ALERTCHAIN_GRAPHQL"},
			Value:       true,
			Destination: &graphQL,
		},
		&cli.BoolFlag{
			Name:        "playground",
			Usage:       "Enable GraphQL playground",
			EnvVars:     []string{"ALERTCHAIN_PLAYGROUND"},
			Value:       false,
			Destination: &playground,
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

			ctx := model.NewContext(model.WithBase(c.Context))

			// Build chain
			var chainOpt []core.Option
			if disableAction {
				chainOpt = append(chainOpt, core.WithDisableAction())
			}

			dbClient, dbCloser, err := dbCfg.New(ctx)
			if err != nil {
				return err
			}
			defer dbCloser()
			chainOpt = append(chainOpt, core.WithDatabase(dbClient))

			sentryCloser, err := sentryCfg.Configure()
			if err != nil {
				return err
			}
			defer sentryCloser()

			chain, err := buildChain(&policyCfg, chainOpt...)
			if err != nil {
				return err
			}

			// Build server
			var serverOpt []server.Option

			authz, err := policyCfg.Load("authz")
			if err != nil {
				return err
			}
			serverOpt = append(serverOpt, server.WithAuthzPolicy(authz))

			if graphQL {
				resolver := graphql.NewResolver(service.New(dbClient))
				serverOpt = append(serverOpt, server.WithResolver(resolver))
			}
			if playground {
				serverOpt = append(serverOpt, server.WithEnableGraphiQL())
			}

			srv := server.New(chain.HandleAlert, serverOpt...)

			// Starting server
			utils.Logger().Info("starting alertchain with serve mode", slog.String("addr", addr))
			if err := srv.Run(addr); err != nil {
				utils.HandleError(err)
				return err
			}

			return nil
		},
	}
}
