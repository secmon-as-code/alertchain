package utils

import (
	"github.com/getsentry/sentry-go"
	"github.com/m-mizutani/goerr"
)

func HandleError(err error) {
	// Logging error
	Logger().Error("runtime error", ErrLog(err))

	// Sending error to Sentry
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		if goErr := goerr.Unwrap(err); goErr != nil {
			for k, v := range goErr.Values() {
				scope.SetExtra(k, v)
			}
		}
	})
	hub.CaptureException(err)
}
