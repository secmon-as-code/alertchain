package interfaces

import (
	"time"

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
	NewAlertLogger(alert *model.Alert) AlertLogger
	LogError(err error)
	Flush() error
}

// AlertLogger records the "play" Action results of the chain, which is used for debugging and testing purposes. An AlertLogger should be created by the ScenarioLogger for each alert. The ScenarioLogger is registered as an option within the chain.Chain.
type AlertLogger interface {
	NewActionLogger() ActionLogger
}

// ActionLogger records the "play" result of each action, which is used for debugging and testing purposes. An ActionLogger should be created by the AlertLogger for each action. The AlertLogger is registered as an option within the chain.Chain.
type ActionLogger interface {
	LogInit(logs []model.Next)
	LogRun(logs []model.Action)
	LogExit(logs []model.Next)
}

// Router is a function to route the alert to the next action. The router is registered as an option within the chain.Chain.
type Router func(ctx *model.Context, schema types.Schema, data any) error

type Env func() types.EnvVars

type TxProc func(ctx *model.Context, input model.Attributes) (model.Attributes, error)

type Database interface {
	GetAttrs(ctx *model.Context, ns types.Namespace) (model.Attributes, error)
	PutAttrs(ctx *model.Context, ns types.Namespace, attrs model.Attributes) error
	PutWorkflow(ctx *model.Context, workflow model.WorkflowRecord) error
	GetWorkflows(ctx *model.Context, offset, limit int) ([]model.WorkflowRecord, error)
	Lock(ctx *model.Context, ns types.Namespace, timeout time.Time) error
	Unlock(ctx *model.Context, ns types.Namespace) error
	Close() error
}
