package usecase

import (
	"time"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type Job struct {
	Timeout   time.Duration
	ExitOnErr bool
	Tasks     []*Task
}

type Task struct {
	Name    string
	Execute func(ctx *types.Context, alert *ent.Alert) (*ChangeRequest, error)
}

type Action struct {
	ID         string
	Name       string
	Executable func(attr *ent.Attribute) bool
	Execute    func(ctx *types.Context, attr *ent.Attribute) error
}

type newAnnotationRequest struct {
	attr *ent.Attribute
	ann  *ent.Annotation
}

type ChangeRequest struct {
	newAttrs       []*ent.Attribute
	newReferences  []*ent.Reference
	newStatus      *types.AlertStatus
	newSeverity    *types.Severity
	newAnnotations []*newAnnotationRequest
}

func (x *ChangeRequest) UpdateStatus(status types.AlertStatus) {
	x.newStatus = &status
}

func (x *ChangeRequest) UpdateSeverity(sev types.Severity) {
	x.newSeverity = &sev
}

func (x *ChangeRequest) AddAttributes(attrs []*ent.Attribute) {
	x.newAttrs = append(x.newAttrs, attrs...)
}

func (x *ChangeRequest) AddReference(ref *ent.Reference) {
	x.newReferences = append(x.newReferences, ref)
}

func (x *ChangeRequest) AddAnnotation(attr *ent.Attribute, ann *ent.Annotation) {
	x.newAnnotations = append(x.newAnnotations, &newAnnotationRequest{
		attr: attr,
		ann:  ann,
	})
}
