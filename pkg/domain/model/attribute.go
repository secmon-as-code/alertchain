package model

import (
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

type Attribute struct {
	ID      types.AttrID    `json:"id" firestore:"id"`
	Key     types.AttrKey   `json:"key" firestore:"key"`
	Value   types.AttrValue `json:"value" firestore:"value"`
	Type    types.AttrType  `json:"type" firestore:"type"`
	Persist bool            `json:"persist" firestore:"persist"`
	TTL     int             `json:"ttl" firestore:"ttl"`
}

func (x Attribute) Copy() Attribute {
	copied := x
	return copied
}

type Attributes []Attribute

func (x Attributes) Copy() Attributes {
	newAttrs := make(Attributes, len(x))
	for i, p := range x {
		newAttrs[i] = Attribute{
			ID:      p.ID,
			Key:     p.Key,
			Value:   p.Value,
			Type:    p.Type,
			TTL:     p.TTL,
			Persist: p.Persist,
		}
	}
	return newAttrs
}

func (x Attributes) Tidy() Attributes {
	var ret Attributes

	idMap := map[types.AttrID]int{}

	for _, p := range x {
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
