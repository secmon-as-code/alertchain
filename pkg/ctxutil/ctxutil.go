package ctxutil

import (
	"context"
	"log/slog"
	"time"

	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/logging"
)

type ctxStackKey struct{}

func InjectStack(ctx context.Context, stack int) context.Context {
	return context.WithValue(ctx, ctxStackKey{}, stack)
}

func GetStack(ctx context.Context) int {
	v := ctx.Value(ctxStackKey{})
	if v == nil {
		return 0
	}
	return v.(int)
}

type ctxAlertKey struct{}

func InjectAlert(ctx context.Context, alert *model.Alert) context.Context {
	return context.WithValue(ctx, ctxAlertKey{}, alert)
}

func GetAlert(ctx context.Context) *model.Alert {
	v := ctx.Value(ctxAlertKey{})
	if v == nil {
		return nil
	}
	return v.(*model.Alert)
}

type ctxDryRunKey struct{}

func SetDryRun(ctx context.Context, dryRun bool) context.Context {
	return context.WithValue(ctx, ctxDryRunKey{}, dryRun)
}

func IsDryRun(ctx context.Context) bool {
	v := ctx.Value(ctxDryRunKey{})
	if v == nil {
		return false
	}
	return v.(bool)
}

type ctxClockKey struct{}

func InjectClock(ctx context.Context, clock model.Clock) context.Context {
	return context.WithValue(ctx, ctxClockKey{}, clock)
}

func Now(ctx context.Context) time.Time {
	v := ctx.Value(ctxClockKey{})
	if v == nil {
		return time.Now()
	}
	return v.(model.Clock)()
}

type ctxCLIKey struct{}

func SetCLI(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxCLIKey{}, true)
}

func IsCLI(ctx context.Context) bool {
	v := ctx.Value(ctxCLIKey{})
	if v == nil {
		return false
	}
	return v.(bool)
}

type ctxLoggerKey struct{}

func InjectLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, logger)
}

func Logger(ctx context.Context) *slog.Logger {
	v := ctx.Value(ctxLoggerKey{})
	if v == nil {
		return logging.Default()
	}
	return v.(*slog.Logger)
}
