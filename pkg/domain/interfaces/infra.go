package interfaces

import (
	"context"
	"time"

	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

type GenAI interface {
	Generate(ctx context.Context, prompts ...string) ([]string, error)
}

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
