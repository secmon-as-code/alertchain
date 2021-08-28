package usecase

import (
	"context"
	"sync"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

func (x *usecase) GetAlerts(ctx context.Context) ([]*ent.Alert, error) {
	return x.clients.DB.GetAlerts(ctx)
}

func (x *usecase) GetAlert(ctx context.Context, id types.AlertID) (*ent.Alert, error) {
	return x.clients.DB.GetAlert(ctx, id)
}

type ctxKey string

const (
	ctxKeyWaitGroup ctxKey = "WaitGroup"
)

func getWaitGroupFromCtx(ctx context.Context) *sync.WaitGroup {
	obj := ctx.Value(ctxKeyWaitGroup)
	if obj == nil {
		return nil
	}
	wg, ok := obj.(*sync.WaitGroup)
	if !ok {
		return nil
	}
	return wg
}

func ContextWithWaitGroup(ctx context.Context) (context.Context, *sync.WaitGroup) {
	wg := new(sync.WaitGroup)
	resp := context.WithValue(ctx, ctxKeyWaitGroup, wg)
	return resp, wg
}

func (x *usecase) RecvAlert(ctx context.Context, recvAlert *alertchain.Alert) (*alertchain.Alert, error) {
	return x.chain.InvokeTasks(ctx, recvAlert, &x.clients)
}
