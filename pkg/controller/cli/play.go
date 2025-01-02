package cli

import (
	"context"

	"github.com/secmon-lab/alertchain/pkg/controller/cli/config"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/usecase"
	"github.com/urfave/cli/v3"
)

func cmdPlay() *cli.Command {
	var (
		input usecase.PlayInput

		policyCfg config.Policy
	)

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "scenario",
			Aliases:     []string{"s"},
			Usage:       "scenario directory",
			Sources:     cli.EnvVars("ALERTCHAIN_SCENARIO"),
			Required:    true,
			Destination: &input.ScenarioPath,
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       "output directory",
			Sources:     cli.EnvVars("ALERTCHAIN_OUTPUT"),
			Destination: &input.OutDir,
			Value:       "./output",
		},
		&cli.StringSliceFlag{
			Name:        "target",
			Aliases:     []string{"t"},
			Usage:       "Target scenario ID to play. If not specified, all scenarios are played",
			Sources:     cli.EnvVars("ALERTCHAIN_TARGET"),
			Destination: &input.Targets,
		},
	}
	flags = append(flags, policyCfg.Flags()...)

	return &cli.Command{
		Name:    "play",
		Aliases: []string{"p"},
		Usage:   "Simulate alertchain policy",
		Flags:   flags,

		Action: func(ctx context.Context, cmd *cli.Command) error {
			ctx = ctxutil.SetCLI(ctx)

			coreOptions, err := policyCfg.CoreOption(ctx)
			if err != nil {
				return err
			}
			input.CoreOptions = coreOptions

			if err := usecase.Play(ctx, input); err != nil {
				return err
			}

			return nil
		},
	}
}
