package db_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
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
		alert, err := client.PutAlert(ctx, &ent.Alert{
			Title:       "five",
			Detector:    "blue",
			Description: "insane",
		})
		require.NoError(t, err)
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
			assert.NotEqual(t, types.StatusClosed, alert.Status)
			assert.NotEqual(t, types.SevAffected, alert.Severity)

			// Should be updated
			alert.Status = types.StatusClosed
			alert.Severity = types.SevAffected

			require.NoError(t, client.UpdateAlertSeverity(ctx, alert.ID, types.SevAffected))
			require.NoError(t, client.UpdateAlertStatus(ctx, alert.ID, types.StatusClosed))
			require.NoError(t, client.UpdateAlertClosedAt(ctx, alert.ID, 1234))

			t.Run("Get updated alert", func(t *testing.T) {
				got, err := client.GetAlert(ctx, alert.ID)
				require.NoError(t, err)
				assert.Equal(t, got.ID, alert.ID)
				assert.Equal(t, types.StatusClosed, alert.Status)
				assert.Equal(t, types.SevAffected, alert.Severity)
				assert.Equal(t, int64(1234), got.ClosedAt)
			})
		})
	})

	t.Run("Create a new alert with attributes", func(t *testing.T) {
		client := setupDB(t)
		alert, _ := client.PutAlert(ctx, &ent.Alert{
			Title: "five",
		})

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

func TestReference(t *testing.T) {
	t.Run("Add Reference", func(t *testing.T) {
		client := setupDB(t)
		ctx := newContext()
		alert, _ := client.PutAlert(ctx, &ent.Alert{})

		ref1 := &ent.Reference{
			Source:  "blue",
			URL:     "https://example.com/b1",
			Title:   "b1",
			Comment: "pity",
		}
		ref2 := &ent.Reference{
			Source:  "blue",
			URL:     "https://example.com/b2",
			Title:   "b2",
			Comment: "regression",
		}
		ref6 := &ent.Reference{
			Source: "orange",
			URL:    "https://example.com/b6",
			Title:  "b6",
		}

		require.NoError(t, client.AddReferences(ctx, alert.ID, []*ent.Reference{ref1, ref2}))
		require.NoError(t, client.AddReferences(ctx, alert.ID, []*ent.Reference{ref6}))

		t.Run("get added references", func(t *testing.T) {
			got, err := client.GetAlert(ctx, alert.ID)
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Len(t, got.Edges.References, 3)

			assert.Equal(t, "blue", got.Edges.References[0].Source)
			assert.Equal(t, "https://example.com/b1", got.Edges.References[0].URL)
			assert.Equal(t, "b1", got.Edges.References[0].Title)
			assert.Equal(t, "pity", got.Edges.References[0].Comment)

			assert.Equal(t, "blue", got.Edges.References[1].Source)
			assert.Equal(t, "https://example.com/b2", got.Edges.References[1].URL)
			assert.Equal(t, "b2", got.Edges.References[1].Title)
			assert.Equal(t, "regression", got.Edges.References[1].Comment)

			assert.Equal(t, "orange", got.Edges.References[2].Source)
			assert.Equal(t, "https://example.com/b6", got.Edges.References[2].URL)
			assert.Equal(t, "b6", got.Edges.References[2].Title)
			assert.Empty(t, got.Edges.References[2].Comment)
		})
	})
}
