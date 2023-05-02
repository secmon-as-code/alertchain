package logger_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/infra/logger"
	"github.com/m-mizutani/gt"
)

// bufferWriteCloser is a wrapper around bytes.Buffer that implements io.WriteCloser.
type bufferWriteCloser struct {
	bytes.Buffer
}

// NewBufferWriteCloser creates a new bufferWriteCloser.
func NewBufferWriteCloser() *bufferWriteCloser {
	return &bufferWriteCloser{
		Buffer: bytes.Buffer{},
	}
}

// Close does nothing as bytes.Buffer does not need any cleanup.
func (bwc *bufferWriteCloser) Close() error {
	return nil
}

func TestJSONLogger(t *testing.T) {
	scenario := &model.Scenario{
		ID: "test-scenario",
	}

	buf := NewBufferWriteCloser()
	jsonLogger := logger.NewJSONLogger(buf, scenario)

	alertLog := &model.AlertLog{
		Alert: model.Alert{
			ID: "test-alert",
		},
	}

	alertLogger := jsonLogger.NewAlertLogger(alertLog)
	actionLog := &model.ActionLog{
		Action: model.Action{ID: "test-action"},
	}

	alertLogger.Log(actionLog)

	err := jsonLogger.Flush()
	gt.NoError(t, err)

	var resultLog model.ScenarioLog
	err = json.Unmarshal(buf.Bytes(), &resultLog)
	gt.NoError(t, err)

	gt.V(t, scenario.ID).Equal(resultLog.ID)
	gt.A(t, resultLog.AlertLog).Length(1)
	gt.V(t, alertLog.Alert.ID).Equal(resultLog.AlertLog[0].Alert.ID)
	gt.A(t, resultLog.AlertLog[0].Actions).Length(1)
	gt.V(t, actionLog.Action.ID).Equal(resultLog.AlertLog[0].Actions[0].Action.ID)
}
