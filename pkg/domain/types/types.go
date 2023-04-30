package types

import "github.com/google/uuid"

type (
	AlertID    string
	ProbeID    string
	ProbeName  string
	ActionID   string
	ActionName string
	EnricherID string

	Schema string

	ScenarioID    string
	ScenarioTitle string
)

func NewAlertID() AlertID {
	return AlertID(uuid.New().String())
}
