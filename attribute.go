package alertchain

import "github.com/m-mizutani/alertchain/pkg/infra/ent"

type Attribute struct {
	ent.Attribute
	alert          *Alert
	newAnnotations []*Annotation

	// To remove "edges" in JSON. DO NOT USE as data field
	EdgesOverride interface{} `json:"edges,omitempty"`
}

type Attributes []*Attribute

func (x *Attribute) AddFinding(ann *Annotation) {
	x.newAnnotations = append(x.newAnnotations, ann)
}

type Annotation struct {
	ent.Annotation
}
