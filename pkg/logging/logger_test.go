package logging_test

import (
	"bytes"
	"testing"

	"log/slog"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/alertchain/pkg/logging"
)

func TestLogger(t *testing.T) {
	t.Run("default logger", func(t *testing.T) {
		var buf bytes.Buffer
		logging.ReconfigureLogger(&buf, slog.LevelInfo, logging.FormatJSON)
		logging.Default().Info("hello",
			slog.String("secret_key", "xxx"),
			slog.String("normal_key", "aaa"),
		)

		gt.S(t, buf.String()).Contains("aaa").NotContains("xxx")
	})
}
