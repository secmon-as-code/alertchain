package utils

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/logging"
)

func HandleError(ctx context.Context, err error) {
	// Logging error
	ctxutil.Logger(ctx).Error("runtime error", logging.ErrAttr(err))

	// Sending error to Sentry
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		if goErr := goerr.Unwrap(err); goErr != nil {
			for k, v := range goErr.Values() {
				scope.SetExtra(fmt.Sprintf("%v", k), v)
			}
		}
	})
	hub.CaptureException(err)
}
