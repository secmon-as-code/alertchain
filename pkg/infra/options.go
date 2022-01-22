package infra

import (
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
)

type Option func(*Clients)

func WithDB(client *db.Client) Option {
	return func(c *Clients) {
		c.db = client
	}
}

func WithPolicy(client policy.Client) Option {
	return func(c *Clients) {
		c.policy = client
	}
}
