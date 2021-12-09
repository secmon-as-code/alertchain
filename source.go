package alertchain

import (
	"context"
)

type Handler func(ctx context.Context, alert *Alert) error

type Source interface {
	Name() string
	Run(handler Handler) error
}
