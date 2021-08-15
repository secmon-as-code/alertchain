package alertchain

type Task interface {
	Name() string
	Description() string
	Execute(alert *Alert) error
	IsExecutable(alert *Alert) bool
}

type Tasks []Task
