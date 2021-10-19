package alertchain

import (
	"time"

	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

type Alert struct {
	Title       string
	Description string
	Detector    string
	Status      types.AlertStatus
	Severity    types.Severity

	DetectedAt time.Time

	Attributes Attributes
	References References

	base      *ent.Alert
	createdAt time.Time
	closedAt  time.Time

	changeRequest
}

func (x *Alert) ID() types.AlertID {
	if x.base == nil {
		return ""
	}
	return x.base.ID
}
func (x *Alert) CreatedAt() time.Time { return x.createdAt }
func (x *Alert) ClosedAt() time.Time  { return x.closedAt }

func (x *Alert) toEnt() *ent.Alert {
	alert := &ent.Alert{
		Title:       x.Title,
		Description: x.Description,
		Detector:    x.Detector,
		Status:      x.Status,
		Severity:    x.Severity,
		DetectedAt:  x.DetectedAt.Unix(),
		CreatedAt:   x.createdAt.Unix(),
		ClosedAt:    x.closedAt.Unix(),
	}
	if x.base != nil {
		alert.ID = x.base.ID
	}
	return alert
}

func newAlert(base *ent.Alert) *Alert {
	alert := &Alert{
		base:        base,
		Title:       base.Title,
		Description: base.Description,
		Detector:    base.Detector,
		Status:      base.Status,
		Severity:    base.Severity,
		DetectedAt:  time.Unix(base.DetectedAt, 0),
		createdAt:   time.Unix(base.CreatedAt, 0),
		closedAt:    time.Unix(base.ClosedAt, 0),
	}

	alert.Attributes = newAttributes(alert, base.Edges.Attributes)
	alert.References = newReferences(base.Edges.References)

	return alert
}

func (x *Alert) validate() error {
	if x.Title == "" {
		return goerr.Wrap(types.ErrInvalidInput, "'title' field is required")
	}
	if x.Detector == "" {
		return goerr.Wrap(types.ErrInvalidInput, "'detector' field is required")
	}

	return nil
}

func insertAlert(ctx *types.Context, alert *Alert, client db.Interface) (types.AlertID, error) {
	alert.Status = types.StatusNew
	alert.createdAt = time.Now().UTC()

	added, err := client.PutAlert(ctx, alert.toEnt())
	if err != nil {
		return "", err
	}

	if err := client.AddAttributes(ctx, added.ID, alert.Attributes.toEnt()); err != nil {
		return "", err
	}

	if err := client.AddReferences(ctx, added.ID, alert.References.toEnt()); err != nil {
		return "", err
	}

	return added.ID, nil
}
