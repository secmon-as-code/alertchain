package recorder

import (
	"encoding/json"
	"io"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/alertchain/pkg/domain/interfaces"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
)

type JSONLogger struct {
	w   io.WriteCloser
	log model.ScenarioLog
}

var _ interfaces.ScenarioRecorder = &JSONLogger{}

func NewJsonRecorder(w io.WriteCloser, s *model.Scenario) *JSONLogger {
	return &JSONLogger{
		w:   w,
		log: s.ToLog(),
	}
}

func (x *JSONLogger) NewAlertRecorder(alert *model.Alert) interfaces.AlertRecorder {
	copied := alert.Copy()

	// Remove redundant data from alert
	copied.Data = nil
	copied.Raw = ""

	log := &model.PlayLog{
		Alert: copied,
	}
	x.log.Results = append(x.log.Results, log)

	return &JSONAlertRecorder{
		log: log,
	}
}

func (x *JSONLogger) LogError(err error) {
	if gErr := goerr.Unwrap(err); gErr != nil {
		x.log.Error = gErr.Printable()
		return
	}

	x.log.Error = err.Error()
}

func (x *JSONLogger) Flush() error {
	encoder := json.NewEncoder(x.w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(x.log); err != nil {
		return goerr.Wrap(err, "Failed to encode JSON scenario log")
	}

	return nil
}

var _ interfaces.AlertRecorder = &JSONAlertRecorder{}

type JSONAlertRecorder struct {
	seq int
	log *model.PlayLog
}

// NewJSONActionRecorder implements interfaces.AlertRecorder.
func (x *JSONAlertRecorder) NewActionRecorder() interfaces.ActionRecorder {
	logger := &JSONActionRecorder{
		seq: x.seq,
		log: x.log,
	}

	x.seq++
	return logger
}

type JSONActionRecorder struct {
	seq int
	log *model.PlayLog
}

// Add implements interfaces.AlertRecorder.
func (x *JSONActionRecorder) Add(log model.Action) {
	x.log.Actions = append(x.log.Actions, &model.ActionLog{
		Seq:    x.seq,
		Action: log,
	})
}
