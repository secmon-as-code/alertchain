package utils

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/zlog"
)

var Logger *zlog.Logger

func initLogger() {
	Logger = zlog.New()
}

func SetLogLevel(logLevel string) error {
	return Logger.SetLogLevel(logLevel)
}

func SetLogFormatter(logfmt string) error {
	switch logfmt {
	case "console":
		Logger.Formatter = zlog.NewConsoleFormatter()
	case "json":
		Logger.Formatter = zlog.NewJsonFormatter()
	default:
		return goerr.New("Unsupported log format").With("format", logfmt)
	}

	return nil
}
