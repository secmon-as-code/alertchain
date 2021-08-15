package alertchain

import "github.com/m-mizutani/alertchain/pkg/infra/ent"

type Attribute struct {
	ent.Attribute
	alert       *Alert
	newFindings []*Finding
}

type Attributes []*Attribute

func (x *Attribute) AddFinding(finding *Finding) {
	x.newFindings = append(x.newFindings, finding)
}
