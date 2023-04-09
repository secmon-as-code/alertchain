package utils

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/m-mizutani/goerr"
	"golang.org/x/exp/slog"
)

var logger = slog.Default()
var loggerMutex sync.Mutex

func Logger() *slog.Logger {
	return logger
}

func ReconfigureLogger(format, level, output string) error {
	logLevelMap := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	logLevel, ok := logLevelMap[level]
	if !ok {
		return goerr.New("invalid log format, should be 'debug', 'info', 'warn' or 'error'").With("actual", level)
	}

	var w io.Writer
	switch output {
	case "-", "stdout":
		w = os.Stdout
	case "stderr":
		w = os.Stderr
	default:
		fd, err := os.Create(filepath.Clean(output))
		if err != nil {
			return goerr.Wrap(err, "opening log file")
		}
		w = fd
	}

	opt := slog.HandlerOptions{
		AddSource: logLevel <= slog.LevelDebug,
		Level:     logLevel,
	}

	var newLogger *slog.Logger
	switch format {
	case "text":
		newLogger = slog.New(opt.NewTextHandler(w))
	case "json":
		newLogger = slog.New(opt.NewJSONHandler(w))
	default:
		return goerr.New("invalid log format, should be 'text' or 'json'").With("actual", format)
	}

	loggerMutex.Lock()
	logger = newLogger
	loggerMutex.Unlock()

	return nil
}
