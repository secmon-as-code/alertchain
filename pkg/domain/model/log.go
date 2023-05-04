package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type ScenarioLog struct {
	ID    types.ScenarioID    `json:"id"`
	Title types.ScenarioTitle `json:"title"`

	AlertLog []*AlertLog `json:"alerts,omitempty"`
}

type AlertLog struct {
	Alert     Alert `json:"alert"`
	CreatedAt int   `json:"created_at"`

	Actions []*ActionLog `json:"actions"`
}

type ActionLog struct {
	Action Proc   `json:"action"`
	Next   []Proc `json:"next"`

	StartedAt int `json:"started_at"`
	EndedAt   int `json:"ended_at"`
}
