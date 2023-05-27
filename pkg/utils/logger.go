package utils

import (
	"reflect"
	"sync"
	"time"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/masq"
	"github.com/m-mizutani/slogger"
	"golang.org/x/exp/slog"
)

var logger = slog.Default()
var loggerMutex sync.Mutex

func Logger() *slog.Logger {
	return logger
}

func ReconfigureLogger(options ...slogger.Option) error {
	filter := masq.New(
		masq.WithTag("secret"),
		masq.WithFieldPrefix("secret_"),
		masq.WithAllowedType(reflect.TypeOf(time.Time{})),
	)

	options = append(options, slogger.WithReplacer(filter))

	newLogger, err := slogger.NewWithError(options...)
	if err != nil {
		return goerr.Wrap(err, "fail to initialize logger")
	}

	loggerMutex.Lock()
	logger = newLogger
	loggerMutex.Unlock()

	return nil
}

func ErrToAttrs(err error) []any {
	if err == nil {
		return nil
	}

	attrs := []any{
		slog.String("errmsg", err.Error()),
	}
	if e := goerr.Unwrap(err); e != nil {
		for k, v := range e.Values() {
			attrs = append(attrs, slog.Any("error."+k, v))
		}

		attrs = append(attrs, slog.Any("stacktrace", e.StackTrace()))
	}

	return attrs
}
