package alert

import (
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra"
)

type Service struct {
	alert   *model.Alert
	clients *infra.Clients
}

func New(alert *model.Alert, clients *infra.Clients) *Service {
	return &Service{
		alert:   alert,
		clients: clients,
	}
}

func (x *Service) Alert() *model.Alert { return x.alert }

func (x *Service) HandleChangeRequest(ctx *types.Context, req *model.ChangeRequest) error {
	panic("not implemented")
}

func (x *Service) Refresh(ctx *types.Context) error {
	got, err := x.clients.DB().GetAlert(ctx, x.alert.ID)
	if err != nil {
		return err
	}

	x.alert = got
	return nil
}
