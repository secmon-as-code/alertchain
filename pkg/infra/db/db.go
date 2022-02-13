package db

import (
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type Client interface {
	GetAlert(ctx *types.Context, id types.AlertID) (*model.Alert, error)
	SaveAlert(ctx *types.Context, alert *model.Alert) error
	UpdateAlertSeverity(ctx *types.Context, id types.AlertID, sev types.Severity) error
	UpdateAlertClosedAt(ctx *types.Context, id types.AlertID, closedAt time.Time) error
	AddAnnotation(ctx *types.Context, annotations []*model.Annotation) error
	AddAttributes(ctx *types.Context, attrs []*model.Attribute) error
	AddReferences(ctx *types.Context, refs []*model.Reference) error
}
