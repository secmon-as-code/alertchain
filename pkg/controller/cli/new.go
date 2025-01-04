package cli

import (
	"context"

	"github.com/secmon-lab/alertchain/pkg/usecase"
	"github.com/urfave/cli/v3"
)

func cmdNew() *cli.Command {
	var (
		dir string
	)

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "dir",
			Aliases:     []string{"d"},
			Usage:       "Directory path to create new AlertChain policy repository",
			Sources:     cli.EnvVars("ALERTCHAIN_DIR"),
			Destination: &dir,
			Value:       ".",
		},
	}

	return &cli.Command{
		Name:  "new",
		Usage: "Create new AlertChain policy repository",
		Flags: flags,

		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := usecase.NewPolicyDirectory(ctx, dir); err != nil {
				return err
			}

			return nil
		},
	}
}
