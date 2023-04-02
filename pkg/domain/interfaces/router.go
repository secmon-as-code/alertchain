package interfaces

import (
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type Router func(ctx *types.Context, schema types.Schema, data any) error
