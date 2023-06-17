package flag

import (
	"golang.org/x/exp/slog"

	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

type LogLevel struct {
	level slog.Level
}

func (x *LogLevel) Set(value string) error {
	levelMap := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	level, ok := levelMap[value]
	if !ok {
		return goerr.Wrap(types.ErrInvalidOption, "Invalid log level").With("level", value)
	}

	x.level = level
	return nil
}

func (x *LogLevel) String() string {
	return x.Level().String()
}

func (x *LogLevel) Level() slog.Level {
	if x.level == 0 {
		return slog.LevelInfo
	}

	return x.level
}
