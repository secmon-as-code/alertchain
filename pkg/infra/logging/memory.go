package logging

import (
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
)

type Memory struct {
	Log model.ScenarioLog
}

// NewAlertLogger implements interfaces.ScenarioLogger.
func (x *Memory) NewAlertLogger(alert *model.Alert) interfaces.AlertLogger {
	log := &model.PlayLog{
		Alert: *alert,
	}
	x.Log.Results = append(x.Log.Results, log)

	return &MemoryAlertLogger{
		log: log,
	}
}

var _ interfaces.ScenarioLogger = &Memory{}

func NewMemory(s *model.Scenario) *Memory {
	return &Memory{
		Log: s.ToLog(),
	}
}

func (x *Memory) LogError(err error) {
	x.Log.Error = err.Error()
}

func (x *Memory) Flush() error {
	return nil
}

type MemoryAlertLogger struct {
	seq int
	log *model.PlayLog
}

func (x *MemoryAlertLogger) NewActionLogger() interfaces.ActionLogger {
	log := &model.ActionLog{
		Seq: x.seq,
	}
	x.seq++

	x.log.Actions = append(x.log.Actions, log)
	return &MemoryActionLogger{
		log: log,
	}
}

type MemoryActionLogger struct {
	log *model.ActionLog
}

// LogInit implements interfaces.AlertLogger.
func (x *MemoryActionLogger) LogInit(logs []model.Next) {
	x.log.Init = append(x.log.Init, logs...)
}

// LogExit implements interfaces.AlertLogger.
func (x *MemoryActionLogger) LogExit(logs []model.Next) {
	x.log.Exit = append(x.log.Exit, logs...)
}

// LogRun implements interfaces.AlertLogger.
func (x *MemoryActionLogger) LogRun(logs []model.Action) {
	x.log.Run = append(x.log.Run, logs...)
}

var _ interfaces.AlertLogger = &MemoryAlertLogger{}
