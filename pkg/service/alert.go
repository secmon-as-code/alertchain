package service

import (
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

func (x *Service) InvokeChain(ctx *types.Context, alert *model.Alert) error {
	if err := x.clients.DB().PutAlert(ctx, alert); err != nil {
		return err
	}

	return nil
}
