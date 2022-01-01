package model

import (
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type newAnnotationRequest struct {
	attr *Attribute
	ann  *Annotation
}

type changeRequest struct {
	NewAttrs       []*Attribute
	NewReferences  []*Reference
	NewStatus      *types.AlertStatus
	NewSeverity    *types.Severity
	NewAnnotations []*newAnnotationRequest
}

func (x *changeRequest) UpdateStatus(status types.AlertStatus) {
	x.NewStatus = &status
}

func (x *changeRequest) UpdateSeverity(sev types.Severity) {
	x.NewSeverity = &sev
}

func (x *changeRequest) AddAttributes(attrs []*Attribute) {
	x.NewAttrs = append(x.NewAttrs, attrs...)
}

func (x *changeRequest) AddReference(ref *Reference) {
	x.NewReferences = append(x.NewReferences, ref)
}

func (x *changeRequest) AddAnnotation(attr *Attribute, ann *Annotation) {
	x.NewAnnotations = append(x.NewAnnotations, &newAnnotationRequest{
		attr: attr,
		ann:  ann,
	})
}
