package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type ScenarioLog struct {
	ID    types.ScenarioID    `json:"id"`
	Title types.ScenarioTitle `json:"title"`

	Results []*PlayLog `json:"results,omitempty"`
	Error   string     `json:"error,omitempty"`
}

type PlayLog struct {
	Alert Alert `json:"alert"`

	Actions []*ActionLog `json:"actions"`
}

type ActionLog struct {
	Seq  int      `json:"seq"`
	Init []Chore  `json:"init,omitempty"`
	Run  []Action `json:"run,omitempty"`
	Exit []Chore  `json:"exit,omitempty"`
}
