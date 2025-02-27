package utils

import (
	"context"
	"io"

	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/logging"
)

type Closer interface {
	Close() error
}

func SafeClose(ctx context.Context, c Closer) {
	if err := c.Close(); err != nil {
		ctxutil.Logger(ctx).Error("Fail to close io.WriteCloser", logging.ErrAttr(goerr.Wrap(err, "Fail to close io.WriteCloser")))
	}
}

func SafeWrite(ctx context.Context, w io.Writer, b []byte) {
	if _, err := w.Write(b); err != nil {
		ctxutil.Logger(ctx).Error("Fail to write", logging.ErrAttr(goerr.Wrap(err, "Fail to write")))
	}
}
