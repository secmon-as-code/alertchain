package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type Alert struct {
	ID types.AlertID

	Title       string
	Description string
	Detector    string
	Severity    types.Severity

	DetectedAt time.Time

	Attributes Attributes `spanner:"-"`
	References References `spanner:"-"`

	CreatedAt time.Time
	ClosedAt  time.Time
}

// NewAlert creates and returns initialized Alert instance. If base is not nil, values will be copied to a new instance. If nil, nothing to do.
func NewAlert(base *Alert) *Alert {
	var alert Alert
	if base != nil {
		alert = *base
	}

	// set initial values
	alert.ID = types.NewAlertID()
	alert.Severity = types.SevUnclassified

	return &alert
}

// NewAlertWithID also creates and returns a new Alert instance, but it can set id.
func NewAlertWithID(id types.AlertID, base *Alert) *Alert {
	alert := NewAlert(base)
	alert.ID = id
	return alert
}

func (x *Alert) Validate() error {
	if err := validation.ValidateStruct(x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Title, validation.Required),
		validation.Field(&x.Detector, validation.Required),
		validation.Field(&x.Severity, validation.Required),
		validation.Field(&x.CreatedAt, validation.Required),
	); err != nil {
		return types.ErrInvalidInput.Wrap(err)
	}

	return nil
}
