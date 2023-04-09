package cli

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
)

func cmdRun(cfg *model.Config) *cli.Command {
	var (
		input  string
		schema types.Schema
	)

	return &cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "input",
				Aliases:     []string{"i"},
				Usage:       "input file or '-' for stdin",
				EnvVars:     []string{"ALERTCHAIN_INPUT"},
				Required:    true,
				Destination: &input,
			},
			&cli.StringFlag{
				Name:        "schema",
				Aliases:     []string{"s"},
				Usage:       "schema type",
				EnvVars:     []string{"ALERTCHAIN_SCHEMA"},
				Required:    true,
				Destination: (*string)(&schema),
			},
		},

		Action: func(c *cli.Context) error {
			chain, err := buildChain(*cfg)
			if err != nil {
				return err
			}

			var r io.Reader
			if input != "-" {
				fd, err := os.Open(filepath.Clean(input))
				if err != nil {
					return goerr.Wrap(err, "failed to open input file")
				}
				r = fd
			} else {
				r = os.Stdin
			}

			var data any
			if err := json.NewDecoder(r).Decode(&data); err != nil {
				return goerr.Wrap(err, "failed to decode input data")
			}

			ctx := types.NewContext(types.WithBase(c.Context))
			if err := chain.HandleAlert(ctx, schema, data); err != nil {
				return goerr.Wrap(err, "failed to handle alert")
			}

			return nil
		},
	}
}
