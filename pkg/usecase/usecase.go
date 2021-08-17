package usecase

import (
	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/alertchain/types"
)

var logger = utils.Logger

type Interface interface {
	RecvAlert(alert *alertchain.Alert) (*alertchain.Alert, error)
	GetAlerts() ([]*ent.Alert, error)
	GetAlert(id types.AlertID) (*ent.Alert, error)
}

type usecase struct {
	clients infra.Clients
	chain   *alertchain.Chain
}

func New(clients infra.Clients, chain *alertchain.Chain) Interface {
	return &usecase{
		clients: clients,
		chain:   chain,
	}
}
