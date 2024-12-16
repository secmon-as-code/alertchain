package service

import (
	"github.com/secmon-lab/alertchain/pkg/domain/interfaces"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

type ActionService struct {
	db interfaces.Database
}

func NewActionService(db interfaces.Database) *ActionService {
	return &ActionService{db: db}
}

func (x *ActionService) Fetch(ctx *model.Context, wfID types.WorkflowID) ([]model.ActionRecord, error) {
	return nil, nil
}
