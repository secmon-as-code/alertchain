package flag

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mizutani/goerr"
)

type LogOutput struct {
	v string
	w io.Writer
}

func (x *LogOutput) Set(value string) error {
	var w io.Writer
	switch strings.ToLower(value) {
	case "-", "stdout":
		w = os.Stdout
	case "stderr":
		w = os.Stderr
	default:
		f, err := os.Create(filepath.Clean(value))
		if err != nil {
			return goerr.Wrap(err, "Failed to open log file").With("path", value)
		}
		w = f
	}

	x.v = value
	x.w = w

	return nil
}

func (x *LogOutput) String() string {
	return x.v
}

func (x *LogOutput) Writer() io.Writer {
	if x.w == nil {
		return os.Stdout
	}
	return x.w
}
