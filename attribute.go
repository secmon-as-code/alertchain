package alertchain

import (
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/usecase"
	"github.com/m-mizutani/alertchain/types"
)

type Attribute struct {
	ent.Attribute
	Annotations []*Annotation     `json:"annotations"`
	Actions     []*usecase.Action `json:"actions"`

	Alert *Alert `json:"-"`

	changeRequest *usecase.ChangeRequest

	// To remove "edges" in JSON. DO NOT USE as data field
	EdgesOverride interface{} `json:"edges,omitempty"`
}

func (x *Alert) pushAttribute(attr *ent.Attribute) *Attribute {
	created := newAttribute(attr)
	x.Attributes = append(x.Attributes, created)

	created.Alert = x
	created.changeRequest = &x.ChangeRequest

	return created
}

func newAttribute(attr *ent.Attribute) *Attribute {
	annotations := make([]*Annotation, len(attr.Edges.Annotations))
	for j, ann := range attr.Edges.Annotations {
		annotations[j] = &Annotation{Annotation: *ann}
	}

	created := &Attribute{
		Attribute:     *attr,
		Alert:         NewAlert(attr.Edges.Alert),
		Annotations:   annotations,
		changeRequest: &usecase.ChangeRequest{},
	}
	return created
}

func (x *Attribute) HasContext(ctx types.AttrContext) bool {
	for _, c := range x.Context {
		if c == string(ctx) {
			return true
		}
	}
	return false
}

type Attributes []*Attribute

func (x Attributes) toEnt() []*ent.Attribute {
	resp := make([]*ent.Attribute, len(x))
	for i := range x {
		resp[i] = &x[i].Attribute
	}
	return resp
}

func (x Attributes) FindByKey(key string) Attributes {
	var resp Attributes
	for _, attr := range x {
		if attr.Key == key {
			resp = append(resp, attr)
		}
	}
	return resp
}

func (x Attributes) FindByValue(value string) Attributes {
	var resp Attributes
	for _, attr := range x {
		if attr.Value == value {
			resp = append(resp, attr)
		}
	}
	return resp
}

func (x Attributes) FindByType(attrType types.AttrType) Attributes {
	var resp Attributes
	for _, attr := range x {
		if attr.Type == attrType {
			resp = append(resp, attr)
		}
	}
	return resp

}

func (x Attributes) FindByContext(ctx types.AttrContext) Attributes {
	var resp Attributes
	for _, attr := range x {
		if attr.HasContext(ctx) {
			resp = append(resp, attr)
		}
	}
	return resp
}

func (x *Attribute) Annotate(ann *Annotation) {
	x.changeRequest.AddAnnotation(&x.Attribute, &ann.Annotation)
}

type Annotation struct {
	ent.Annotation
}
