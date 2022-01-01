package model

import (
	"strings"
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type Attribute struct {
	Key      string
	Value    string
	Type     types.AttrType
	Contexts Contexts

	Annotations []*Annotation

	id int
}

func (x *Attribute) ID() int { return x.id }

func MakeAttribute(id int, attr *Attribute) *Attribute {
	newAttr := *attr
	newAttr.id = id
	return &newAttr
}

type Contexts []types.AttrContext

func (x Contexts) String() string {
	s := make([]string, len(x))
	for i, c := range x {
		s[i] = string(c)
	}
	return strings.Join(s, ", ")
}

func (x *Attribute) HasContext(ctx types.AttrContext) bool {
	for _, c := range x.Contexts {
		if c == ctx {
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

type Annotation struct {
	Timestamp time.Time
	Source    string
	Name      string
	Value     string
}
