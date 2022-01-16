package utils

import (
	"os"

	"github.com/m-mizutani/zlog"
)

var Logger = zlog.New()

func init() {
	if logLevel, ok := os.LookupEnv("LOG_LEVEL"); ok {
		Logger = zlog.New(zlog.WithLogLevel(logLevel))
	}
}
