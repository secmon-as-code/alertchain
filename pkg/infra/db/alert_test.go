package db_test

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlert(t *testing.T) {
	ctx := types.NewContext()
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

func TestReference(t *testing.T) {
	t.Run("Add Reference", func(t *testing.T) {
		client := setupDB(t)
		alert, _ := client.NewAlert(types.NewContext())

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

		require.NoError(t, client.AddReference(types.NewContext(), alert.ID, ref1))
		require.NoError(t, client.AddReference(types.NewContext(), alert.ID, ref2))
		require.NoError(t, client.AddReference(types.NewContext(), alert.ID, ref6))

		t.Run("get added references", func(t *testing.T) {
			got, err := client.GetAlert(types.NewContext(), alert.ID)
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
