package alertchain

import (
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type Attribute struct {
	ent.Attribute
	Annotations []*Annotation  `json:"annotations"`
	Actions     []*ActionEntry `json:"actions"`

	Alert *Alert `json:"-"`

	newAnnotations []*Annotation

	// To remove "edges" in JSON. DO NOT USE as data field
	EdgesOverride interface{} `json:"edges,omitempty"`
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
	x.newAnnotations = append(x.newAnnotations, ann)
}

type Annotation struct {
	ent.Annotation
}
