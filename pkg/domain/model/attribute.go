package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type Attribute struct {
	ID    types.AttrID        `json:"id" firestore:"id"`
	Key   types.AttrKey       `json:"key" firestore:"key"`
	Value types.AttrValue     `json:"value" firestore:"value"`
	Type  types.AttributeType `json:"type" firestore:"type"`
	TTL   int64               `json:"ttl" firestore:"ttl"`

	// for DB only
	ExpiresAt int64 `json:"-" firestore:"expires_at"`
}

type Attributes []Attribute

func (x Attributes) Copy() Attributes {
	newAttrs := make(Attributes, len(x))
	for i, p := range x {
		newAttrs[i] = Attribute{
			ID:    p.ID,
			Key:   p.Key,
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
