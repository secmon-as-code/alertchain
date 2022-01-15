package infra

import (
	"github.com/m-mizutani/alertchain/pkg/infra/db"
)

type Clients struct {
	DB *db.Client
}

func New(options ...Option) (*Clients, error) {
	clients := &Clients{}
	for _, opt := range options {
		if err := opt(clients); err != nil {
			return nil, err
		}
	}

	return clients, nil
}
