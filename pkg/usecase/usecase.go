package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/types"
)

// var logger = utils.Logger

type Interface interface {
	RecvAlert(ctx context.Context, alert *alertchain.Alert) (*alertchain.Alert, error)
	GetAlerts(ctx context.Context) ([]*alertchain.Alert, error)
	GetAlert(ctx context.Context, id types.AlertID) (*alertchain.Alert, error)

	// Action
	GetExecutableActions(ctx context.Context, attr *alertchain.Attribute) ([]*alertchain.ActionEntry, error)
	ExecuteAction(ctx context.Context, actionID string, attrID int) (*alertchain.ActionLog, error)
}

type usecase struct {
	clients infra.Clients
	chain   *alertchain.Chain
	actions map[string]*alertchain.ActionEntry
}

func New(clients infra.Clients, chain *alertchain.Chain) Interface {
	uc := &usecase{
		clients: clients,
		chain:   chain,
		actions: make(map[string]*alertchain.ActionEntry),
	}

	for _, action := range chain.Actions {
		entry := &alertchain.ActionEntry{
			ID:     uuid.NewString(),
			Action: action,
		}
		uc.actions[entry.ID] = entry
	}

	return uc
}
