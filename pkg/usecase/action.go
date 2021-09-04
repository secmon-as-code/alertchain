package usecase

import (
	"bytes"
	"time"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

func (x *usecase) GetExecutableActions(attr *ent.Attribute) []*Action {
	var resp []*Action
	for _, action := range x.actions {
		if action.Executable(attr) {
			resp = append(resp, action)
		}
	}
	return resp
}

func (x *usecase) GetActionLog(ctx *types.Context, actionLogID int) (*ent.ActionLog, error) {
	return x.clients.DB.GetActionLog(ctx, actionLogID)
}

func (x *usecase) ExecuteAction(ctx *types.Context, actionID string, attrID int) (*ent.ActionLog, error) {
	action, ok := x.actions[actionID]
	if !ok {
		return nil, goerr.Wrap(types.ErrInvalidInput, "no such action ID")
	}
	attr, err := x.clients.DB.GetAttribute(ctx, attrID)
	if err != nil {
		return nil, err
	}

	actionLog, err := x.clients.DB.NewActionLog(ctx, action.Name, attrID)
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
