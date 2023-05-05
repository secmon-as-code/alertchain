package interfaces

import (
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

// RunAction is a function to run an action. The function is registered as an option within the chain.Chain.
type RunAction func(ctx *model.Context, alert model.Alert, args model.ActionArgs) (any, error)

// ActionMock is an interface for "play" mode. The mock should be registered as an option within the chain.Chain. This mock only returns the prepared result for each action ID.
type ActionMock interface {
	GetResult(name types.ActionName) any
}

// ScenarioLogger records the "play" result of the alert chain, which is used for debugging and testing purposes. A logger should be created by the LoggerFactory for each scenario. The LoggerFactory is registered as an option within the chain.Chain.
type ScenarioLogger interface {
	NewAlertLogger(log *model.AlertLog) AlertLogger
	Flush() error
}

// AlertLogger records the "play" Action results of the chain, which is used for debugging and testing purposes. An AlertLogger should be created by the ScenarioLogger for each alert. The ScenarioLogger is registered as an option within the chain.Chain.
type AlertLogger interface {
	Log(log *model.ActionLog)
}

// Router is a function to route the alert to the next action. The router is registered as an option within the chain.Chain.
type Router func(ctx *model.Context, schema types.Schema, data any) error
