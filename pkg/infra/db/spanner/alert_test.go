package spanner_test

import (
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlert(t *testing.T) {
	ctx := types.NewContext()
	t.Run("Create a new alert", func(t *testing.T) {
		client := setupDB(t)
		alert := model.NewAlert(&model.Alert{
			Title:       "five",
			Detector:    "blue",
			Description: "insane",
			CreatedAt:   time.Now(),
		})

		require.NoError(t, client.SaveAlert(ctx, alert))
		assert.NotEmpty(t, alert.ID)
		assert.Equal(t, "five", alert.Title)
		assert.Equal(t, "blue", alert.Detector)
		assert.Equal(t, "insane", alert.Description)

		t.Run("Get alert by ID", func(t *testing.T) {
			got, err := client.GetAlert(ctx, alert.ID)
			require.NoError(t, err)
			assert.Equal(t, got.ID, alert.ID)
		})

		t.Run("Update alert", func(t *testing.T) {
			assert.NotEqual(t, types.SevAffected, alert.Severity)

			// Should be updated
			alert.Severity = types.SevAffected
			now := time.Now()
			require.NoError(t, client.UpdateAlertSeverity(ctx, alert.ID, types.SevAffected))
			require.NoError(t, client.UpdateAlertClosedAt(ctx, alert.ID, now))

			t.Run("Get updated alert", func(t *testing.T) {
				got, err := client.GetAlert(ctx, alert.ID)
				require.NoError(t, err)
				assert.Equal(t, got.ID, alert.ID)
				assert.Equal(t, types.SevAffected, alert.Severity)
				assert.Equal(t, now.Unix(), got.ClosedAt.Unix())
			})
		})
	})

	t.Run("Create a new alert with attributes", func(t *testing.T) {
		client := setupDB(t)
		alert := model.NewAlert(&model.Alert{
			Title:     "five",
			Detector:  "blue",
			CreatedAt: time.Now(),
		})

		require.NoError(t, client.SaveAlert(ctx, alert))

		attrs := []*model.Attribute{
			alert.NewAttribute("srcaddr", "10.1.2.3", types.AttrIPAddr),
			alert.NewAttribute("fqdn", "example.com", types.AttrDomain, types.CtxLocal),
		}

		require.NoError(t, client.AddAttributes(ctx, attrs))

		t.Run("Get alert with attributes", func(t *testing.T) {
			got, err := client.GetAlert(ctx, alert.ID)
			require.NoError(t, err)
			assert.Len(t, got.Attributes, 2)
			assert.Contains(t, got.Attributes, attrs[0])
			assert.Contains(t, got.Attributes, attrs[1])
		})
	})
}

func TestReference(t *testing.T) {
	t.Run("Add Reference", func(t *testing.T) {
		client := setupDB(t)
		ctx := types.NewContext()
		alert := model.NewAlert(&model.Alert{
			Title:     "five",
			Detector:  "blue",
			CreatedAt: time.Now(),
		})
		require.NoError(t, client.SaveAlert(ctx, alert))

		ref1 := alert.NewReference("blue", "b1", "https://example.com/b1", "pity")
		ref2 := alert.NewReference("blue", "b2", "https://example.com/b2", "regression")
		ref6 := alert.NewReference("orange", "b6", "https://example.com/b6", "")

		require.NoError(t, client.AddReferences(ctx, []*model.Reference{ref1, ref2}))
		require.NoError(t, client.AddReferences(ctx, []*model.Reference{ref6}))

		got, err := client.GetAlert(ctx, alert.ID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Len(t, got.References, 3)

		assert.Contains(t, got.References, ref1)
		assert.Contains(t, got.References, ref2)
		assert.Contains(t, got.References, ref6)
	})
}
