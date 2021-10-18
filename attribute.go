package alertchain

import (
	"strings"
	"time"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type Attribute struct {
	Key      string
	Value    string
	Type     types.AttrType
	Contexts Contexts

	Annotations []*Annotation

	id    int
	alert *Alert
}

func (x *Attribute) toEnt() *ent.Attribute {
	ctx := make([]string, len(x.Contexts))
	for p, c := range x.Contexts {
		ctx[p] = string(c)
	}

	return &ent.Attribute{
		ID:      x.id,
		Key:     x.Key,
		Value:   x.Value,
		Type:    x.Type,
		Context: ctx,
	}
}

func (x Attributes) toEnt() []*ent.Attribute {
	resp := make([]*ent.Attribute, len(x))
	for i, ref := range x {
		resp[i] = ref.toEnt()
	}
	return resp
}

func newAttributes(alert *Alert, bases []*ent.Attribute) Attributes {
	attrs := make(Attributes, len(bases))
	for i, base := range bases {
		attrs[i] = newAttribute(alert, base)
	}
	return attrs
}

func newAttribute(alert *Alert, base *ent.Attribute) *Attribute {
	ctx := make([]types.AttrContext, len(base.Context))
	for p, c := range base.Context {
		ctx[p] = types.AttrContext(c)
	}
	ann := make([]*Annotation, len(base.Edges.Annotations))
	for p, a := range base.Edges.Annotations {
		ann[p] = newAnnotation(a)
	}

	return &Attribute{
		Key:         base.Key,
		Value:       base.Value,
		Type:        base.Type,
		Contexts:    ctx,
		Annotations: ann,

		id:    base.ID,
		alert: alert,
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

func (x *Attribute) Annotate(ann *Annotation) {
	x.alert.AddAnnotation(x, ann)
}

type Annotation struct {
	Timestamp time.Time
	Source    string
	Name      string
	Value     string
}

func (x *Annotation) toEnt() *ent.Annotation {
	return &ent.Annotation{
		Timestamp: x.Timestamp.UTC().Unix(),
		Source:    x.Source,
		Name:      x.Name,
		Value:     x.Value,
	}
}

func newAnnotation(base *ent.Annotation) *Annotation {
	return &Annotation{
		Timestamp: time.Unix(base.Timestamp, 0).UTC(),
		Source:    base.Source,
		Name:      base.Name,
		Value:     base.Value,
	}
}
