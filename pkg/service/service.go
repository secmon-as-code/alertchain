package service

import "github.com/m-mizutani/alertchain/pkg/infra"

type Service struct {
	clients *infra.Clients
}

func New(clients *infra.Clients) *Service {
	return &Service{
		clients: clients,
	}
}
