package logging

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

var _ interfaces.ScenarioLogger = &JSONLogger{}

func NewJSONLogger(w io.WriteCloser, s *model.Scenario) *JSONLogger {
	return &JSONLogger{
		w:   w,
		log: s.ToLog(),
	}
}

func (x *JSONLogger) NewAlertLogger(alert *model.Alert) interfaces.AlertLogger {
	copied := alert.Copy()

	// Remove redundant data from alert
	copied.Data = nil
	copied.Raw = ""

	log := &model.PlayLog{
		Alert: copied,
	}
	x.log.Results = append(x.log.Results, log)

	return &JSONAlertLogger{
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

var _ interfaces.AlertLogger = &JSONAlertLogger{}

type JSONAlertLogger struct {
	seq int
	log *model.PlayLog
}

// NewJSONActionLogger implements interfaces.AlertLogger.
func (x *JSONAlertLogger) NewActionLogger() interfaces.ActionLogger {
	logger := &JSONActionLogger{
		seq: x.seq,
		log: x.log,
	}

	x.seq++
	return logger
}

type JSONActionLogger struct {
	seq int
	log *model.PlayLog
}

// LogRun implements interfaces.AlertLogger.
func (x *JSONActionLogger) LogRun(logs []model.Action) {
	for _, log := range logs {
		x.log.Actions = append(x.log.Actions, &model.ActionLog{
			Seq:    x.seq,
			Action: log,
		})
	}
}
