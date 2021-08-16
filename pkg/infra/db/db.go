package db

import (
	"context"

	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	mock   bool
	ctx    context.Context
	client *ent.Client
}

func newClient() *Client {
	return &Client{
		ctx: context.Background(),
	}
}

func New(dbType, dbConfig string) (infra.DBClient, error) {
	client := newClient()
	if err := client.init(dbType, dbConfig); err != nil {
		return nil, err
	}
	return client, nil
}

func NewDBMock() infra.DBClient {
	db := newClient()

	if err := db.init("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1"); err != nil {
		panic(err.Error())
	}
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
