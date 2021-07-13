package cli

import (
	"fmt"

	"github.com/m-mizutani/alertchain/pkg/controller/api"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/urfave/cli/v2"
)

type serveConfig struct {
	Addr string
	Port int
}

func serveCommand(cliCfg *cliConfig) *cli.Command {
	var cfg serveConfig
	return &cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "Addr",
				Usage:       "server binding address",
				Aliases:     []string{"a"},
				EnvVars:     []string{"ALERTCHAIN_ADDR"},
				Destination: &cfg.Addr,
				Value:       "127.0.0.1",
			},
			&cli.IntFlag{
				Name:        "Port",
				Usage:       "Port number",
				Aliases:     []string{"p"},
				EnvVars:     []string{"ALERTCHAIN_PORT"},
				Destination: &cfg.Port,
				Value:       9080,
			},
		},
		Action: func(c *cli.Context) error {
			serverAddr := fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)

			uc := cliCfg.newUsecase(&model.Config{})
			engine := api.New(uc)
			if err := engine.Run(serverAddr); err != nil {
				logger.Error().Err(err).Interface("config", cfg).Msg("Server error")
			}
			return nil
		},
	}
}
