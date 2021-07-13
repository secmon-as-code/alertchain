package cli

import (
	"os"

	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/usecase"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/urfave/cli/v2"
)

var logger = utils.Logger

type cliConfig struct {
	logLevel   string
	newUsecase interfaces.NewUsecase
}

type CLI struct {
	config cliConfig
}

func New() *CLI {
	return &CLI{
		config: cliConfig{
			newUsecase: usecase.New,
		},
	}
}

func (x *CLI) Run(argv []string) error {
	app := &cli.App{
		Name:        "alertchain",
		Description: "Alert handling service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"l"},
				EnvVars:     []string{"ALERTCHAIN_LOG_LEVEL"},
				Usage:       "LogLevel [trace|debug|info|warn|error]",
				Destination: &x.config.logLevel,
			},
		},
		Commands: []*cli.Command{
			serveCommand(&x.config),
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Error().Err(err).Msg("Failed")
		return err
	}

	return nil
}
