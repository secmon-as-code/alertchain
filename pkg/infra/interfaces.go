package infra

import (
	"github.com/m-mizutani/alertchain/pkg/infra/db"
)

type Clients struct {
	DB db.Interface
}
