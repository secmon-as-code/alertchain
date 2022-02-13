package spanner

import (
	"time"

	"cloud.google.com/go/spanner"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

func (x *Client) GetAlert(ctx *types.Context, id types.AlertID) (*model.Alert, error) {
	txn := x.client.ReadOnlyTransaction()
	defer txn.Close()

	var alert *model.Alert

	if err := txn.Query(ctx, spanner.Statement{
		SQL: `SELECT * FROM ` + tblAlert + ` WHERE ID = @ID`,
		Params: map[string]interface{}{
			"ID": id,
		},
	}).Do(func(r *spanner.Row) error {
		var record model.Alert
		if err := r.ToStruct(&record); err != nil {
			return goerr.Wrap(err)
		}
		alert = &record
		return nil
	}); err != nil {
		return nil, goerr.Wrap(err)
	}

	if err := txn.Query(ctx, spanner.Statement{
		SQL: `SELECT * FROM ` + tblAttribute + `@{FORCE_INDEX=AttributesByAlertID} WHERE AlertID = @ID`,
		Params: map[string]interface{}{
			"ID": id,
		},
	}).Do(func(r *spanner.Row) error {
		var record model.Attribute
		if err := r.ToStruct(&record); err != nil {
			return goerr.Wrap(err)
		}
		alert.Attributes = append(alert.Attributes, &record)
		return nil
	}); err != nil {
		return nil, goerr.Wrap(err)
	}

	if err := txn.Query(ctx, spanner.Statement{
		SQL: `SELECT * FROM ` + tblAnnotation + `@{FORCE_INDEX=AnnotationsByAlertID} WHERE AlertID = @ID`,
		Params: map[string]interface{}{
			"ID": id,
		},
	}).Do(func(r *spanner.Row) error {
		var record model.Annotation
		if err := r.ToStruct(&record); err != nil {
			return goerr.Wrap(err)
		}

		for _, attr := range alert.Attributes {
			if attr.ID == record.AttributeID {
				attr.Annotations = append(attr.Annotations, &record)
				break
			}
		}
		return nil
	}); err != nil {
		return nil, goerr.Wrap(err)
	}

	if err := txn.Query(ctx, spanner.Statement{
		SQL: `SELECT * FROM ` + tblReferences + `@{FORCE_INDEX=ReferencesByAlertID} WHERE AlertID = @ID`,
		Params: map[string]interface{}{
			"ID": id,
		},
	}).Do(func(r *spanner.Row) error {
		var record model.Reference
		if err := r.ToStruct(&record); err != nil {
			return goerr.Wrap(err)
		}

		alert.References = append(alert.References, &record)
		return nil
	}); err != nil {
		return nil, goerr.Wrap(err)
	}

	return alert, nil
}

func (x *Client) SaveAlert(ctx *types.Context, alert *model.Alert) error {
	m, err := spanner.InsertStruct(tblAlert, alert)
	if err != nil {
		return goerr.Wrap(err)
	}
	if _, err := x.client.Apply(ctx, []*spanner.Mutation{m}); err != nil {
		return goerr.Wrap(err)
	}
	return nil
}

func (x *Client) UpdateAlertSeverity(ctx *types.Context, id types.AlertID, sev types.Severity) error {
	m := []*spanner.Mutation{
		spanner.Update(tblAlert, []string{"ID", "Severity"}, []interface{}{id, sev}),
	}

	if _, err := x.client.Apply(ctx, m); err != nil {
		return goerr.Wrap(err)
	}

	return nil

}

func (x *Client) UpdateAlertClosedAt(ctx *types.Context, id types.AlertID, ts time.Time) error {
	m := []*spanner.Mutation{
		spanner.Update(tblAlert, []string{"ID", "ClosedAt"}, []interface{}{id, ts}),
	}

	if _, err := x.client.Apply(ctx, m); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}
