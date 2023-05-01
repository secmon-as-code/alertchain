package chain

import (
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
)

type dummyScenarioLogger struct{}

func (x *dummyScenarioLogger) NewAlertLogger(log *model.AlertLog) interfaces.AlertLogger {
	return &dummyAlertLogger{}
}
func (x *dummyScenarioLogger) Flush() error { return nil }

type dummyAlertLogger struct{}

func (x *dummyAlertLogger) Log(log *model.ActionLog) {}
