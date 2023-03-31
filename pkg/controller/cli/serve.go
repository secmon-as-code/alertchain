package cli

import (
	"github.com/m-mizutani/alertchain/pkg/controller/server"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/urfave/cli/v2"
)

func cmdServe(cfg *model.Config) *cli.Command {
	var (
		addr string
	)
	return &cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "addr",
				Usage:       "Bind address",
				EnvVars:     []string{"ALERTCHAIN_ADDR"},
				Value:       "127.0.0.1:3000",
				Destination: &addr,
			},
		},

		Action: func(ctx *cli.Context) error {
			chain, err := buildChain(*cfg)
			if err != nil {
				return err
			}

			if err := server.New(chain.HandleAlert).Run(addr); err != nil {
				return err
			}

			return nil
		},
	}
}
