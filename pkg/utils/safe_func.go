package utils

import (
	"io"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/alertchain/pkg/logging"
)

type Closer interface {
	Close() error
}

func SafeClose(c Closer) {
	if err := c.Close(); err != nil {
		logging.Default().Error("Fail to close io.WriteCloser", logging.ErrAttr(goerr.Wrap(err)))
	}
}

func SafeWrite(w io.Writer, b []byte) {
	if _, err := w.Write(b); err != nil {
		logging.Default().Error("Fail to write", logging.ErrAttr(goerr.Wrap(err)))
	}
}
