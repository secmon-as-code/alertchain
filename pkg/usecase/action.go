package usecase

import (
	"bytes"
	"time"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

func (x *usecase) GetExecutableActions(ctx *types.Context, attr *ent.Attribute) ([]*Action, error) {
	var resp []*Action
	for _, action := range x.actions {
		if action.Executable(attr) {
			resp = append(resp, action)
		}
	}
	return resp, nil
}

func (x *usecase) ExecuteAction(ctx *types.Context, actionID string, attrID int) (*ent.ActionLog, error) {
	action, ok := x.actions[actionID]
	if !ok {
		return nil, goerr.Wrap(types.ErrInvalidInput, "invalid action ID")
	}
	attr, err := x.clients.DB.GetAttribute(ctx, attrID)
	if err != nil {
		return nil, err
	}

	actionLog, err := x.clients.DB.NewActionLog(ctx, attr.Edges.Alert.ID, action.Name, attrID)
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
		logW := &bytes.Buffer{}
		ctx := types.NewContext().InjectWriter(logW)
		execLog := &ent.ExecLog{
			Status: types.ExecSucceed,
		}
		defer func() {
			execLog.Timestamp = time.Now().UnixNano()
			execLog.Log = logW.String()
			if err := x.clients.DB.AppendActionLog(ctx, actionLog.ID, execLog); err != nil {
				utils.HandleError(err)
			}
		}()

		if err := action.Execute(ctx, attr); err != nil {
			utils.CopyErrorToExecLog(err, execLog)
			utils.HandleError(err)
			return
		}
		execLog.Timestamp = time.Now().UTC().UnixNano()
	}()

	return actionLog, nil
}
