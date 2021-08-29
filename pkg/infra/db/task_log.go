package db

import (
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

func (x *Client) NewTaskLog(ctx *types.Context, id types.AlertID, taskName string, stage int64) (*ent.TaskLog, error) {
	if id == "" {
		return nil, goerr.Wrap(types.ErrDatabaseInvalidInput, "AlertID is not set")
	}
	if taskName == "" {
		return nil, goerr.Wrap(types.ErrDatabaseInvalidInput, "Name is not set")
	}

	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}
	taskLog, err := x.client.TaskLog.Create().
		SetName(taskName).
		SetStage(stage).
		Save(ctx)
	if err != nil {
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	if _, err := x.client.Alert.UpdateOneID(id).AddTaskLogIDs(taskLog.ID).Save(ctx); err != nil {
		if ent.IsConstraintError(err) {
			return nil, types.ErrDatabaseInvalidInput.Wrap(err)
		}
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return taskLog, nil
}

func (x *Client) AppendTaskLog(ctx *types.Context, taskID int, execLog *ent.ExecLog) error {
	created, err := x.appendExecLog(ctx, execLog)
	if err != nil {
		return err
	}

	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}
	if _, err := x.client.TaskLog.UpdateOneID(taskID).AddExecLogs(created).Save(ctx); err != nil {
		if ent.IsConstraintError(err) {
			return types.ErrDatabaseInvalidInput.Wrap(err)
		}
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}
