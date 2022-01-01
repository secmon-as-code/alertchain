package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/zlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newContext() *types.Context {
	var logger = zlog.New()
	return types.NewContextWith(context.Background(), logger)
}

func TestAlert(t *testing.T) {
	ctx := newContext()
	t.Run("Create a new alert", func(t *testing.T) {
		client := setupDB(t)
		alert := model.NewAlert(&model.Alert{
			Title:       "five",
			Detector:    "blue",
			Description: "insane",
			CreatedAt:   time.Now(),
		})

		require.NoError(t, client.PutAlert(ctx, alert))
		assert.NotEmpty(t, alert.ID)
		assert.Equal(t, "five", alert.Title)
		assert.Equal(t, "blue", alert.Detector)
		assert.Equal(t, "insane", alert.Description)

		t.Run("Get alert by ID", func(t *testing.T) {
			got, err := client.GetAlert(ctx, alert.ID())
			require.NoError(t, err)
			assert.Equal(t, got.ID(), alert.ID())
		})

		t.Run("Update alert", func(t *testing.T) {
			assert.NotEqual(t, types.StatusClosed, alert.Status)
			assert.NotEqual(t, types.SevAffected, alert.Severity)

			// Should be updated
			alert.Status = types.StatusClosed
			alert.Severity = types.SevAffected

			require.NoError(t, client.UpdateAlertSeverity(ctx, alert.ID(), types.SevAffected))
			require.NoError(t, client.UpdateAlertStatus(ctx, alert.ID(), types.StatusClosed))
			require.NoError(t, client.UpdateAlertClosedAt(ctx, alert.ID(), 1234))

			t.Run("Get updated alert", func(t *testing.T) {
				got, err := client.GetAlert(ctx, alert.ID())
				require.NoError(t, err)
				assert.Equal(t, got.ID(), alert.ID())
				assert.Equal(t, types.StatusClosed, alert.Status)
				assert.Equal(t, types.SevAffected, alert.Severity)
				assert.Equal(t, int64(1234), got.ClosedAt.Unix())
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

		require.NoError(t, client.PutAlert(ctx, alert))

		attrs := []*model.Attribute{
			{
				Key:      "srcaddr",
				Value:    "10.1.2.3",
				Type:     types.AttrIPAddr,
				Contexts: model.Contexts{types.CtxRemote},
			},
			{
				Key:      "fqdn",
				Value:    "example.com",
				Type:     types.AttrDomain,
				Contexts: model.Contexts{types.CtxLocal},
			},
		}

		require.NoError(t, client.AddAttributes(ctx, alert.ID(), attrs))

		t.Run("Get alert with attributes", func(t *testing.T) {
			got, err := client.GetAlert(ctx, alert.ID())
			require.NoError(t, err)
			assert.Len(t, got.Attributes, 2)
			equalAttributes(t, got.Attributes[0], attrs[0])
			equalAttributes(t, got.Attributes[1], attrs[1])
		})
	})
}

func TestReference(t *testing.T) {
	t.Run("Add Reference", func(t *testing.T) {
		client := setupDB(t)
		ctx := newContext()
		alert := model.NewAlert(&model.Alert{
			Title:     "five",
			Detector:  "blue",
			CreatedAt: time.Now(),
		})
		require.NoError(t, client.PutAlert(ctx, alert))

		ref1 := &model.Reference{
			Source:  "blue",
			URL:     "https://example.com/b1",
			Title:   "b1",
			Comment: "pity",
		}
		ref2 := &model.Reference{
			Source:  "blue",
			URL:     "https://example.com/b2",
			Title:   "b2",
			Comment: "regression",
		}
		ref6 := &model.Reference{
			Source: "orange",
			URL:    "https://example.com/b6",
			Title:  "b6",
		}

		require.NoError(t, client.AddReferences(ctx, alert.ID(), []*model.Reference{ref1, ref2}))
		require.NoError(t, client.AddReferences(ctx, alert.ID(), []*model.Reference{ref6}))

		t.Run("get added references", func(t *testing.T) {
			got, err := client.GetAlert(ctx, alert.ID())
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Len(t, got.References, 3)

			assert.Equal(t, "blue", got.References[0].Source)
			assert.Equal(t, "https://example.com/b1", got.References[0].URL)
			assert.Equal(t, "b1", got.References[0].Title)
			assert.Equal(t, "pity", got.References[0].Comment)

			assert.Equal(t, "blue", got.References[1].Source)
			assert.Equal(t, "https://example.com/b2", got.References[1].URL)
			assert.Equal(t, "b2", got.References[1].Title)
			assert.Equal(t, "regression", got.References[1].Comment)

			assert.Equal(t, "orange", got.References[2].Source)
			assert.Equal(t, "https://example.com/b6", got.References[2].URL)
			assert.Equal(t, "b6", got.References[2].Title)
			assert.Empty(t, got.References[2].Comment)
		})
	})
}
