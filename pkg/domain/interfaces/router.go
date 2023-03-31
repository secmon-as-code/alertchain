package interfaces

import (
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type Router func(ctx *types.Context, label string, data any) error
