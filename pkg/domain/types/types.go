package types

import (
	"github.com/google/uuid"
)

type (
	AlertID    string
	ActionName string
	ActionID   string

	Schema  string
	RawData any

	ScenarioID    string
	ScenarioTitle string

	EnvVarName  string
	EnvVarValue string

	ActionSecret any

	Namespace string

	WorkflowID string
)

// EnvVars is a set of environment variables
type EnvVars map[EnvVarName]EnvVarValue

func NewAlertID() AlertID   { return AlertID(uuid.New().String()) }
func NewActionID() ActionID { return ActionID(uuid.New().String()) }
func NewWorkflowID() WorkflowID {
	return WorkflowID(uuid.NewString())
}

func (x AlertID) String() string    { return string(x) }
func (x WorkflowID) String() string { return string(x) }
