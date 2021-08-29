package usecase

import (
	"context"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/types"
)

func (x *usecase) GetAlerts(ctx context.Context) ([]*alertchain.Alert, error) {
	got, err := x.clients.DB.GetAlerts(ctx)
	if err != nil {
		return nil, err
	}

	alerts := make([]*alertchain.Alert, len(got))
	for i, alert := range got {
		alerts[i] = alertchain.NewAlert(alert, x.clients.DB)
	}

	return alerts, nil
}

func (x *usecase) GetAlert(ctx context.Context, id types.AlertID) (*alertchain.Alert, error) {
	got, err := x.clients.DB.GetAlert(ctx, id)
	if err != nil {
		return nil, err
	}

	alert := alertchain.NewAlert(got, x.clients.DB)

	for _, attr := range alert.Attributes {
		var resp []*alertchain.ActionEntry
		for _, entry := range x.actions {
			if entry.Action.Executable(attr) {
				resp = append(resp, entry)
			}
		}
		attr.Actions = resp
	}
	return alert, nil
}

func (x *usecase) RecvAlert(ctx context.Context, recvAlert *alertchain.Alert) (*alertchain.Alert, error) {
	return x.chain.InvokeTasks(ctx, recvAlert, &x.clients)
}
