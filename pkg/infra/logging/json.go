package logging

import (
	"encoding/json"
	"io"

	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/goerr"
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
	x.log.Error = err.Error()
}

func (x *JSONLogger) Flush() error {
	if err := json.NewEncoder(x.w).Encode(x.log); err != nil {
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
	log := &model.ActionLog{
		Seq: x.seq,
	}
	x.seq++

	x.log.Actions = append(x.log.Actions, log)
	return &JSONActionLogger{
		log: log,
	}
}

type JSONActionLogger struct {
	log *model.ActionLog
}

// LogInit implements interfaces.AlertLogger.
func (x *JSONActionLogger) LogInit(logs []model.Next) {
	x.log.Init = append(x.log.Init, logs...)
}

// LogRun implements interfaces.AlertLogger.
func (x *JSONActionLogger) LogRun(logs []model.Action) {
	x.log.Run = append(x.log.Run, logs...)
}

// LogExit implements interfaces.AlertLogger.
func (x *JSONActionLogger) LogExit(logs []model.Next) {
	x.log.Exit = append(x.log.Exit, logs...)
}
