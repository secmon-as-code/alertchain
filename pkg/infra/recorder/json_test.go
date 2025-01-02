package recorder_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/infra/recorder"
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
	jsonLogger := recorder.NewJsonRecorder(buf, scenario)

	alert := model.Alert{
		ID: "test-alert",
	}

	AlertRecorder := jsonLogger.NewAlertRecorder(&alert)

	// first process
	ActionRecorder := AlertRecorder.NewActionRecorder()
	ActionRecorder.Add(model.Action{
		ID:   "test-action",
		Name: "test-action-name",
	})

	// second process, but not action recorded
	_ = AlertRecorder.NewActionRecorder()

	err := jsonLogger.Flush()
	gt.NoError(t, err)

	var resultLog model.ScenarioLog
	err = json.Unmarshal(buf.Bytes(), &resultLog)
	gt.NoError(t, err)

	gt.V(t, scenario.ID).Equal(resultLog.ID)
	gt.A(t, resultLog.Results).Length(1)

	r := resultLog.Results[0]
	gt.V(t, r.Alert.ID).Equal("test-alert")
	gt.A(t, r.Actions).Length(1)

	gt.N(t, r.Actions[0].Seq).Equal(0)
	gt.V(t, r.Actions[0].ID).Equal("test-action")
	gt.V(t, r.Actions[0].Name).Equal("test-action-name")
}
