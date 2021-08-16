package db

import (
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"

	entAlert "github.com/m-mizutani/alertchain/pkg/infra/ent/alert"
)

func (x *Client) GetAlert(id types.AlertID) (*ent.Alert, error) {
	fetched, err := x.client.Alert.Query().
		Where(entAlert.ID(id)).WithAttributes(func(aq *ent.AttributeQuery) {
		aq.WithFindings()
	}).Only(x.ctx)
	if err != nil {
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return fetched, nil
}

func (x *Client) GetAlerts() ([]*ent.Alert, error) {
	fetched, err := x.client.Alert.Query().All(x.ctx)
	if err != nil {
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return fetched, nil
}

func (x *Client) NewAlert() (*ent.Alert, error) {
	newAlert, err := x.client.Alert.Create().SetID(types.NewAlertID()).Save(x.ctx)
	if err != nil {
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return newAlert, nil
}

func (x *Client) SaveAlert(alert *ent.Alert) error {
	q := alert.Update().
		SetTitle(alert.Title).
		SetDescription(alert.Description).
		SetDetector(alert.Detector).
		SetStatus(alert.Status).
		SetSeverity(alert.Severity)

	if alert.ClosedAt != nil {
		q.SetClosedAt(*alert.ClosedAt)
	}

	if _, err := q.Save(x.ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}
	return nil
}

func (x *Client) AddAttributes(alert *ent.Alert, newAttrs []*ent.Attribute) error {
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

	if _, err := alert.Update().AddAttributes(added...).Save(x.ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func (x *Client) AddFindings(attr *ent.Attribute, findings []*ent.Finding) error {
	if len(findings) == 0 {
		return nil
	}

	builders := make([]*ent.FindingCreate, len(findings))
	for i, finding := range findings {
		builders[i] = x.client.Finding.Create().
			SetSource(finding.Source).
			SetValue(finding.Value).
			SetTime(finding.Time)
	}

	added, err := x.client.Finding.CreateBulk(builders...).Save(x.ctx)
	if err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	if _, err := attr.Update().AddFindings(added...).Save(x.ctx); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}
