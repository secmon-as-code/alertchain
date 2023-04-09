package interfaces

import (
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type Action interface {
	ID() types.ActionID
	Run(ctx *types.Context, alert model.Alert, params model.ActionParams) error
}

type ActionFactory interface {
	Name() types.ActionName
	New(id types.ActionID, cfg model.ActionConfigValues) (Action, error)
}
