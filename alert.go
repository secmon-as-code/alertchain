package alertchain

import (
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type Alert struct {
	ent.Alert
	Attributes []*Attribute `json:"attributes"`

	id       types.AlertID // Immutable AlertID copied from ent.Alert.ID
	db       infra.DBClient
	newAttrs []*Attribute
}

func (x *Alert) AddAttributes(attrs []*Attribute) {
	x.newAttrs = append(x.newAttrs, attrs...)
}

func NewAlert(alert *ent.Alert, db infra.DBClient) *Alert {
	newAlert := &Alert{
		Alert: *alert,
		id:    alert.ID,
		db:    db,
	}
	attrs := make(Attributes, len(alert.Edges.Attributes))
	for i, attr := range alert.Edges.Attributes {
		attrs[i] = &Attribute{
			Attribute: *attr,
			alert:     newAlert,
		}
	}

	return newAlert
}

func (x *Alert) Commit() error {
	if err := x.db.UpdateAlert(x.id, &x.Alert); err != nil {
		return err
	}

	attrs := make([]*ent.Attribute, len(x.newAttrs))
	for i, a := range x.newAttrs {
		attrs[i] = &a.Attribute
	}
	if err := x.db.AddAttributes(x.id, attrs); err != nil {
		return err
	}

	for _, attr := range x.Attributes {
		if len(attr.newFindings) == 0 {
			continue
		}

		findings := make([]*ent.Finding, len(attr.newFindings))
		for i, finding := range attr.newFindings {
			findings[i] = &finding.Finding
		}
		if err := x.db.AddFindings(&attr.Attribute, findings); err != nil {
			return err
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
