package usecase

import (
	"context"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/alertchain/types"
)

var logger = utils.Logger

type Interface interface {
	RecvAlert(ctx context.Context, alert *alertchain.Alert) (*alertchain.Alert, error)
	GetAlerts(ctx context.Context) ([]*ent.Alert, error)
	GetAlert(ctx context.Context, id types.AlertID) (*ent.Alert, error)
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
