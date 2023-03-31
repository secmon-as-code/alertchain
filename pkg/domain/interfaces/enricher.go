package interfaces

import (
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type Enricher interface {
	Enrich(ctx *types.Context, id types.Parameter) ([]model.Reference, error)
}
