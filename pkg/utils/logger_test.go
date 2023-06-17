package utils_test

import (
	"bytes"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/controller/cli/flag"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/gt"
	"golang.org/x/exp/slog"
)

func TestLogger(t *testing.T) {
	t.Run("default logger", func(t *testing.T) {
		var buf bytes.Buffer
		utils.ReconfigureLogger(&buf, slog.LevelInfo, flag.LogFormatJSON)
		utils.Logger().Info("hello",
			slog.String("secret_key", "xxx"),
			slog.String("normal_key", "aaa"),
		)

		gt.S(t, buf.String()).Contains("aaa").NotContains("xxx")
	})
}
