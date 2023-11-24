package utils

import (
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/m-mizutani/alertchain/pkg/controller/cli/flag"
	"github.com/m-mizutani/clog"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/masq"
)

var (
	logger      = slog.Default()
	loggerMutex sync.Mutex
	logFormat   flag.LogFormatType
)

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
			clog.WithSource(true),
			// clog.WithTimeFmt("2006-01-02 15:04:05"),
			clog.WithColorMap(&clog.ColorMap{
				Level: map[slog.Level]*color.Color{
					slog.LevelDebug: color.New(color.FgGreen, color.Bold),
					slog.LevelInfo:  color.New(color.FgCyan, color.Bold),
					slog.LevelWarn:  color.New(color.FgYellow, color.Bold),
					slog.LevelError: color.New(color.FgRed, color.Bold),
				},
				LevelDefault: color.New(color.FgBlue, color.Bold),
				Time:         color.New(color.FgWhite),
				Message:      color.New(color.FgHiWhite),
				AttrKey:      color.New(color.FgHiCyan),
				AttrValue:    color.New(color.FgHiWhite),
			}),
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
	logFormat = format
	loggerMutex.Unlock()
}

func ErrLog(err error) any {
	if err == nil {
		return nil
	}

	attrs := []any{
		slog.String("message", err.Error()),
	}

	if goErr := goerr.Unwrap(err); goErr != nil {
		var values []any
		for k, v := range goErr.Values() {
			values = append(values, slog.Any(k, v))
		}
		attrs = append(attrs, slog.Group("values", values...))

		var stacktrace any
		if logFormat == flag.LogFormatJSON {
			var traces []string
			for _, st := range goErr.StackTrace() {
				traces = append(traces, fmt.Sprintf("%+v", st))
			}
			stacktrace = traces
		} else {
			stacktrace = goErr.StackTrace()
		}

		attrs = append(attrs, slog.Any("stacktrace", stacktrace))
	}

	return slog.Group("error", attrs...)
}
