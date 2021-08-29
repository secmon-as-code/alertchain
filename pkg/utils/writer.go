package utils

import (
	"context"
	"io"
)

type ctxKey string

const ctxLogWriter ctxKey = "LogWriter"

func InjectLogWriter(ctx context.Context, w io.Writer) context.Context {
	return context.WithValue(ctx, ctxLogWriter, w)
}

// LogWriter returns io.Writer to output log message. Log message output via io.Writer will be stored into TaskLog and displayed in Web UI.
func LogWriter(ctx context.Context) io.Writer {
	value := ctx.Value(ctxLogWriter)
	if value == nil {
		panic("LogWriter is not set in context")
	}
	w, ok := value.(io.Writer)
	if !ok {
		panic("LogWriter is not io.Writer")
	}
	return w
}
