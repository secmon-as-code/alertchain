package usecase

import (
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
)

type Usecase struct {
	cfg *model.Config
}

func New(cfg *model.Config) interfaces.Usecase {
	return &Usecase{cfg: cfg}
}
