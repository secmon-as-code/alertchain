package db

import (
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"

	entAlert "github.com/m-mizutani/alertchain/pkg/infra/ent/alert"
	"github.com/m-mizutani/alertchain/pkg/infra/ent/attribute"
)

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

func (x *Client) GetAlert(ctx *types.Context, id types.AlertID) (*ent.Alert, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	fetched, err := getAlertQuery(x.client).Where(entAlert.ID(id)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return fetched, nil
}

func (x *Client) GetAlerts(ctx *types.Context, offset, limit int) ([]*ent.Alert, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	fetched, err := getAlertQuery(x.client).Order(ent.Desc(entAlert.FieldCreatedAt)).Offset(offset).Limit(limit).All(ctx)
	if err != nil {
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return fetched, nil

}

func (x *Client) PutAlert(ctx *types.Context, alert *ent.Alert) (*ent.Alert, error) {
	q := x.client.Alert.Create().
		SetID(types.NewAlertID()).
		SetTitle(alert.Title).
		SetDescription(alert.Description).
		SetDetector(alert.Detector).
		SetStatus(alert.Status).
		SetSeverity(alert.Severity).
		SetDetectedAt(alert.DetectedAt).
		SetCreatedAt(alert.CreatedAt).
		SetClosedAt(alert.ClosedAt)

	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}
	created, err := q.Save(ctx)

	if err != nil {
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return created, nil
}

func (x *Client) UpdateAlert(ctx *types.Context, alert *ent.Alert) error {
	q := x.client.Alert.UpdateOneID(alert.ID).
		SetClosedAt(alert.ClosedAt).
		SetSeverity(alert.Severity).
		SetStatus(alert.Status)

	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}
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

func (x *Client) AddAttributes(ctx *types.Context, id types.AlertID, newAttrs []*ent.Attribute) error {
	if len(newAttrs) == 0 {
		return nil // nothing to do
	}

	builders := make([]*ent.AttributeCreate, len(newAttrs))
	for i, attr := range newAttrs {
		builders[i] = x.client.Attribute.Create().
			SetKey(attr.Key).
			SetValue(attr.Value).
			SetType(attr.Type).
			SetContext(attr.Context).
			SetAlertID(id)
	}

	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}
	added, err := x.client.Attribute.CreateBulk(builders...).Save(ctx)
	if err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	if _, err := x.client.Alert.UpdateOneID(id).AddAttributes(added...).Save(ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func (x *Client) GetAttribute(ctx *types.Context, id int) (*ent.Attribute, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	attr, err := x.client.Attribute.Query().Where(attribute.ID(id)).WithAlert().First(ctx)
	if err != nil {
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return attr, nil
}

func (x *Client) AddAnnotation(ctx *types.Context, attr *ent.Attribute, annotations []*ent.Annotation) error {
	if len(annotations) == 0 {
		return nil
	}

	builders := make([]*ent.AnnotationCreate, len(annotations))
	for i, ann := range annotations {
		builders[i] = x.client.Annotation.Create().
			SetName(ann.Name).
			SetSource(ann.Source).
			SetValue(ann.Value).
			SetTimestamp(ann.Timestamp)
	}

	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}
	added, err := x.client.Annotation.CreateBulk(builders...).Save(ctx)
	if err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	if _, err := x.client.Attribute.UpdateOneID(attr.ID).AddAnnotations(added...).Save(ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func (x *Client) AddReferences(ctx *types.Context, id types.AlertID, refs []*ent.Reference) error {
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
