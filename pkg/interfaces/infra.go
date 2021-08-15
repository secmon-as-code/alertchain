package interfaces

import (
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type DBClient interface {
	Init(dbType, dbConfig string) error
	Close() error

	GetAlert(id types.AlertID) (*ent.Alert, error)
	NewAlert() (*ent.Alert, error)
	SaveAlert(alert *ent.Alert) error

	AddAttributes(alert *ent.Alert, newAttrs []*ent.Attribute) error
	AddFindings(attr *ent.Attribute, findings []*ent.Finding) error
}

type DBClientFactory func() DBClient
