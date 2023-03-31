package cli

import (
	"encoding/json"
	"os"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
)

func cmdConfig(cfg *model.Config) *cli.Command {
	return &cli.Command{
		Name:    "config",
		Aliases: []string{"c"},
		Subcommands: []*cli.Command{
			cmdConfigShow(cfg),
			cmdConfigValidate(cfg),
		},
	}
}

func cmdConfigShow(cfg *model.Config) *cli.Command {
	return &cli.Command{
		Name:    "show",
		Aliases: []string{"s"},
		Action: func(ctx *cli.Context) error {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(cfg); err != nil {
				return goerr.Wrap(err, "Fail to encode config")
			}

			return nil
		},
	}
}

func cmdConfigValidate(cfg *model.Config) *cli.Command {
	return &cli.Command{
		Name:    "validate",
		Aliases: []string{"v"},
		Action: func(ctx *cli.Context) error {
			if _, err := buildChain(*cfg); err != nil {
				return err
			}

			return nil
		},
	}
}
