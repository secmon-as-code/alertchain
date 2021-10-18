package alertchain

import "context"

type Handler func(ctx context.Context, alert *Alert) error
