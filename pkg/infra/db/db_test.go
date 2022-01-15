package db_test

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) *db.Client {
	client := db.NewDBMock(t)
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Logf("Warning failed to close DB: %+v", err)
		}
	})

	return client
}

func equalAttributes(t *testing.T, expected, actual *model.Attribute) {
	assert.Equal(t, expected.Key, actual.Key)
	assert.Equal(t, expected.Value, actual.Value)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.Contexts, actual.Contexts)
}
