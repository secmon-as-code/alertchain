package utils_test

import (
	"bytes"
	"testing"

	"log/slog"

	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/gt"
)

func TestLogger(t *testing.T) {
	t.Run("default logger", func(t *testing.T) {
		var buf bytes.Buffer
		utils.ReconfigureLogger(&buf, slog.LevelInfo, utils.LogFormatJSON)
		utils.Logger().Info("hello",
			slog.String("secret_key", "xxx"),
			slog.String("normal_key", "aaa"),
		)

		gt.S(t, buf.String()).Contains("aaa").NotContains("xxx")
	})
}
