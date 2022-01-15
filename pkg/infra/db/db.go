package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/m-mizutani/alertchain/gen/ent"
	"github.com/m-mizutani/alertchain/gen/ent/enttest"
	"github.com/m-mizutani/alertchain/pkg/domain/types"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	client *ent.Client
}

func newClient() *Client {
	return &Client{}
}

func New(dbType, dbConfig string) (*Client, error) {
	client := newClient()
	if err := client.init(dbType, dbConfig); err != nil {
		return nil, err
	}
	return client, nil
}

func NewDBMock(t *testing.T) *Client {
	db := newClient()
	db.client = enttest.Open(t, "sqlite3", "file:"+uuid.NewString()+"?mode=memory&cache=private&_fk=1")
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
