package db_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionLog(t *testing.T) {
	ctx := context.Background()
	t.Run("Add ActionLog", func(t *testing.T) {
		client := setupDB(t)
		created, _ := client.NewAlert(ctx)
		err := client.AddAttributes(ctx, created.ID, []*ent.Attribute{
			{
				Key:   "color",
				Value: "blue",
				Type:  types.AttrNoType,
			},
		})
		require.NoError(t, err)

		alert, _ := client.GetAlert(ctx, created.ID)
		log1, err := client.NewActionLog(ctx, created.ID, "blue", alert.Edges.Attributes[0].ID)
		require.NoError(t, err)

		require.NoError(t, client.AppendActionLog(ctx, log1.ID, &ent.ExecLog{
			Timestamp: 1001,
			Status:    types.ExecStart,
		}))
		require.NoError(t, client.AppendActionLog(ctx, log1.ID, &ent.ExecLog{
			Timestamp: 1003,
			Status:    types.ExecSucceed,
		}))

		got, err := client.GetAlert(ctx, alert.ID)
		require.NoError(t, err)
		require.Len(t, got.Edges.ActionLogs, 1)
		require.Len(t, got.Edges.ActionLogs[0].Edges.ExecLogs, 2)
		assert.Equal(t, int64(1003), got.Edges.ActionLogs[0].Edges.ExecLogs[0].Timestamp)
	})

	t.Run("AlertID not found", func(t *testing.T) {
		client := setupDB(t)
		_, err := client.NewActionLog(ctx, "xxx", "blue", 0)
		assert.ErrorIs(t, err, types.ErrDatabaseInvalidInput)
	})

	t.Run("ActionLogID not found", func(t *testing.T) {
		client := setupDB(t)
		err := client.AppendActionLog(ctx, 404, &ent.ExecLog{})
		assert.ErrorIs(t, err, types.ErrDatabaseInvalidInput)
	})
}
