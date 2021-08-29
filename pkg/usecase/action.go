package usecase

import (
	"context"

	"github.com/m-mizutani/alertchain"
)

func (x *usecase) GetExecutableActions(ctx context.Context, attr *alertchain.Attribute) ([]*alertchain.ActionEntry, error) {
	var resp []*alertchain.ActionEntry
	for _, entry := range x.actions {
		if entry.Action.Executable(attr) {
			resp = append(resp, entry)
		}
	}
	return resp, nil
}

func (x *usecase) ExecuteAction(ctx context.Context, actionID string, attrID int) (*alertchain.ActionLog, error) {
	/*
		attr, err := x.clients.DB.GetAttribute(ctx, attrID)
		if err != nil {
			return nil, err
		}
	*/

	return nil, nil
}
