package cli

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"log/slog"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/alertchain/pkg/chain"
	"github.com/secmon-lab/alertchain/pkg/controller/cli/config"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/urfave/cli/v3"
)

func cmdRun() *cli.Command {
	var (
		input     string
		schema    types.Schema
		policyCfg config.Policy
	)

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "input",
			Aliases:     []string{"i"},
			Usage:       "input file or '-' for stdin",
			Sources:     cli.EnvVars("ALERTCHAIN_INPUT"),
			Required:    true,
			Destination: &input,
			Category:    "run",
		},
		&cli.StringFlag{
			Name:        "schema",
			Aliases:     []string{"s"},
			Usage:       "schema type",
			Sources:     cli.EnvVars("ALERTCHAIN_SCHEMA"),
			Required:    true,
			Destination: (*string)(&schema),
		},
	}
	flags = append(flags, policyCfg.Flags()...)

	return &cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "Run alertchain policy at once and exit in",
		Flags:   flags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var chainOptions []chain.Option

			chain, err := buildChain(&policyCfg, chainOptions...)
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

			ctxutil.Logger(ctx).Info("starting alertchain with run mode", slog.Any("data", data))

			if _, err := chain.HandleAlert(ctx, schema, data); err != nil {
				return goerr.Wrap(err, "failed to handle alert")
			}

			return nil
		},
	}
}
