package usecase

import (
	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type usecase struct {
	infra infra.Infra
	chain *alertchain.Chain
}

func New(infra infra.Infra, chain *alertchain.Chain) Usecase {
	return &usecase{
		infra: infra,
		chain: chain,
	}
}

func (x *usecase) Initialize() error {
	panic("not implemented")
}

func (x *usecase) RecvAlert(alert *ent.Alert) (*ent.Alert, error) {
	panic("not implemented")
}

func (x *usecase) GetAlerts() ([]*ent.Alert, error) {
	return x.infra.DB.GetAlerts()
}

func (x *usecase) GetAlert(id types.AlertID) (*ent.Alert, error) {
	return x.infra.DB.GetAlert(id)
}

func (x *usecase) Shutdown() error {
	panic("not implemented")
}
