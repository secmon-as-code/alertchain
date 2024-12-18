package interfaces

import (
	"context"
	"time"

	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

// RunAction is a function to run an action. The function is registered as an option within the chain.Chain.
type RunAction func(ctx context.Context, alert model.Alert, args model.ActionArgs) (any, error)

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
	LogRun(logs []model.Action)
}

// AlertHandler is a function to handle the alert from data source. The handler is registered as an option within the chain.Chain.
type AlertHandler func(ctx context.Context, schema types.Schema, data any) ([]*model.Alert, error)

type Env func() types.EnvVars

type TxProc func(ctx context.Context, input model.Attributes) (model.Attributes, error)

type Database interface {
	GetAttrs(ctx context.Context, ns types.Namespace) (model.Attributes, error)
	PutAttrs(ctx context.Context, ns types.Namespace, attrs model.Attributes) error
	PutWorkflow(ctx context.Context, workflow model.WorkflowRecord) error
	GetWorkflows(ctx context.Context, offset, limit int) ([]model.WorkflowRecord, error)
	GetWorkflow(ctx context.Context, id types.WorkflowID) (*model.WorkflowRecord, error)
	PutAlert(ctx context.Context, alert model.Alert) error
	GetAlert(ctx context.Context, id types.AlertID) (*model.Alert, error)
	Lock(ctx context.Context, ns types.Namespace, timeout time.Time) error
	Unlock(ctx context.Context, ns types.Namespace) error
	Close() error
}
