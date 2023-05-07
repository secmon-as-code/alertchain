package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type ScenarioLog struct {
	ID    types.ScenarioID    `json:"id"`
	Title types.ScenarioTitle `json:"title"`

	Results []*AlertLog `json:"results,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type AlertLog struct {
	Alert     Alert `json:"alert"`
	CreatedAt int64 `json:"created_at"`

	Actions []*ActionLog `json:"actions"`
}

type ActionLog struct {
	Action Action `json:"action"`

	StartedAt int64 `json:"started_at"`
	EndedAt   int64 `json:"ended_at"`
}
