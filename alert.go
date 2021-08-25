package alertchain

import (
	"context"
	"time"

	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type Alert struct {
	ent.Alert
	Attributes []*Attribute `json:"attributes"`

	id          types.AlertID // Immutable AlertID copied from ent.Alert.ID
	db          infra.DBClient
	newAttrs    []*Attribute
	newStatus   *types.AlertStatus
	newSeverity *types.Severity

	// To remove "edges" in JSON. DO NOT USE as data field
	EdgesOverride interface{} `json:"edges,omitempty"`
}

func (x *Alert) AddAttributes(attrs []*Attribute) {
	x.newAttrs = append(x.newAttrs, attrs...)
}

func (x *Alert) UpdateStatus(status types.AlertStatus) {
	x.newStatus = &status
}

func (x *Alert) UpdateSeverity(sev types.Severity) {
	x.newSeverity = &sev
}

func NewAlert(alert *ent.Alert, db infra.DBClient) *Alert {
	newAlert := &Alert{
		Alert: *alert,
		id:    alert.ID,
		db:    db,
	}
	if len(alert.Edges.Attributes) > 0 {
		attrs := make(Attributes, len(alert.Edges.Attributes))
		for i, attr := range alert.Edges.Attributes {
			attrs[i] = &Attribute{
				Attribute: *attr,
				alert:     newAlert,
			}
		}
		newAlert.Attributes = attrs
	}

	return newAlert
}

func (x *Alert) Commit(ctx context.Context) error {
	ts := time.Now().UTC().Unix()
	if x.newStatus != nil {
		if err := x.db.UpdateAlertStatus(ctx, x.id, *x.newStatus, ts); err != nil {
			return err
		}
	}
	if x.newSeverity != nil {
		if err := x.db.UpdateAlertSeverity(ctx, x.id, *x.newSeverity, ts); err != nil {
			return err
		}
	}

	if len(x.newAttrs) > 0 {
		attrs := make([]*ent.Attribute, len(x.newAttrs))
		for i, a := range x.newAttrs {
			attrs[i] = &a.Attribute
		}
		if err := x.db.AddAttributes(ctx, x.id, attrs); err != nil {
			return err
		}

		for _, attr := range x.Attributes {
			if len(attr.newAnnotations) == 0 {
				continue
			}

			annotations := make([]*ent.Annotation, len(attr.newAnnotations))
			for i, ann := range attr.newAnnotations {
				annotations[i] = &ann.Annotation
			}
			if err := x.db.AddAnnotation(ctx, &attr.Attribute, annotations); err != nil {
				return err
			}
		}
	}

	return nil
}

func (x *Alert) Abort() error {
	return nil
}

func (x *Alert) Close() error {
	return nil
}
