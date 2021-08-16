package usecase

import (
	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

func (x *usecase) GetAlerts() ([]*ent.Alert, error) {
	return x.infra.DB.GetAlerts()
}

func (x *usecase) GetAlert(id types.AlertID) (*ent.Alert, error) {
	return x.infra.DB.GetAlert(id)
}

func (x *usecase) RecvAlert(alert *alertchain.Alert) (*ent.Alert, error) {
	if err := validateAlert(alert); err != nil {
		return nil, goerr.Wrap(err)
	}

	created, err := x.infra.DB.NewAlert()
	if err != nil {
		return nil, err
	}

	if err := x.infra.DB.UpdateAlert(created.ID, &alert.Alert); err != nil {
		return nil, err
	}

	attrs := make([]*ent.Attribute, len(alert.Attributes))
	for i, attr := range alert.Attributes {
		attrs[i] = &attr.Attribute
	}
	if err := x.infra.DB.AddAttributes(created.ID, attrs); err != nil {
		return nil, err
	}

	newAlert, err := x.infra.DB.GetAlert(created.ID)
	if err != nil {
		return nil, err
	}

	return newAlert, nil
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
