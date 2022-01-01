package db

import (
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"

	entAlert "github.com/m-mizutani/alertchain/pkg/infra/ent/alert"
)

func entToAlert(alert *ent.Alert) *model.Alert {
	return model.NewAlertWithID(alert.ID, &model.Alert{
		Title:       alert.Title,
		Description: alert.Description,
		Detector:    alert.Detector,
		Status:      alert.Status,
		Severity:    alert.Severity,

		DetectedAt: time.Unix(alert.DetectedAt, 0),

		CreatedAt: time.Unix(alert.CreatedAt, 0),
		ClosedAt:  time.Unix(alert.ClosedAt, 0),

		Attributes: entToAttributes(alert.Edges.Attributes),
		References: entToReferences(alert.Edges.References),
	})
}

func getAlertQuery(client *ent.Client) *ent.AlertQuery {
	return client.Alert.Query().
		WithTaskLogs(func(q *ent.TaskLogQuery) {
			q.WithExecLogs(func(q *ent.ExecLogQuery) {
				q.Order(ent.Desc("timestamp"))
			})
		}).
		WithActionLogs(func(q *ent.ActionLogQuery) {
			q.WithExecLogs(func(q *ent.ExecLogQuery) {
				q.Order(ent.Desc("timestamp"))
			})
		}).
		WithReferences().
		WithAttributes(func(q *ent.AttributeQuery) {
			q.WithAnnotations()
		})
}

func (x *Client) GetAlert(ctx *types.Context, id types.AlertID) (*model.Alert, error) {
	fetched, err := getAlertQuery(x.client).Where(entAlert.ID(id)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return entToAlert(fetched), nil
}

func (x *Client) GetAlerts(ctx *types.Context, offset, limit int) ([]*model.Alert, error) {
	fetched, err := getAlertQuery(x.client).
		Order(ent.Desc(entAlert.FieldCreatedAt)).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	alerts := make([]*model.Alert, len(fetched))
	for i := range fetched {
		alerts[i] = entToAlert(fetched[i])
	}
	return alerts, nil
}

func (x *Client) PutAlert(ctx *types.Context, alert *model.Alert) error {
	if err := alert.Validate(); err != nil {
		return err
	}

	q := x.client.Alert.Create().
		SetID(alert.ID()).
		SetTitle(alert.Title).
		SetDescription(alert.Description).
		SetDetector(alert.Detector).
		SetStatus(alert.Status).
		SetSeverity(alert.Severity).
		SetDetectedAt(alert.DetectedAt.Unix()).
		SetCreatedAt(alert.CreatedAt.Unix()).
		SetClosedAt(alert.ClosedAt.Unix())

	if _, err := q.Save(ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func (x *Client) UpdateAlert(ctx *types.Context, alert *model.Alert) error {
	q := x.client.Alert.UpdateOneID(alert.ID()).
		SetClosedAt(alert.ClosedAt.Unix()).
		SetSeverity(alert.Severity).
		SetStatus(alert.Status)

	if _, err := q.Save(ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func (x *Client) UpdateAlertSeverity(ctx *types.Context, id types.AlertID, sev types.Severity) error {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}
	if _, err := x.client.Alert.UpdateOneID(id).
		SetSeverity(sev).Save(ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func (x *Client) UpdateAlertStatus(ctx *types.Context, id types.AlertID, status types.AlertStatus) error {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}
	if _, err := x.client.Alert.UpdateOneID(id).
		SetStatus(status).Save(ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func (x *Client) UpdateAlertClosedAt(ctx *types.Context, id types.AlertID, ts int64) error {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}
	if _, err := x.client.Alert.UpdateOneID(id).
		SetClosedAt(ts).Save(ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil

}
