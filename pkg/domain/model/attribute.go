package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type Attribute struct {
	ID      types.AttributeID
	AlertID types.AlertID

	Key      string
	Value    string
	Type     types.AttrType
	Contexts Contexts

	Annotations []*Annotation `spanner:"-"`
}

func (x *Alert) NewAttribute(key, value string, t types.AttrType, ctxs ...types.AttrContext) *Attribute {
	return &Attribute{
		ID:      types.AttributeID(uuid.NewString()),
		AlertID: x.ID,

		Key:      key,
		Value:    value,
		Type:     t,
		Contexts: ctxs,
	}
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
	ID          types.AnnotationID
	AlertID     types.AlertID
	AttributeID types.AttributeID

	Timestamp *time.Time
	Source    string
	Name      string
	Value     string
	Tags      []string
	URI       string
}

func (x *Attribute) NewAnnotation(src, name, value string, t types.AttrType) *Annotation {
	return &Annotation{
		ID:          types.AnnotationID(uuid.NewString()),
		AlertID:     x.AlertID,
		AttributeID: x.ID,

		Source: src,
		Name:   name,
		Value:  value,
	}
}
