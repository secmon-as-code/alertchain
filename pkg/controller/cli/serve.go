package cli

import (
	"context"
	"log/slog"

	"github.com/secmon-lab/alertchain/pkg/chain"
	"github.com/secmon-lab/alertchain/pkg/controller/cli/config"
	"github.com/secmon-lab/alertchain/pkg/controller/graphql"
	"github.com/secmon-lab/alertchain/pkg/controller/server"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/service"
	"github.com/secmon-lab/alertchain/pkg/utils"
	"github.com/urfave/cli/v3"
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
			Sources:     cli.EnvVars("ALERTCHAIN_ADDR"),
			Value:       "127.0.0.1:8080",
			Destination: &addr,
		},
		&cli.BoolFlag{
			Name:        "graphql",
			Usage:       "Enable GraphQL",
			Sources:     cli.EnvVars("ALERTCHAIN_GRAPHQL"),
			Value:       true,
			Destination: &graphQL,
		},
		&cli.BoolFlag{
			Name:        "playground",
			Usage:       "Enable GraphQL playground",
			Sources:     cli.EnvVars("ALERTCHAIN_PLAYGROUND"),
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

		Action: func(ctx context.Context, cmd *cli.Command) error {
			ctxutil.Logger(ctx).Info("starting alertchain with serve mode",
				slog.String("addr", addr),
				slog.Bool("disable-action", disableAction),
				slog.Any("database", dbCfg),
				slog.Any("sentry", sentryCfg),
			)

			// Build chain
			var chainOpt []chain.Option

			dbClient, dbCloser, err := dbCfg.New(ctx)
			if err != nil {
				return err
			}
			defer dbCloser()
			chainOpt = append(chainOpt, chain.WithDatabase(dbClient))

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
			ctxutil.Logger(ctx).Info("starting alertchain with serve mode", slog.String("addr", addr))
			if err := srv.Run(addr); err != nil {
				utils.HandleError(ctx, err)
				return err
			}

			return nil
		},
	}
}
