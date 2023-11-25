package utils

import "io"

func SafeClose(c io.WriteCloser) {
	if err := c.Close(); err != nil {
		Logger().Error("Fail to close io.WriteCloser", ErrLog(err))
	}
}
