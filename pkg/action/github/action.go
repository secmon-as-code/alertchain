package github

import (
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type Action struct {
	id types.ActionID
}

type Factory struct{}

func (x *Factory) Name() types.ActionName {
	return "github"
}

func (x *Factory) New(id types.ActionID, cfg model.ActionConfigValues) (interfaces.Action, error) {
	return &Action{
		id: id,
	}, nil
}

func (x *Action) ID() types.ActionID { return x.id }

func (x *Action) Run(ctx *types.Context, alert model.Alert, params model.ActionParams) error {
	return nil
}
