package db

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/infra/ent/enttest"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Interface interface {
	Close() error

	GetAlert(ctx *types.Context, id types.AlertID) (*model.Alert, error)
	GetAlerts(ctx *types.Context, offset, limit int) ([]*model.Alert, error)
	PutAlert(ctx *types.Context, alert *model.Alert) error
	UpdateAlertStatus(ctx *types.Context, id types.AlertID, status types.AlertStatus) error
	UpdateAlertSeverity(ctx *types.Context, id types.AlertID, sev types.Severity) error
	UpdateAlertClosedAt(ctx *types.Context, id types.AlertID, ts int64) error

	AddAttributes(ctx *types.Context, id types.AlertID, newAttrs []*model.Attribute) error
	AddReferences(ctx *types.Context, id types.AlertID, refs []*model.Reference) error
	AddAnnotation(ctx *types.Context, attr *model.Attribute, ann []*model.Annotation) error
}

type Client struct {
	client *ent.Client

	lock  bool
	mutex sync.Mutex
}

func newClient() *Client {
	return &Client{}
}

func New(dbType, dbConfig string) (Interface, error) {
	client := newClient()
	if err := client.init(dbType, dbConfig); err != nil {
		return nil, err
	}
	return client, nil
}

func NewDBMock(t *testing.T) Interface {
	db := newClient()
	db.client = enttest.Open(t, "sqlite3", "file:"+uuid.NewString()+"?mode=memory&cache=shared&_fk=1")
	db.lock = true
	return db
}

func (x *Client) init(dbType, dbConfig string) error {
	client, err := ent.Open(dbType, dbConfig)
	if err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}
	x.client = client

	if err := client.Schema.Create(context.Background()); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func (x *Client) Close() error {
	if err := x.client.Close(); err != nil {
		return types.ErrDatabaseUnexpected.Wrap(err)
	}
	return nil
}
