package logger

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

func NewJSONLogger(w io.WriteCloser, s *model.Scenario) *JSONLogger {
	return &JSONLogger{
		w:   w,
		log: s.ToLog(),
	}
}

func (x *JSONLogger) NewAlertLogger(log *model.AlertLog) interfaces.AlertLogger {
	x.log.AlertLog = append(x.log.AlertLog, log)
	return &JSONAlertLogger{
		alertLog: log,
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

type JSONAlertLogger struct {
	alertLog *model.AlertLog
}

func (x *JSONAlertLogger) Log(log *model.ActionLog) {
	x.alertLog.Actions = append(x.alertLog.Actions, log)
}
