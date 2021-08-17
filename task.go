package alertchain

import "context"

type Task interface {
	Name() string
	Description() string
	Execute(ctx context.Context, alert *Alert) error
	IsExecutable(alert *Alert) bool
}
