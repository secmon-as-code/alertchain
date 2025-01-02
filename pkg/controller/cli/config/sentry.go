package config

import (
	"context"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/logging"
	"github.com/urfave/cli/v3"
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
			Sources:     cli.EnvVars("ALERTCHAIN_SENTRY_DSN"),
		},
		&cli.StringFlag{
			Name:        "sentry-env",
			Usage:       "Sentry environment",
			Category:    category,
			Destination: &x.env,
			Sources:     cli.EnvVars("ALERTCHAIN_SENTRY_ENV"),
		},
	}
}

func (x *Sentry) Configure(ctx context.Context) (func(), error) {
	if x.dsn == "" {
		ctxutil.Logger(ctx).Warn("Sentry is not configured")
		return func() {}, nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:         x.dsn,
		Environment: x.env,
		Release:     fmt.Sprintf("alertchain@%s", types.AppVersion),
		Debug:       false,
	})
	if err != nil {
		ctxutil.Logger(ctx).Warn("failed to initialize Sentry", logging.ErrAttr(err))
		return nil, err
	}

	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	return func() {
		sentry.Recover()
		sentry.Flush(2 * time.Second)
	}, nil
}
