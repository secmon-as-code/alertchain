package logging

import (
	"github.com/secmon-lab/alertchain/pkg/domain/interfaces"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
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
	logger := &MemoryActionLogger{
		seq: x.seq,
		log: x.log,
	}
	x.seq++

	return logger
}

type MemoryActionLogger struct {
	seq int
	log *model.PlayLog
}

// LogRun implements interfaces.AlertLogger.
func (x *MemoryActionLogger) LogRun(logs []model.Action) {
	for _, log := range logs {
		x.log.Actions = append(x.log.Actions, &model.ActionLog{
			Seq:    x.seq,
			Action: log,
		})
	}
}

var _ interfaces.AlertLogger = &MemoryAlertLogger{}
