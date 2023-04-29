package cli

import (
	"os"

	"github.com/m-mizutani/alertchain/pkg/chain"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
)

func cmdPlay(cfg *model.Config) *cli.Command {
	var (
		playbookPath string
		enablePrint  bool
	)

	return &cli.Command{
		Name:    "play",
		Aliases: []string{"p"},
		Usage:   "Simulate alertchain policy",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "playbook",
				Aliases:     []string{"b"},
				Usage:       "playbook file",
				EnvVars:     []string{"ALERTCHAIN_PLAYBOOK"},
				Required:    true,
				Destination: &playbookPath,
			},
			&cli.BoolFlag{
				Name:        "enable-print",
				Aliases:     []string{"p"},
				Usage:       "enable print feature in Rego",
				EnvVars:     []string{"ALERTCHAIN_ENABLE_PRINT"},
				Destination: &enablePrint,
			},
		},

		Action: func(c *cli.Context) error {
			var chainOptions []chain.Option
			if enablePrint {
				chainOptions = append(chainOptions, chain.WithEnablePrint())
			}

			// Load playbook
			var playbook model.Playbook
			if err := model.ParsePlaybook(playbookPath, os.ReadFile, &playbook); err != nil {
				return goerr.Wrap(err, "failed to parse playbook")
			}

			ctx := model.NewContext(model.WithBase(c.Context))
			for _, s := range playbook.Scenarios {
				chain, err := buildChain(*cfg, chainOptions...)
				if err != nil {
					return err
				}

				if err := chain.HandleAlert(ctx, s.Schema, s.Alert); err != nil {
					return goerr.Wrap(err, "failed to handle alert")
				}
			}

			return nil
		},
	}
}
