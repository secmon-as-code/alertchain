package utils_test

import (
	"bytes"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/slogger"
	"golang.org/x/exp/slog"
)

func TestLogger(t *testing.T) {
	t.Run("default logger", func(t *testing.T) {
		defer gt.NoError(t, utils.ReconfigureLogger())
		var buf bytes.Buffer
		gt.NoError(t, utils.ReconfigureLogger(slogger.WithWriter(&buf)))
		utils.Logger().Info("hello",
			slog.String("secret_key", "xxx"),
			slog.String("normal_key", "aaa"),
		)

		gt.S(t, buf.String()).Contains("aaa").NotContains("xxx")
	})
}
