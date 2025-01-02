package interfaces

import (
	"context"

	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

// ActionMock is an interface for "play" mode. The mock should be registered as an option within the chain.Chain. This mock only returns the prepared result for each action ID.
type ActionMock interface {
	GetResult(name types.ActionName) any
}

// ScenarioRecorder records the "play" result of the alert chain, which is used for debugging and testing purposes. A logger should be created by the LoggerFactory for each scenario. The LoggerFactory is registered as an option within the chain.Chain.
type ScenarioRecorder interface {
	NewAlertRecorder(alert *model.Alert) AlertRecorder
	LogError(err error)
	Flush() error
}

// AlertRecorder records the "play" Action results of the chain, which is used for debugging and testing purposes. An AlertRecorder should be created by the ScenarioRecorder for each alert. The ScenarioRecorder is registered as an option within the chain.Chain.
type AlertRecorder interface {
	NewActionRecorder() ActionRecorder
}

// ActionRecorder records the "play" result of each action, which is used for debugging and testing purposes. An ActionRecorder should be created by the AlertRecorder for each action. The AlertRecorder is registered as an option within the chain.Chain.
type ActionRecorder interface {
	Add(action model.Action)
}

// AlertHandler is a function to handle the alert from data source. The handler is registered as an option within the chain.Chain.
type AlertHandler func(ctx context.Context, schema types.Schema, data any) ([]*model.Alert, error)

type Env func() types.EnvVars
