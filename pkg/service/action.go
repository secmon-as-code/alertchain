package service

import (
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type ActionService struct {
	db interfaces.Database
}

func NewActionService(db interfaces.Database) *ActionService {
	return &ActionService{db: db}
}

func (x *ActionService) Get(ctx *model.Context, wfID types.WorkflowID) ([]model.ActionRecord, error) {
	return nil, nil
}
