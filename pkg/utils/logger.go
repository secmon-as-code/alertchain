package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func initLogger() {
	Logger = zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func SetLogLevel(logLevel string) error {
	level := zerolog.InfoLevel
	switch strings.ToLower(logLevel) {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		return fmt.Errorf("invalid log level, choose one of debug, info, warn or error")
	}
	Logger = Logger.Level(level)
	return nil
}

func OutputError(logger zerolog.Logger, err error) {
	ev := logger.Error()

	var goErr *goerr.Error
	if errors.As(err, &goErr) {
		for k, v := range goErr.Values() {
			ev.Interface(k, v)
		}
		for _, st := range goErr.Stacks() {
			logger.Debug().Interface("stack", st).Send()
		}
	}
	ev.Msg(err.Error())
}
