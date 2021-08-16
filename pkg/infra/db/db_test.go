package db_test

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) infra.DBClient {
	client := db.NewDBMock(t)
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
