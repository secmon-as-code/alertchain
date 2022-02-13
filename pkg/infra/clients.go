package infra

import (
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
)

type Clients struct {
	db     db.Client
	policy policy.Client
}

func (x *Clients) DB() db.Client         { return x.db }
func (x *Clients) Policy() policy.Client { return x.policy }

func New(dbClient db.Client, policyClient policy.Client) *Clients {
	clients := &Clients{
		db:     dbClient,
		policy: policyClient,
	}

	return clients
}
