package utils

type Closer interface {
	Close() error
}

func SafeClose(c Closer) {
	if err := c.Close(); err != nil {
		Logger().Error("Fail to close io.WriteCloser", ErrLog(err))
	}
}
