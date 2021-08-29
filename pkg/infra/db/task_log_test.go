package db_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskLog(t *testing.T) {
	ctx := context.Background()
	t.Run("Add TaskLog", func(t *testing.T) {
		client := setupDB(t)
		alert, _ := client.NewAlert(context.Background())

		t1, err := client.NewTaskLog(ctx, alert.ID, "blue", 0)
		require.NoError(t, err)
		t2, err := client.NewTaskLog(ctx, alert.ID, "orange", 1)
		require.NoError(t, err)
		t3, err := client.NewTaskLog(ctx, alert.ID, "red", 2)
		require.NoError(t, err)

		require.NoError(t, client.AppendTaskLog(ctx, t1.ID, &ent.ExecLog{
			Timestamp: 1001,
			Status:    types.ExecStart,
		}))
		require.NoError(t, client.AppendTaskLog(ctx, t1.ID, &ent.ExecLog{
			Timestamp: 1010,
			Log:       "timeless",
			Status:    types.ExecSucceed,
		}))
		require.NoError(t, client.AppendTaskLog(ctx, t2.ID, &ent.ExecLog{
			Timestamp: 1002,
			Log:       "rune",
			Status:    types.ExecFailure,
		}))
		require.NoError(t, client.AppendTaskLog(ctx, t3.ID, &ent.ExecLog{
			Timestamp: 1003,
			Log:       "scar",
			Errmsg:    "x",
			Status:    types.ExecFailure,
		}))

		got, err := client.GetAlert(ctx, alert.ID)
		require.NoError(t, err)
		require.Len(t, got.Edges.TaskLogs, 3)
		require.Len(t, got.Edges.TaskLogs[0].Edges.ExecLogs, 2)
		assert.Equal(t, "timeless", got.Edges.TaskLogs[0].Edges.ExecLogs[0].Log)

		assert.Equal(t, "rune", got.Edges.TaskLogs[1].Edges.ExecLogs[0].Log)
		assert.Equal(t, types.ExecFailure, got.Edges.TaskLogs[1].Edges.ExecLogs[0].Status)

		assert.Equal(t, "scar", got.Edges.TaskLogs[2].Edges.ExecLogs[0].Log)
		assert.Equal(t, "x", got.Edges.TaskLogs[2].Edges.ExecLogs[0].Errmsg)
	})

	t.Run("AlertID not found", func(t *testing.T) {
		client := setupDB(t)
		_, err := client.NewTaskLog(ctx, "xxx", "blue", 0)
		assert.ErrorIs(t, err, types.ErrDatabaseInvalidInput)
	})

	t.Run("TaskID not found", func(t *testing.T) {
		client := setupDB(t)
		err := client.AppendTaskLog(ctx, 1234, &ent.ExecLog{})
		assert.ErrorIs(t, err, types.ErrDatabaseInvalidInput)
	})
}
