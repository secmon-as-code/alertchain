package usecase

import (
	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra"
)

type usecase struct {
	clients infra.Clients
	chain   *alertchain.Chain
}

func New(clients infra.Clients, chain *alertchain.Chain) Usecase {
	return &usecase{
		clients: clients,
		chain:   chain,
	}
}
