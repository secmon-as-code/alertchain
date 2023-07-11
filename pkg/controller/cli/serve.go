package cli

import (
	"github.com/getsentry/sentry-go"
	"github.com/m-mizutani/alertchain/pkg/chain"
	"github.com/m-mizutani/alertchain/pkg/controller/server"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/firestore"
	"github.com/m-mizutani/alertchain/pkg/infra/memory"
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
		dbType        string

		firestoreProjectID  string
		firestoreCollection string
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

			&cli.StringFlag{
				Name:        "db-type",
				Usage:       "Database type (memory, firestore)",
				Aliases:     []string{"t"},
				EnvVars:     []string{"ALERTCHAIN_DB_TYPE"},
				Value:       "memory",
				Destination: &dbType,
			},
			&cli.StringFlag{
				Name:        "firestore-project-id",
				Usage:       "Project ID of Firestore",
				Category:    "firestore",
				EnvVars:     []string{"ALERTCHAIN_FIRESTORE_PROJECT_ID"},
				Destination: &firestoreProjectID,
			},
			&cli.StringFlag{
				Name:        "firestore-collection",
				Usage:       "Collection name of Firestore",
				Category:    "firestore",
				EnvVars:     []string{"ALERTCHAIN_FIRESTORE_COLLECTION"},
				Destination: &firestoreCollection,
			},
		}, cfg.Flags()...),

		Action: func(c *cli.Context) error {
			utils.Logger().Info("starting alertchain with serve mode",
				slog.String("addr", addr),
				slog.Bool("disable-action", disableAction),
				slog.Bool("enable-sentry", enableSentry),
				slog.String("db-type", dbType),
				slog.String("firestore-project-id", firestoreProjectID),
				slog.String("firestore-collection", firestoreCollection),
				slog.Any("config", cfg),
			)

			var options []chain.Option
			if disableAction {
				options = append(options, chain.WithDisableAction())
			}

			ctx := model.NewContext(model.WithBase(c.Context))

			switch dbType {
			case "memory":
				options = append(options, chain.WithDatabase(memory.New()))

			case "firestore":
				if firestoreProjectID == "" {
					return goerr.Wrap(types.ErrInvalidOption, "firestore-project-id is required for firestore")
				}
				if firestoreCollection == "" {
					return goerr.Wrap(types.ErrInvalidOption, "firestore-collection is required for firestore")
				}

				client, err := firestore.New(ctx, firestoreProjectID, firestoreCollection)
				if err != nil {
					return goerr.Wrap(err, "failed to initialize firestore client")
				}
				defer func() {
					if err := client.Close(); err != nil {
						ctx.Logger().Error("Failed to close firestore client", utils.ErrLog(err))
					}
				}()

				options = append(options, chain.WithDatabase(client))

			default:
				return goerr.Wrap(types.ErrInvalidOption, "invalid db-type").With("db-type", dbType)
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
