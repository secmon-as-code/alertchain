package model

import (
	"time"

	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type Parameter struct {
	Key   string
	Value string
}

type AlertMetaData struct {
	Title  string            `json:"title"`
	Source string            `json:"source"`
	Params []types.Parameter `json:"params"`
}

type Alert struct {
	AlertMetaData
	Data       any         `json:"data"`
	CreatedAt  time.Time   `json:"created_at"`
	References []Reference `json:"reference"`
}

type Reference struct {
	types.Parameter
	Source types.EnricherID `json:"source"`
	Data   any              `json:"data"`
}
