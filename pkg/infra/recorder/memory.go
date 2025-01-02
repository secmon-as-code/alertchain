package recorder

import (
	"github.com/secmon-lab/alertchain/pkg/domain/interfaces"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
)

type Memory struct {
	Log model.ScenarioLog
}

// NewAlertRecorder implements interfaces.ScenarioRecorder.
func (x *Memory) NewAlertRecorder(alert *model.Alert) interfaces.AlertRecorder {
	log := &model.PlayLog{
		Alert: *alert,
	}
	x.Log.Results = append(x.Log.Results, log)

	return &MemoryAlertRecorder{
		log: log,
	}
}

var _ interfaces.ScenarioRecorder = &Memory{}

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

type MemoryAlertRecorder struct {
	seq int
	log *model.PlayLog
}

func (x *MemoryAlertRecorder) NewActionRecorder() interfaces.ActionRecorder {
	logger := &MemoryActionRecorder{
		seq: x.seq,
		log: x.log,
	}
	x.seq++

	return logger
}

type MemoryActionRecorder struct {
	seq int
	log *model.PlayLog
}

// LogRun implements interfaces.AlertRecorder.
func (x *MemoryActionRecorder) Add(log model.Action) {
	x.log.Actions = append(x.log.Actions, &model.ActionLog{
		Seq:    x.seq,
		Action: log,
	})
}

var _ interfaces.AlertRecorder = &MemoryAlertRecorder{}
