package model

import "github.com/secmon-lab/alertchain/pkg/domain/types"

type ScenarioLog struct {
	ID    types.ScenarioID    `json:"id"`
	Title types.ScenarioTitle `json:"title"`

	Results []*PlayLog `json:"results,omitempty"`
	Error   any        `json:"error,omitempty"`
}

type PlayLog struct {
	Alert Alert `json:"alert"`

	Actions []*ActionLog `json:"actions"`
}

type ActionLog struct {
	Seq int `json:"seq"`
	Action
}
