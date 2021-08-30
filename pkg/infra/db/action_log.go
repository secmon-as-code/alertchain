package db

import (
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/infra/ent/actionlog"
	"github.com/m-mizutani/alertchain/pkg/infra/ent/attribute"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

func (x *Client) NewActionLog(ctx *types.Context, name string, attrID int) (*ent.ActionLog, error) {
	if name == "" {
		return nil, goerr.Wrap(types.ErrDatabaseInvalidInput, "Name is not set")
	}

	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	got, err := x.client.Attribute.Query().
		Where(attribute.ID(attrID)).WithAlert().Only(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, types.ErrDatabaseInvalidInput.Wrap(err)
		}
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	actionLog, err := x.client.ActionLog.Create().
		SetName(name).
		AddArgumentIDs(attrID).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, types.ErrDatabaseInvalidInput.Wrap(err)
		}
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	if _, err := x.client.Alert.UpdateOneID(got.Edges.Alert.ID).AddActionLogIDs(actionLog.ID).Save(ctx); err != nil {
		if ent.IsConstraintError(err) {
			return nil, types.ErrDatabaseInvalidInput.Wrap(err)
		}
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return actionLog, nil
}

func (x *Client) AppendActionLog(ctx *types.Context, actionID int, execLog *ent.ExecLog) error {
	created, err := x.appendExecLog(ctx, execLog)
	if err != nil {
		return err
	}

	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	if _, err := x.client.ActionLog.UpdateOneID(actionID).AddExecLogs(created).Save(ctx); err != nil {
		if ent.IsConstraintError(err) {
			return types.ErrDatabaseInvalidInput.Wrap(err)
		}
		return types.ErrDatabaseUnexpected.Wrap(err)
	}

	return nil
}

func (x *Client) appendExecLog(ctx *types.Context, execLog *ent.ExecLog) (*ent.ExecLog, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}
	created, err := x.client.ExecLog.Create().
		SetTimestamp(execLog.Timestamp).
		SetLog(execLog.Log).
		SetErrmsg(execLog.Errmsg).
		SetErrValues(execLog.ErrValues).
		SetStackTrace(execLog.StackTrace).
		SetStatus(execLog.Status).
		Save(ctx)
	if err != nil {
		return nil, types.ErrDatabaseUnexpected.Wrap(err)

	}
	return created, nil
}

func (x *Client) GetActionLog(ctx *types.Context, actionLogID int) (*ent.ActionLog, error) {
	if x.lock {
		x.mutex.Lock()
		defer x.mutex.Unlock()
	}

	got, err := x.client.ActionLog.Query().
		Where(actionlog.ID(actionLogID)).
		WithExecLogs().Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, types.ErrItemNotFound.Wrap(err)
		}
		return nil, types.ErrDatabaseUnexpected.Wrap(err)
	}

	return got, nil
}
