package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type ActionArgs map[string]any

type Action struct {
	ID   types.ActionID `json:"id"`
	Args ActionArgs     `json:"args"`

	Params []Parameter `json:"params"`
}
