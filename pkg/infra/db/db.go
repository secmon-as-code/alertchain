package db

import (
	"context"
	"sync"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/infra/ent/enttest"
	"github.com/m-mizutani/alertchain/types"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type Interface interface {
	Close() error

	GetAlert(ctx context.Context, id types.AlertID) (*ent.Alert, error)
	GetAlerts(ctx context.Context) ([]*ent.Alert, error)
	NewAlert(ctx context.Context) (*ent.Alert, error)
	UpdateAlert(ctx context.Context, id types.AlertID, alert *ent.Alert) error
	UpdateAlertStatus(ctx context.Context, id types.AlertID, status types.AlertStatus, ts int64) error
	UpdateAlertSeverity(ctx context.Context, id types.AlertID, status types.Severity, ts int64) error

	AddAttributes(ctx context.Context, id types.AlertID, newAttrs []*ent.Attribute) error
	GetAttribute(ctx context.Context, id int) (*ent.Attribute, error)

	AddAnnotation(ctx context.Context, attr *ent.Attribute, ann []*ent.Annotation) error
	AddReference(ctx context.Context, id types.AlertID, ref *ent.Reference) error

	NewTaskLog(ctx context.Context, id types.AlertID, taskName string, stage int64) (*ent.TaskLog, error)
	AppendTaskLog(ctx context.Context, taskID int, execLog *ent.ExecLog) error

	NewActionLog(ctx context.Context, id types.AlertID, name string, attrID int) (*ent.ActionLog, error)
	AppendActionLog(ctx context.Context, actionID int, execLog *ent.ExecLog) error
}

type Client struct {
	ctx    context.Context
	client *ent.Client

	lock  bool
	mutex sync.Mutex
}

func newClient() *Client {
	return &Client{
		ctx: context.Background(),
	}
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
	db.client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
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
