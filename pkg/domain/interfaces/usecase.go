package interfaces

import "github.com/m-mizutani/alertchain/pkg/domain/model"

type Usecase interface {
}

type NewUsecase func(cfg *model.Config) Usecase
