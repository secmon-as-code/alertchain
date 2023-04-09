package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type ActionParams map[string]any

type Action struct {
	ID     types.ActionID `json:"id"`
	Params ActionParams   `json:"params"`
}
