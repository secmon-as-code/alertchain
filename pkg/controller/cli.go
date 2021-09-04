package controller

import (
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
	chain    *alertchain.Chain

	// serve
	DBType     string
	DBConfig   string
	ServerAddr string
	ServerPort uint64

	ChainPath string
}

var logger = utils.Logger

func (x *Controller) CLIWithChain(args []string, chain *alertchain.Chain) error {
	cfg := config{
		chain: chain,
	}

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
		utils.HandleError(err)
		return err
	}
	return nil
}

func (x *Controller) CLI(args []string) error {
	return x.CLIWithChain(args, nil)
}

func cmdServe(cfg *config) *cli.Command {
	flags := []cli.Flag{
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
	}

	if cfg.chain == nil {
		flags = append(flags, &cli.StringFlag{
			Name:        "chain",
			Aliases:     []string{"c"},
			Usage:       "chain plugin path",
			Required:    true,
			Destination: &cfg.ChainPath,
		})
	}

	return &cli.Command{
		Name:  "serve",
		Flags: flags,

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
			chain := cfg.chain
			if chain == nil {
				c, err := loadChainPlugin(cfg.ChainPath)
				if err != nil {
					return err
				}
				chain = c
			}

			// Setup usecase
			uc := usecase.New(infra.Clients{DB: dbClient}, chain.Jobs.Convert(), chain.Actions.Convert())

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

	makeChain, ok := f.(func() (*alertchain.Chain, error))
	if !ok {
		return nil, goerr.Wrap(types.ErrInvalidChain, "Chain() type mismatch")
	}

	chain, err := makeChain()
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return chain, nil
}
