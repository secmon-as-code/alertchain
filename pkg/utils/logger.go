package utils

import (
	"io"
	"reflect"
	"sync"
	"time"

	"github.com/m-mizutani/alertchain/pkg/controller/cli/flag"
	"github.com/m-mizutani/clog"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/masq"
	"golang.org/x/exp/slog"
)

var logger = slog.Default()
var loggerMutex sync.Mutex

func Logger() *slog.Logger {
	return logger
}

func ReconfigureLogger(w io.Writer, level slog.Level, format flag.LogFormatType) {
	filter := masq.New(
		masq.WithTag("secret"),
		masq.WithFieldPrefix("secret_"),
		masq.WithAllowedType(reflect.TypeOf(time.Time{})),
	)

	var handler slog.Handler
	switch format {
	case flag.LogFormatConsole:
		handler = clog.New(
			clog.WithWriter(w),
			clog.WithLevel(level),
			clog.WithReplaceAttr(filter),
		)

	case flag.LogFormatJSON:
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource:   true,
			Level:       level,
			ReplaceAttr: filter,
		})

	default:
		panic("Log format is not supported: " + format.String())
	}

	loggerMutex.Lock()
	logger = slog.New(handler)
	loggerMutex.Unlock()
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
