package spanner

import (
	"cloud.google.com/go/spanner"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

func (x *Client) AddReferences(ctx *types.Context, refs []*model.Reference) error {
	var mutations []*spanner.Mutation
	for _, ref := range refs {
		m, err := spanner.InsertStruct(tblReferences, ref)
		if err != nil {
			return goerr.Wrap(err)
		}
		mutations = append(mutations, m)
	}
	if _, err := x.client.Apply(ctx, mutations); err != nil {
		return goerr.Wrap(err)
	}
	return nil
}
