package spanner_test

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/infra/db/spanner"
)

func setupDB(t *testing.T) *spanner.Client {
	client := spanner.NewTestDB(t)
	return client
}
