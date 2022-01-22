package infra

import (
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
)

type Clients struct {
	db     *db.Client
	policy policy.Client
}

func (x *Clients) DB() *db.Client        { return x.db }
func (x *Clients) Policy() policy.Client { return x.policy }

func New(options ...Option) *Clients {
	clients := &Clients{}
	for _, opt := range options {
		opt(clients)
	}

	return clients
}
