package cli

import (
	"os"
	"path/filepath"
	"strings"

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

		configFile string
		configData string

		cfg model.Config
	)
	app := &cli.App{
		Name: "alertchain",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Category:    "logging",
				Usage:       "Set log level [debug|info|warn|error]",
				Value:       "info",
				Destination: &logLevel,
			},
			&cli.StringFlag{
				Name:        "log-format",
				Category:    "logging",
				Usage:       "Set log format [text|json]",
				Value:       "text",
				Destination: &logFormat,
			},
			&cli.StringFlag{
				Name:        "log-output",
				Category:    "logging",
				Usage:       "Set log output (create file other than '-', 'stdout', 'stderr')",
				Value:       "-",
				Destination: &logOutput,
			},

			&cli.StringFlag{
				Name:        "config-file",
				Aliases:     []string{"c"},
				Category:    "config",
				EnvVars:     []string{"ALERTCHAIN_CONFIG_FILE"},
				Usage:       "Set config jsonnet file path",
				Destination: &configFile,
			},
			&cli.StringFlag{
				Name:        "config-data",
				Aliases:     []string{"d"},
				Category:    "config",
				EnvVars:     []string{"ALERTCHAIN_CONFIG_DATA"},
				Usage:       "Set config jsonnet data content",
				Destination: &configData,
			},
		},

		Before: func(ctx *cli.Context) error {
			if err := utils.ReconfigureLogger(logFormat, logLevel, logOutput); err != nil {
				return err
			}

			data := configData
			if configFile != "" {
				raw, err := os.ReadFile(filepath.Clean(configFile))
				if err != nil {
					return goerr.Wrap(err, "reading config file")
				}
				data = string(raw)
			}

			var envVars []model.EnvVar
			for _, env := range os.Environ() {
				keyValue := strings.SplitN(env, "=", 2)
				envVars = append(envVars, model.EnvVar{
					Key:   keyValue[0],
					Value: keyValue[1],
				})
			}

			if err := model.ParseConfig(configFile, data, envVars, &cfg); err != nil {
				return err
			}

			utils.Logger().Debug("config loaded", slog.Any("config", cfg))

			return nil
		},

		Commands: []*cli.Command{
			cmdConfig(&cfg),
			cmdServe(&cfg),
			cmdRun(&cfg),
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
