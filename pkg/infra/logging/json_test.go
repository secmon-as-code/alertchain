package logging_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/infra/logging"
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
	jsonLogger := logging.NewJSONLogger(buf, scenario)

	alert := model.Alert{
		ID: "test-alert",
	}

	alertLogger := jsonLogger.NewAlertLogger(&alert)

	// first process
	actionLogger := alertLogger.NewActionLogger()
	actionLogger.LogRun([]model.Action{
		{
			ID:   "test-action",
			Name: "test-action-name",
		},
	})

	// second process
	actionLogger = alertLogger.NewActionLogger()

	err := jsonLogger.Flush()
	gt.NoError(t, err)

	var resultLog model.ScenarioLog
	err = json.Unmarshal(buf.Bytes(), &resultLog)
	gt.NoError(t, err)

	gt.V(t, scenario.ID).Equal(resultLog.ID)
	gt.A(t, resultLog.Results).Length(1)

	r := resultLog.Results[0]
	gt.V(t, r.Alert.ID).Equal("test-alert")
	gt.A(t, r.Actions).Length(2)

	gt.N(t, r.Actions[0].Seq).Equal(0)
	gt.A(t, r.Actions[0].Run).Length(1).At(0, func(t testing.TB, v model.Action) {
		gt.V(t, v.ID).Equal("test-action")
		gt.V(t, v.Name).Equal("test-action-name")
	})

	gt.N(t, r.Actions[1].Seq).Equal(1)
	gt.A(t, r.Actions[1].Run).Length(0)
}
