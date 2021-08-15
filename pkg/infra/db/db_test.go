package db_test

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/infra/ent/enttest"
	"github.com/m-mizutani/alertchain/pkg/interfaces"
	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) interfaces.DBClient {
	entClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")

	client := db.NewClient().(*db.Client)
	client.InjectClient(entClient)

	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning failed to close DB: %+v", err)
		}
	})

	return client
}

func equalAttributes(t *testing.T, expected, actual *ent.Attribute) {
	assert.Equal(t, expected.Key, actual.Key)
	assert.Equal(t, expected.Value, actual.Value)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.Context, actual.Context)
}
