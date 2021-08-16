package usecase

import (
	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra"
)

type usecase struct {
	infra infra.Infra
	chain *alertchain.Chain
}

func New(infra infra.Infra, chain *alertchain.Chain) Usecase {
	return &usecase{
		infra: infra,
		chain: chain,
	}
}
