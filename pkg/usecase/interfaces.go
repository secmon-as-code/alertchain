package usecase

import (
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type Usecase interface {
	RecvAlert(alert *ent.Alert) (*ent.Alert, error)
	GetAlerts() ([]*ent.Alert, error)
	GetAlert(id types.AlertID) (*ent.Alert, error)
}
