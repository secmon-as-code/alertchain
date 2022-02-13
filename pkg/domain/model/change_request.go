package model

import (
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type AnnotationRequest struct {
	attr *Attribute
	ann  *Annotation
}

type ChangeRequest struct {
	addingAttrs       []*Attribute
	addingReferences  []*Reference
	addingAnnotations []*AnnotationRequest
	updatingSeverity  *types.Severity
}

func (x *ChangeRequest) AddingAttrs() []*Attribute               { return x.addingAttrs }
func (x *ChangeRequest) AddingReferences() []*Reference          { return x.addingReferences }
func (x *ChangeRequest) AddingAnnotations() []*AnnotationRequest { return x.addingAnnotations }

func (x *ChangeRequest) UpdateSeverity(sev types.Severity) {
	x.updatingSeverity = &sev
}

func (x *ChangeRequest) AddAttributes(attrs []*Attribute) {
	x.addingAttrs = append(x.addingAttrs, attrs...)
}

func (x *ChangeRequest) AddReference(ref *Reference) {
	x.addingReferences = append(x.addingReferences, ref)
}

func (x *ChangeRequest) AddAnnotation(attr *Attribute, ann *Annotation) {
	x.addingAnnotations = append(x.addingAnnotations, &AnnotationRequest{
		attr: attr,
		ann:  ann,
	})
}
