package infra

import (
	"context"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type Clients struct {
	DB DBClient
}

type DBClient interface {
	Close() error

	GetAlert(ctx context.Context, id types.AlertID) (*ent.Alert, error)
	GetAlerts(ctx context.Context) ([]*ent.Alert, error)
	NewAlert(ctx context.Context) (*ent.Alert, error)
	UpdateAlert(ctx context.Context, id types.AlertID, alert *ent.Alert) error
	UpdateAlertStatus(ctx context.Context, id types.AlertID, status types.AlertStatus) error
	UpdateAlertSeverity(ctx context.Context, id types.AlertID, status types.Severity) error

	AddAttributes(ctx context.Context, id types.AlertID, newAttrs []*ent.Attribute) error
	AddFindings(ctx context.Context, attr *ent.Attribute, findings []*ent.Finding) error
}
