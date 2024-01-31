package config

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/urfave/cli/v2"
)

type Sentry struct {
	dsn string
	env string
}

func (x *Sentry) Flags() []cli.Flag {
	category := "Sentry"

	return []cli.Flag{
		&cli.StringFlag{
			Name:        "sentry-dsn",
			Usage:       "Sentry DSN",
			Category:    category,
			Destination: &x.dsn,
			EnvVars:     []string{"ALERTCHAIN_SENTRY_DSN"},
		},
		&cli.StringFlag{
			Name:        "sentry-env",
			Usage:       "Sentry environment",
			Category:    category,
			Destination: &x.env,
			EnvVars:     []string{"ALERTCHAIN_SENTRY_ENV"},
		},
	}
}

func (x *Sentry) Configure() (func(), error) {
	if x.dsn == "" {
		utils.Logger().Warn("Sentry is not configured")
		return func() {}, nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:         x.dsn,
		Environment: x.env,
		Release:     fmt.Sprintf("alertchain@%s", types.AppVersion),
		Debug:       false,
	})
	if err != nil {
		utils.Logger().Warn("failed to initialize Sentry", utils.ErrLog(err))
		return nil, err
	}

	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	return func() {
		sentry.Recover()
		sentry.Flush(2 * time.Second)
	}, nil
}
