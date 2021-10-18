package alertchain

import (
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type newAnnotationRequest struct {
	attr *ent.Attribute
	ann  *ent.Annotation
}

type changeRequest struct {
	newAttrs       []*ent.Attribute
	newReferences  []*ent.Reference
	newStatus      *types.AlertStatus
	newSeverity    *types.Severity
	newAnnotations []*newAnnotationRequest
}

func (x *changeRequest) UpdateStatus(status types.AlertStatus) {
	x.newStatus = &status
}

func (x *changeRequest) UpdateSeverity(sev types.Severity) {
	x.newSeverity = &sev
}

func (x *changeRequest) AddAttributes(attrs []*Attribute) {
	for _, attr := range attrs {
		x.newAttrs = append(x.newAttrs, attr.toEnt())
	}
}

func (x *changeRequest) AddReference(ref *Reference) {
	x.newReferences = append(x.newReferences, ref.toEnt())
}

func (x *changeRequest) AddAnnotation(attr *Attribute, ann *Annotation) {
	x.newAnnotations = append(x.newAnnotations, &newAnnotationRequest{
		attr: attr.toEnt(),
		ann:  ann.toEnt(),
	})
}

func (x *changeRequest) commit(ctx *types.Context, client db.Interface, id types.AlertID) error {
	if x.newStatus != nil {
		if err := client.UpdateAlertStatus(ctx, id, *x.newStatus); err != nil {
			return err
		}
	}
	if x.newSeverity != nil {
		if err := client.UpdateAlertSeverity(ctx, id, *x.newSeverity); err != nil {
			return err
		}
	}

	if len(x.newAttrs) > 0 {
		if err := client.AddAttributes(ctx, id, x.newAttrs); err != nil {
			return err
		}
	}

	for _, newAnn := range x.newAnnotations {
		if err := client.AddAnnotation(ctx, newAnn.attr, []*ent.Annotation{newAnn.ann}); err != nil {
			return err
		}
	}

	if len(x.newReferences) > 0 {
		if err := client.AddReferences(ctx, id, x.newReferences); err != nil {
			return err
		}
	}

	return nil
}
