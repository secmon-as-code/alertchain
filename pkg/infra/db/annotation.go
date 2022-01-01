package db

import (
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
)

func entToAnnotations(base []*ent.Annotation) []*model.Annotation {
	ann := make([]*model.Annotation, len(base))
	for i, a := range base {
		ann[i] = &model.Annotation{
			Source:    a.Source,
			Name:      a.Name,
			Value:     a.Value,
			Timestamp: time.Unix(a.Timestamp, 0),
		}
	}

	return ann
}

func (x *Client) AddAnnotation(ctx *types.Context, attr *model.Attribute, annotations []*model.Annotation) error {
	if len(annotations) == 0 {
		return nil
	}

	builders := make([]*ent.AnnotationCreate, len(annotations))
	for i, ann := range annotations {
		builders[i] = x.client.Annotation.Create().
			SetName(ann.Name).
			SetSource(ann.Source).
			SetValue(ann.Value).
			SetTimestamp(ann.Timestamp.Unix())
	}

	added, err := x.client.Annotation.CreateBulk(builders...).Save(ctx)
	if err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	if _, err := x.client.Attribute.UpdateOneID(attr.ID()).AddAnnotations(added...).Save(ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}
