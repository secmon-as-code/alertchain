package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type Attribute struct {
	ID    types.AttrID        `json:"id"`
	Name  types.AttrName      `json:"name"`
	Value types.AttrValue     `json:"value"`
	Type  types.AttributeType `json:"type"`
}

type Attributes []Attribute

func (x Attributes) Copy() Attributes {
	newAttrs := make(Attributes, len(x))
	for i, p := range x {
		newAttrs[i] = Attribute{
			ID:    p.ID,
			Name:  p.Name,
			Value: p.Value,
			Type:  p.Type,
		}
	}
	return newAttrs
}

func TidyAttributes(attrs Attributes) Attributes {
	var ret Attributes

	idMap := map[types.AttrID]int{}

	for _, p := range attrs {
		if p.ID == "" {
			p.ID = types.NewAttrID()
		}

		if _, ok := idMap[p.ID]; ok {
			ret[idMap[p.ID]] = p
		} else {
			ret = append(ret, p)
			idMap[p.ID] = len(ret) - 1
		}
	}

	return ret
}
