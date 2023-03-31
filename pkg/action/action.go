package action

import (
	"github.com/m-mizutani/alertchain/pkg/action/github"
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

var factoryMap = map[types.ActionName]interfaces.ActionFactory{}

func init() {
	factories := []interfaces.ActionFactory{
		&github.Factory{},
	}

	for _, fac := range factories {
		factoryMap[fac.Name()] = fac
	}
}

func New(cfg model.ActionConfig) (interfaces.Action, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	factory, ok := factoryMap[cfg.Name]
	if !ok {
		return nil, goerr.Wrap(types.ErrNoSuchActionName).With("name", cfg.Name)
	}

	return factory.New(cfg.ID, cfg.Config)
}
