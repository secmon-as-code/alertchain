package interfaces

import (
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type Action interface {
	ID() types.ActionID
	Run(ctx *model.Context, alert model.Alert, args model.ActionArgs) (any, error)
}

type ActionFactory interface {
	Name() types.ActionName
	New(id types.ActionID, cfg model.ActionConfigValues) (Action, error)
}
