package db

import (
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
)

func (x *Client) AddReferences(ctx *types.Context, id types.AlertID, refs []*model.Reference) error {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	builders := make([]*ent.ReferenceCreate, len(refs))
	for i, ref := range refs {
		builders[i] = x.client.Reference.Create().
			SetSource(ref.Source).
			SetTitle(ref.Title).
			SetURL(ref.URL).
			SetComment(ref.Comment)
	}

	added, err := x.client.Reference.CreateBulk(builders...).Save(ctx)
	if err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	if _, err := x.client.Alert.UpdateOneID(id).AddReferences(added...).Save(ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func entToReferences(bases []*ent.Reference) model.References {
	resp := make(model.References, len(bases))
	for i, ref := range bases {
		resp[i] = entToReference(ref)
	}
	return resp
}

func entToReference(base *ent.Reference) *model.Reference {
	return &model.Reference{
		Source:  base.Source,
		Title:   base.Title,
		URL:     base.URL,
		Comment: base.Comment,
	}
}
