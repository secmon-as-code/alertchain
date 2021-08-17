package usecase

import (
	"context"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

func (x *usecase) GetAlerts(ctx context.Context) ([]*ent.Alert, error) {
	return x.clients.DB.GetAlerts(ctx)
}

func (x *usecase) GetAlert(ctx context.Context, id types.AlertID) (*ent.Alert, error) {
	return x.clients.DB.GetAlert(ctx, id)
}

func (x *usecase) RecvAlert(ctx context.Context, recvAlert *alertchain.Alert) (*alertchain.Alert, error) {
	if err := validateAlert(recvAlert); err != nil {
		return nil, goerr.Wrap(err)
	}

	created, err := x.clients.DB.NewAlert(ctx)
	if err != nil {
		return nil, err
	}

	if err := x.clients.DB.UpdateAlert(ctx, created.ID, &recvAlert.Alert); err != nil {
		return nil, err
	}

	attrs := make([]*ent.Attribute, len(recvAlert.Attributes))
	for i, attr := range recvAlert.Attributes {
		attrs[i] = &attr.Attribute
	}
	if err := x.clients.DB.AddAttributes(ctx, created.ID, attrs); err != nil {
		return nil, err
	}

	newAlert, err := x.clients.DB.GetAlert(ctx, created.ID)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := executeChain(ctx, x.chain, newAlert.ID, x.clients); err != nil {
			utils.OutputError(logger, err)
		}
	}()

	return alertchain.NewAlert(newAlert, x.clients.DB), nil
}

func executeChain(ctx context.Context, chain *alertchain.Chain, alertID types.AlertID, clients infra.Clients) error {
	for _, stage := range chain.Stages {
		alert, err := clients.DB.GetAlert(ctx, alertID)
		if err != nil {
			return err
		}

		wg := sync.WaitGroup{}
		results := make([]error, len(stage.Tasks))

		for i, task := range stage.Tasks {
			wg.Add(1)
			args := new(ent.Alert)
			copier.Copy(&args, alert)

			go func(t alertchain.Task, input *alertchain.Alert, idx int) {
				if err := t.Execute(ctx, input); err != nil {
					results[idx] = goerr.Wrap(err).With("task.Name", t.Name())
					utils.OutputError(logger, results[idx])
				}
			}(task, alertchain.NewAlert(args, clients.DB), i)
		}

		wg.Wait()
		for _, err := range results {
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func validateAlert(alert *alertchain.Alert) error {
	if alert.Title == "" {
		return goerr.Wrap(types.ErrInvalidInput, "'title' field is required")
	}
	if alert.Detector == "" {
		return goerr.Wrap(types.ErrInvalidInput, "'detector' field is required")
	}

	for _, attr := range alert.Attributes {
		if attr.Key == "" {
			return goerr.Wrap(types.ErrInvalidInput, "'key' field is required").With("attr", attr)
		}
		if attr.Value == "" {
			return goerr.Wrap(types.ErrInvalidInput, "'value' field is required").With("attr", attr)
		}

		if err := attr.Type.IsValid(); err != nil {
			return goerr.Wrap(err).With("attr", attr)
		}

		for _, s := range attr.Context {
			ctx := types.AttrContext(s)
			if err := ctx.IsValid(); err != nil {
				return goerr.Wrap(err).With("attr", attr)
			}
		}
	}

	return nil
}
