package chain

import (
	"github.com/secmon-lab/alertchain/pkg/domain/interfaces"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
)

type dummyScenarioRecorder struct{}

func (*dummyScenarioRecorder) NewAlertRecorder(alert *model.Alert) interfaces.AlertRecorder {
	return &dummyAlertRecorder{}
}

var _ interfaces.ScenarioRecorder = &dummyScenarioRecorder{}

func (x *dummyScenarioRecorder) LogError(err error) {}
func (x *dummyScenarioRecorder) Flush() error       { return nil }

type dummyAlertRecorder struct{}

// NewActionRecorder implements interfaces.AlertRecorder.
func (*dummyAlertRecorder) NewActionRecorder() interfaces.ActionRecorder {
	return &dummyActionRecorder{}
}

var _ interfaces.AlertRecorder = &dummyAlertRecorder{}

type dummyActionRecorder struct{}

func (*dummyActionRecorder) Add(log model.Action) {}

var _ interfaces.ActionRecorder = &dummyActionRecorder{}
