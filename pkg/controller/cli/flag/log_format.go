package flag

import (
	"strings"

	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

type LogFormatType int

func (x LogFormatType) String() string {
	switch x {
	case LogFormatConsole:
		return "console"
	case LogFormatJSON:
		return "json"
	default:
		return "unknown"
	}
}

const (
	LogFormatConsole LogFormatType = iota + 1
	LogFormatJSON
)

type LogFormat struct {
	fmtType LogFormatType
}

func NewLogFormat(fmtType LogFormatType) *LogFormat {
	return &LogFormat{fmtType: fmtType}
}

func (x *LogFormat) Set(value string) error {
	formatMap := map[string]LogFormatType{
		"console": LogFormatConsole,
		"c":       LogFormatConsole,
		"json":    LogFormatJSON,
		"j":       LogFormatJSON,
	}

	fmtType, ok := formatMap[strings.ToLower(value)]
	if !ok {
		return goerr.Wrap(types.ErrInvalidOption, "Invalid log format").With("format", value)
	}

	x.fmtType = fmtType
	return nil
}

func (x *LogFormat) String() string {
	return x.Format().String()
}

func (x *LogFormat) Format() LogFormatType {
	if x.fmtType == 0 {
		return LogFormatConsole
	}
	return x.fmtType
}
