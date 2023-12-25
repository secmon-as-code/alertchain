package utils

import (
	"io"

	"github.com/m-mizutani/goerr"
)

type Closer interface {
	Close() error
}

func SafeClose(c Closer) {
	if err := c.Close(); err != nil {
		Logger().Error("Fail to close io.WriteCloser", ErrLog(goerr.Wrap(err)))
	}
}

func SafeWrite(w io.Writer, b []byte) {
	if _, err := w.Write(b); err != nil {
		Logger().Error("Fail to write", ErrLog(goerr.Wrap(err)))
	}
}
