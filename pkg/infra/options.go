package infra

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/infra/db"
)

type Option func(*Clients) error

func WithDB(dbType, dbConfig string) Option {
	return func(c *Clients) error {
		client, err := db.New(dbType, dbConfig)
		if err != nil {
			return err
		}
		c.DB = client
		return nil
	}
}

func WithDBMock(t *testing.T) Option {
	return func(c *Clients) error {
		c.DB = db.NewDBMock(t)
		return nil
	}
}
