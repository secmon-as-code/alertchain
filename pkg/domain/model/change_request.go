package model

import (
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type AnnotationRequest struct {
	attr *Attribute
	ann  *Annotation
}

type ChangeRequest struct {
	NewAttrs       []*Attribute
	NewReferences  []*Reference
	NewStatus      *types.AlertStatus
	NewSeverity    *types.Severity
	NewAnnotations []*AnnotationRequest
}

func (x *ChangeRequest) UpdateStatus(status types.AlertStatus) {
	x.NewStatus = &status
}

func (x *ChangeRequest) UpdateSeverity(sev types.Severity) {
	x.NewSeverity = &sev
}

func (x *ChangeRequest) AddAttributes(attrs []*Attribute) {
	x.NewAttrs = append(x.NewAttrs, attrs...)
}

func (x *ChangeRequest) AddReference(ref *Reference) {
	x.NewReferences = append(x.NewReferences, ref)
}

func (x *ChangeRequest) AddAnnotation(attr *Attribute, ann *Annotation) {
	x.NewAnnotations = append(x.NewAnnotations, &AnnotationRequest{
		attr: attr,
		ann:  ann,
	})
}
