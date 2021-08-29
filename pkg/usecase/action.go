package usecase

import (
	"context"
	"time"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
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
	action, ok := x.actions[actionID]
	if !ok {
		return nil, goerr.Wrap(types.ErrInvalidInput, "invalid action ID")
	}
	attr, err := x.clients.DB.GetAttribute(ctx, attrID)
	if err != nil {
		return nil, err
	}

	actionLog, err := x.clients.DB.NewActionLog(ctx, attr.Edges.Alert.ID, action.Action.Name(), attrID)
	if err != nil {
		return nil, err
	}
	if err := x.clients.DB.AppendActionLog(ctx, actionLog.ID, &ent.ExecLog{
		Timestamp: time.Now().UTC().UnixNano(),
		Status:    types.ExecStart,
	}); err != nil {
		return nil, err
	}

	go func() {
		ctx := context.Background()
		execLog := &ent.ExecLog{
			Status: types.ExecSucceed,
		}
		defer func() {
			execLog.Timestamp = time.Now().UnixNano()
			// execLog.Log =
			if err := x.clients.DB.AppendActionLog(ctx, actionLog.ID, execLog); err != nil {
				utils.HandleError(err)
			}
		}()
		arg := &alertchain.Attribute{Attribute: *attr}
		if err := action.Action.Execute(ctx, arg); err != nil {
			utils.CopyErrorToExecLog(err, execLog)
			utils.HandleError(err)
			return
		}
		execLog.Timestamp = time.Now().UTC().UnixNano()
	}()

	return &alertchain.ActionLog{ActionLog: *actionLog}, nil
}
