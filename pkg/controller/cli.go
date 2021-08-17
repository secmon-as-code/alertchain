package controller

import (
	"os"
	"plugin"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/controller/server"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/usecase"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
	cli "github.com/urfave/cli/v2"
)

type config struct {
	logLevel string

	// serve
	DBType     string
	DBConfig   string
	ServerAddr string
	ServerPort uint64

	ChainPath string
}

var logger = utils.Logger

func (x *Controller) CLI(args []string) {
	var cfg config
	app := cli.App{
		Name:        "alertchain",
		Version:     types.Version,
		Description: "Programmable SOAR (Security Orchestration, Automation and Response) platform and universal alert manager",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"l"},
				Usage:       "[debug|info|warn|error]",
				Value:       "info",
				Destination: &cfg.logLevel,
			},
		},
		Before: func(c *cli.Context) error {
			if err := utils.SetLogLevel(cfg.logLevel); err != nil {
				return err
			}
			return nil
		},
		Commands: []*cli.Command{
			cmdServe(&cfg),
		},
	}

	if err := app.Run(args); err != nil {
		utils.OutputError(logger, err)
		os.Exit(1)
	}
}

func cmdServe(cfg *config) *cli.Command {
	return &cli.Command{
		Name: "serve",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "db-type",
				Aliases:     []string{"t"},
				Usage:       "database type",
				Destination: &cfg.DBType,
				Value:       "sqlite3",
				EnvVars:     []string{"ALERTCHAIN_DB_TYPE"},
			},
			&cli.StringFlag{
				Name:        "db-config",
				Aliases:     []string{"d"},
				Usage:       "database configuration",
				Destination: &cfg.DBConfig,
				Value:       "file:ent?mode=memory&cache=shared&_fk=1",
				EnvVars:     []string{"ALERTCHAIN_DB_CONFIG"},
			},
			&cli.StringFlag{
				Name:        "server-addr",
				Aliases:     []string{"a"},
				Usage:       "server binding adddress",
				Value:       "localhost",
				Destination: &cfg.ServerAddr,
			},
			&cli.Uint64Flag{
				Name:        "server-port",
				Aliases:     []string{"p"},
				Usage:       "server binding port",
				Value:       9080,
				Destination: &cfg.ServerPort,
			},

			&cli.StringFlag{
				Name:        "chain",
				Aliases:     []string{"c"},
				Usage:       "chain plugin path",
				Required:    true,
				Destination: &cfg.ChainPath,
			},
		},

		Action: func(c *cli.Context) error {
			logger.Info().Interface("config", cfg).Msg("Starting AlertChain")

			// Setup database
			dbClient, err := db.New(cfg.DBType, cfg.DBConfig)
			if err != nil {
				return err
			}
			defer func() {
				if err := dbClient.Close(); err != nil {
					logger.Err(err).Msg("Failed to close database conn")
				}
			}()

			// Setup chain
			chain, err := loadChainPlugin(cfg.ChainPath)
			if err != nil {
				return err
			}

			// Setup usecase
			uc := usecase.New(infra.Clients{DB: dbClient}, chain)

			// Starting server
			if err := server.New(uc, cfg.ServerAddr, cfg.ServerPort).Run(); err != nil {
				return err
			}

			return nil
		},
	}
}

func loadChainPlugin(filePath string) (*alertchain.Chain, error) {
	p, err := plugin.Open(filePath)
	if err != nil {
		return nil, types.ErrInvalidChain.Wrap(err)
	}

	f, err := p.Lookup("Chain")
	if err != nil {
		return nil, goerr.Wrap(types.ErrInvalidChain, "Chain() function not found")
	}

	getChain, ok := f.(func() *alertchain.Chain)
	if !ok {
		return nil, goerr.Wrap(types.ErrInvalidChain, "Chain() type mismatch")
	}

	return getChain(), nil
}
