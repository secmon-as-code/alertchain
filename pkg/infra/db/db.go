package db

import (
	"context"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/interfaces"
	"github.com/m-mizutani/alertchain/types"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	mock   bool
	ctx    context.Context
	client *ent.Client
}

func NewClient() interfaces.DBClient {
	return &Client{
		ctx: context.Background(),
	}
}

func NewDBMock() interfaces.DBClient {
	db := NewClient().(*Client)
	db.mock = true
	return db
}

func (x *Client) Init(dbType, dbConfig string) error {
	if x.mock {
		dbType = "sqlite3"
		dbConfig = "file:ent?mode=memory&cache=shared&_fk=1"
	}

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
