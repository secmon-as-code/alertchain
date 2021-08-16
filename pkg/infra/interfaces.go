package infra

import (
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type Infra struct {
	DB DBClient
}

type DBClient interface {
	Close() error

	GetAlert(id types.AlertID) (*ent.Alert, error)
	GetAlerts() ([]*ent.Alert, error)
	NewAlert() (*ent.Alert, error)
	SaveAlert(alert *ent.Alert) error

	AddAttributes(alert *ent.Alert, newAttrs []*ent.Attribute) error
	AddFindings(attr *ent.Attribute, findings []*ent.Finding) error
}
