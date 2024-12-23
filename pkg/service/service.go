package service

import "github.com/secmon-lab/alertchain/pkg/domain/interfaces"

type Services struct {
	Workflow *WorkflowService
}

func New(db interfaces.Database) *Services {
	return &Services{
		Workflow: NewWorkflowService(db),
	}
}
