package usecase

import (
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

func ValidateAlert(x *ent.Alert) error {
	if x.Title == "" {
		return goerr.Wrap(types.ErrInvalidInput, "'title' field is required")
	}
	if x.Detector == "" {
		return goerr.Wrap(types.ErrInvalidInput, "'detector' field is required")
	}

	return nil
}

func ValidateAttribute(attr *ent.Attribute) error {
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

	return nil
}
