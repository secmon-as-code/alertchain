package db

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/infra/ent/enttest"
)

func NewTestClient(t *testing.T) *Client {
	client := newClient()
	client.client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")

	return client
}
