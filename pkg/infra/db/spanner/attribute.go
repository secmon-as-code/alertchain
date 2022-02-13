package spanner

import (
	"cloud.google.com/go/spanner"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

func (x *Client) AddAttributes(ctx *types.Context, newAttrs []*model.Attribute) error {
	var mutations []*spanner.Mutation
	for _, attr := range newAttrs {
		m, err := spanner.InsertStruct(tblAttribute, attr)
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

func (x *Client) AddAnnotation(ctx *types.Context, annotations []*model.Annotation) error {
	var mutations []*spanner.Mutation
	for _, ann := range annotations {
		m, err := spanner.InsertStruct(tblAnnotation, ann)
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
