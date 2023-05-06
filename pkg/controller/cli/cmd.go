package cli

import (
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

type CLI struct {
	app *cli.App
}

func New() *CLI {
	var (
		logLevel  string
		logFormat string
		logOutput string

		cfg model.Config
	)

	app := &cli.App{
		Name: "alertchain",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"l"},
				Category:    "logging",
				Usage:       "Set log level [debug|info|warn|error]",
				Value:       "info",
				Destination: &logLevel,
			},
			&cli.StringFlag{
				Name:        "log-format",
				Aliases:     []string{"f"},
				Category:    "logging",
				Usage:       "Set log format [text|json]",
				Value:       "text",
				Destination: &logFormat,
			},
			&cli.StringFlag{
				Name:        "log-output",
				Aliases:     []string{"o"},
				Category:    "logging",
				Usage:       "Set log output (create file other than '-', 'stdout', 'stderr')",
				Value:       "-",
				Destination: &logOutput,
			},
			&cli.BoolFlag{
				Name:        "enable-print",
				Aliases:     []string{"p"},
				Category:    "logging",
				EnvVars:     []string{"ALERTCHAIN_ENABLE_PRINT"},
				Usage:       "Enable print feature in Rego. The cli option is priority than config file.",
				Value:       false,
				Destination: &cfg.Policy.Print,
			},
			&cli.StringFlag{
				Name:        "policy-dir",
				Aliases:     []string{"d"},
				Usage:       "directory path of policy files",
				EnvVars:     []string{"ALERTCHAIN_POLICY_DIR"},
				Required:    true,
				Destination: &cfg.Policy.Path,
			},
		},

		Before: func(ctx *cli.Context) error {
			if err := utils.ReconfigureLogger(logFormat, logLevel, logOutput); err != nil {
				return err
			}

			utils.Logger().Debug("config loaded", slog.Any("config", cfg))

			return nil
		},

		Commands: []*cli.Command{
			cmdServe(&cfg),
			cmdRun(&cfg),
			cmdPlay(&cfg),
		},
	}
	return &CLI{app: app}
}

func (x *CLI) Run(argv []string) error {
	if err := x.app.Run(argv); err != nil {
		attrs := []any{
			slog.String("error", err.Error()),
		}

		if goErr := goerr.Unwrap(err); goErr != nil {
			for k, v := range goErr.Values() {
				attrs = append(attrs, slog.Any(k, v))
			}
		}

		utils.Logger().Error("cli failed", attrs...)
		return err
	}

	return nil
}
