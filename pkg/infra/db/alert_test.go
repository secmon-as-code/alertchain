package db_test

import (
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlert(t *testing.T) {
	t.Run("Create a new alert", func(t *testing.T) {
		client := setupDB(t)
		alert, err := client.NewAlert()
		require.NoError(t, err)
		assert.NotEmpty(t, alert.ID)
		assert.NotEmpty(t, alert.CreatedAt)
		assert.Empty(t, alert.Title)

		t.Run("Get alert by ID", func(t *testing.T) {
			got, err := client.GetAlert(alert.ID)
			require.NoError(t, err)
			assert.Equal(t, got.ID, alert.ID)
			assert.Empty(t, alert.Title)
		})

		t.Run("Update alert", func(t *testing.T) {
			now := time.Now().UTC().Add(time.Hour)
			alert.Title = "five"
			alert.Detector = "blue"
			alert.Description = "insane"
			alert.Status = types.StatusClosed
			alert.Severity = types.SevAffected
			alert.ClosedAt = &now

			require.NoError(t, client.SaveAlert(alert))

			t.Run("Get updated alert", func(t *testing.T) {
				got, err := client.GetAlert(alert.ID)
				require.NoError(t, err)
				assert.Equal(t, got.ID, alert.ID)
				assert.Equal(t, "five", got.Title)
				assert.Equal(t, "blue", got.Detector)
				assert.Equal(t, "insane", got.Description)
				assert.Equal(t, types.StatusClosed, got.Status)
				assert.Equal(t, types.SevAffected, got.Severity)
				assert.Equal(t, now, *got.ClosedAt)
			})
		})
	})

	t.Run("Create a new alert with attributes", func(t *testing.T) {
		client := setupDB(t)
		alert, _ := client.NewAlert()
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
		require.NoError(t, client.SaveAlert(alert))
		require.NoError(t, client.AddAttributes(alert, attrs))

		t.Run("Get alert with attributes", func(t *testing.T) {
			got, err := client.GetAlert(alert.ID)
			require.NoError(t, err)
			assert.Len(t, got.Edges.Attributes, 2)
			equalAttributes(t, got.Edges.Attributes[0], attrs[0])
			equalAttributes(t, got.Edges.Attributes[1], attrs[1])
		})
	})

}
