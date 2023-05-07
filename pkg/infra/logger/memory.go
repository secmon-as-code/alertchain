package logger

import (
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
)

type Memory struct {
	Log model.ScenarioLog
}

func NewMemory(s *model.Scenario) *Memory {
	return &Memory{
		Log: s.ToLog(),
	}
}

func (x *Memory) NewAlertLogger(log *model.AlertLog) interfaces.AlertLogger {
	x.Log.Results = append(x.Log.Results, log)
	return &MemoryAlertLogger{
		alertLog: log,
	}
}

func (x *Memory) LogError(err error) {
	x.Log.Error = err.Error()
}

func (x *Memory) Flush() error {
	return nil
}

type MemoryAlertLogger struct {
	alertLog *model.AlertLog
}

func (x *MemoryAlertLogger) Log(log *model.ActionLog) {
	x.alertLog.Actions = append(x.alertLog.Actions, log)
}
