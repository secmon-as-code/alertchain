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
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Params      []types.Parameter `json:"params"`
}

type Alert struct {
	AlertMetaData
	Schema     types.Schema `json:"schema"`
	Data       any          `json:"data"`
	CreatedAt  time.Time    `json:"created_at"`
	References []Reference  `json:"reference"`

	Raw string `json:"-"`
}

type Reference struct {
	types.Parameter
	Source types.EnricherID `json:"source"`
	Data   any              `json:"data"`
}
