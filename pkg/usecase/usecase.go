package usecase

import (
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

// var logger = utils.Logger

type Interface interface {
	GetAlerts(ctx *types.Context) ([]*ent.Alert, error)
	GetAlert(ctx *types.Context, id types.AlertID) (*ent.Alert, error)
	HandleAlert(ctx *types.Context, alert *ent.Alert, attrs []*ent.Attribute) (*ent.Alert, error)

	// Action
	GetExecutableActions(attr *ent.Attribute) []*Action
	GetActionLog(ctx *types.Context, actionLogID int) (*ent.ActionLog, error)
	ExecuteAction(ctx *types.Context, actionID string, attrID int) (*ent.ActionLog, error)
}

type usecase struct {
	clients infra.Clients
	jobs    []*Job
	actions map[string]*Action
}

func New(clients infra.Clients, jobs []*Job, actions []*Action) Interface {
	uc := &usecase{
		clients: clients,
		jobs:    jobs,
		actions: make(map[string]*Action),
	}

	for _, action := range actions {
		uc.actions[action.ID] = action
	}

	return uc
}
