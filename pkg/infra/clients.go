package infra

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/infra/db"
)

type Clients struct {
	DB *db.Client
}

func NewMock(t *testing.T) *Clients {
	return &Clients{
		DB: db.NewDBMock(t),
	}
}
