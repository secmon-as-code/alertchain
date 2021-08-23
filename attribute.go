package alertchain

import "github.com/m-mizutani/alertchain/pkg/infra/ent"

type Attribute struct {
	ent.Attribute
	alert          *Alert
	newAnnotations []*Annotation
}

type Attributes []*Attribute

func (x *Attribute) AddFinding(ann *Annotation) {
	x.newAnnotations = append(x.newAnnotations, ann)
}

type Annotation struct {
	ent.Annotation
}
