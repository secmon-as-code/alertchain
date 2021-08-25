package db

import (
	"context"
	"time"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"

	entAlert "github.com/m-mizutani/alertchain/pkg/infra/ent/alert"
)

func (x *Client) GetAlert(ctx context.Context, id types.AlertID) (*ent.Alert, error) {
	fetched, err := x.client.Alert.Query().
		Where(entAlert.ID(id)).
		WithTaskLogs().
		WithReferences().
		WithAttributes(func(aq *ent.AttributeQuery) {
			aq.WithAnnotations()
		}).Only(x.ctx)
	if err != nil {
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return fetched, nil
}

func (x *Client) GetAlerts(ctx context.Context) ([]*ent.Alert, error) {
	fetched, err := x.client.Alert.Query().All(x.ctx)
	if err != nil {
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return fetched, nil
}

func (x *Client) NewAlert(ctx context.Context) (*ent.Alert, error) {
	newAlert, err := x.client.Alert.Create().
		SetID(types.NewAlertID()).
		SetCreatedAt(time.Now().UTC().Unix()).
		Save(x.ctx)
	if err != nil {
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return newAlert, nil
}

func (x *Client) UpdateAlert(ctx context.Context, id types.AlertID, alert *ent.Alert) error {
	q := x.client.Alert.UpdateOneID(id).
		SetTitle(alert.Title).
		SetDescription(alert.Description).
		SetDetector(alert.Detector)

	if alert.DetectedAt != nil {
		q = q.SetDetectedAt(*alert.DetectedAt)
	}
	if alert.ClosedAt != nil {
		q = q.SetClosedAt(*alert.ClosedAt)
	}

	if _, err := q.Save(x.ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}
	return nil
}

func (x *Client) UpdateAlertStatus(ctx context.Context, id types.AlertID, status types.AlertStatus, ts int64) error {
	q := x.client.Alert.UpdateOneID(id).SetStatus(status)
	if status == types.StatusClosed {
		q = q.SetClosedAt(ts)
	}
	if _, err := q.Save(x.ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}
	return nil
}

func (x *Client) UpdateAlertSeverity(ctx context.Context, id types.AlertID, sev types.Severity, ts int64) error {
	if _, err := x.client.Alert.UpdateOneID(id).SetSeverity(sev).Save(x.ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}
	return nil
}

func (x *Client) AddAttributes(ctx context.Context, id types.AlertID, newAttrs []*ent.Attribute) error {
	if len(newAttrs) == 0 {
		return nil // nothing to do
	}

	builders := make([]*ent.AttributeCreate, len(newAttrs))
	for i, attr := range newAttrs {
		builders[i] = x.client.Attribute.Create().
			SetKey(attr.Key).
			SetValue(attr.Value).
			SetType(attr.Type).
			SetContext(attr.Context)
	}
	added, err := x.client.Attribute.CreateBulk(builders...).Save(x.ctx)
	if err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	if _, err := x.client.Alert.UpdateOneID(id).AddAttributes(added...).Save(x.ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func (x *Client) AddAnnotation(ctx context.Context, attr *ent.Attribute, annotations []*ent.Annotation) error {
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

	added, err := x.client.Annotation.CreateBulk(builders...).Save(x.ctx)
	if err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	if _, err := attr.Update().AddAnnotations(added...).Save(x.ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func (x *Client) AddReference(ctx context.Context, id types.AlertID, ref *ent.Reference) error {
	if id == "" {
		return goerr.Wrap(types.ErrInvalidInput, "AlertID is not set")
	}
	if ref.Source == "" {
		return goerr.Wrap(types.ErrInvalidInput, "Reference.Source is not set")
	}
	if ref.Title == "" {
		return goerr.Wrap(types.ErrInvalidInput, "Reference.Title is not set")
	}
	if ref.URL == "" {
		return goerr.Wrap(types.ErrInvalidInput, "Reference.URL is not set")
	}

	added, err := x.client.Reference.Create().
		SetSource(ref.Source).
		SetTitle(ref.Title).
		SetURL(ref.URL).
		SetComment(ref.Comment).
		Save(ctx)
	if err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	if _, err := x.client.Alert.UpdateOneID(id).AddReferenceIDs(added.ID).Save(ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func (x *Client) NewTaskLog(ctx context.Context, id types.AlertID, taskName string, ts, stage int64, optional bool) (*ent.TaskLog, error) {
	if id == "" {
		return nil, goerr.Wrap(types.ErrInvalidInput, "AlertID is not set")
	}
	if taskName == "" {
		return nil, goerr.Wrap(types.ErrInvalidInput, "Reference.Source is not set")
	}

	taskLog, err := x.client.TaskLog.Create().
		SetTaskName(taskName).
		SetStage(stage).
		SetStartedAt(ts).
		SetOptional(optional).
		Save(ctx)
	if err != nil {
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	if _, err := x.client.Alert.UpdateOneID(id).AddTaskLogIDs(taskLog.ID).Save(ctx); err != nil {
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return taskLog, nil
}

func (x *Client) UpdateTaskLog(ctx context.Context, task *ent.TaskLog) error {
	if task.ID == 0 {
		return goerr.Wrap(types.ErrInvalidInput, "task.ID is not set")
	}

	q := x.client.TaskLog.UpdateOneID(task.ID).
		SetExitedAt(task.ExitedAt).
		SetLog(task.Log).
		SetErrmsg(task.Errmsg).
		SetStatus(task.Status)

	if _, err := q.Save(ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}
