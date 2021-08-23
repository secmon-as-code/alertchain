package db_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func int64ptr(v int64) *int64 {
	return &v
}

func TestAlert(t *testing.T) {
	ctx := context.Background()
	t.Run("Create a new alert", func(t *testing.T) {
		client := setupDB(t)
		alert, err := client.NewAlert(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, alert.ID)
		assert.NotEmpty(t, alert.CreatedAt)
		assert.Empty(t, alert.Title)

		t.Run("Get alert by ID", func(t *testing.T) {
			got, err := client.GetAlert(ctx, alert.ID)
			require.NoError(t, err)
			assert.Equal(t, got.ID, alert.ID)
			assert.Empty(t, alert.Title)
		})

		t.Run("Update alert", func(t *testing.T) {
			alert.Title = "five"
			alert.Detector = "blue"
			alert.Description = "insane"

			// Should not be updated
			alert.Status = types.StatusClosed
			alert.Severity = types.SevAffected

			require.NoError(t, client.UpdateAlert(ctx, alert.ID, alert))

			t.Run("Get updated alert", func(t *testing.T) {
				got, err := client.GetAlert(ctx, alert.ID)
				require.NoError(t, err)
				assert.Equal(t, got.ID, alert.ID)
				assert.Equal(t, "five", got.Title)
				assert.Equal(t, "blue", got.Detector)
				assert.Equal(t, "insane", got.Description)
				assert.NotEqual(t, types.StatusClosed, got.Status)
				assert.NotEqual(t, types.SevAffected, got.Severity)

				t.Run("status can not be updated via SaveAlert", func(t *testing.T) {
					assert.NotEqual(t, types.StatusClosed, got.Status)
				})
			})
		})
	})

	t.Run("Create a new alert with attributes", func(t *testing.T) {
		client := setupDB(t)
		alert, _ := client.NewAlert(ctx)
		alert.Title = "five"
		attrs := []*ent.Attribute{
			{
				Key:     "srcaddr",
				Value:   "10.1.2.3",
				Type:    types.AttrIPAddr,
				Context: []string{string(types.CtxRemote)},
			},
			{
				Key:     "fqdn",
				Value:   "example.com",
				Type:    types.AttrDomain,
				Context: []string{string(types.CtxLocal)},
			},
		}
		require.NoError(t, client.UpdateAlert(ctx, alert.ID, alert))
		require.NoError(t, client.AddAttributes(ctx, alert.ID, attrs))

		t.Run("Get alert with attributes", func(t *testing.T) {
			got, err := client.GetAlert(ctx, alert.ID)
			require.NoError(t, err)
			assert.Len(t, got.Edges.Attributes, 2)
			equalAttributes(t, got.Edges.Attributes[0], attrs[0])
			equalAttributes(t, got.Edges.Attributes[1], attrs[1])
		})
	})

}
