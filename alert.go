package alertchain

import (
	"github.com/google/uuid"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type Alert struct {
	ent.Alert
	db       infra.DBClient
	newAttrs []*Attribute
}

func (x *Alert) Attributes() Attributes {
	attrs := make(Attributes, len(x.Edges.Attributes))
	for i, attr := range x.Edges.Attributes {
		attrs[i] = &Attribute{
			Attribute: *attr,
			alert:     x,
		}
	}
	return attrs
}

func (x *Alert) AddAttributes(attrs []*Attribute) {
	x.newAttrs = append(x.newAttrs, attrs...)
}

func NewAlert(db infra.DBClient) *Alert {
	return &Alert{
		Alert: ent.Alert{
			ID: types.AlertID(uuid.New().String()),
		},
		db: db,
	}
}

func (x *Alert) Commit() error {
	if err := x.db.SaveAlert(&x.Alert); err != nil {
		return err
	}

	attrs := make([]*ent.Attribute, len(x.newAttrs))
	for i, a := range x.newAttrs {
		attrs[i] = &a.Attribute
	}
	if err := x.db.AddAttributes(&x.Alert, attrs); err != nil {
		return err
	}

	for _, attr := range x.Attributes() {
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
