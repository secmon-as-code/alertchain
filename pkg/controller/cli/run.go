package cli

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/m-mizutani/alertchain/pkg/chain"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

func cmdRun(cfg *model.Config) *cli.Command {
	var (
		input       string
		schema      types.Schema
		enablePrint bool
	)

	return &cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "Run alertchain policy at once and exit in",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "input",
				Aliases:     []string{"i"},
				Usage:       "input file or '-' for stdin",
				EnvVars:     []string{"ALERTCHAIN_INPUT"},
				Required:    true,
				Destination: &input,
				Category:    "run",
			},
			&cli.StringFlag{
				Name:        "schema",
				Aliases:     []string{"s"},
				Usage:       "schema type",
				EnvVars:     []string{"ALERTCHAIN_SCHEMA"},
				Required:    true,
				Destination: (*string)(&schema),
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

			chain, err := buildChain(*cfg, chainOptions...)
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

			ctx := model.NewContext(model.WithBase(c.Context))
			ctx.Logger().Info("starting alertchain with run mode", slog.Any("data", data))

			if err := chain.HandleAlert(ctx, schema, data); err != nil {
				return goerr.Wrap(err, "failed to handle alert")
			}

			return nil
		},
	}
}
