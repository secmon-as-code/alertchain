package db

import (
	"github.com/m-mizutani/alertchain/gen/ent"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

func entToAttributes(base []*ent.Attribute) []*model.Attribute {
	attrs := make([]*model.Attribute, len(base))
	for i, b := range base {
		ctxSet := make([]types.AttrContext, len(b.Context))
		for c := range b.Context {
			ctxSet[c] = types.AttrContext(b.Context[c])
		}
		attrs[i] = &model.Attribute{
			Key:         b.Key,
			Value:       b.Value,
			Type:        b.Type,
			Contexts:    ctxSet,
			Annotations: entToAnnotations(b.Edges.Annotations),
		}
	}

	return attrs
}

func (x *Client) AddAttributes(ctx *types.Context, id types.AlertID, newAttrs []*model.Attribute) error {
	if len(newAttrs) == 0 {
		return nil // nothing to do
	}

	builders := make([]*ent.AttributeCreate, len(newAttrs))
	for i, attr := range newAttrs {
		ctxSet := make([]string, len(attr.Contexts))
		for i := range attr.Contexts {
			ctxSet[i] = string(attr.Contexts[i])
		}

		builders[i] = x.client.Attribute.Create().
			SetKey(attr.Key).
			SetValue(attr.Value).
			SetType(attr.Type).
			SetContext(ctxSet).
			SetAlertID(id)
	}

	_, err := x.client.Attribute.CreateBulk(builders...).Save(ctx)
	if err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}
