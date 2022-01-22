package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type Action interface {
	Run(ctx *types.Context, alert *Alert, args ...*Attribute) (*ChangeRequest, error)
}

type ActionDefinition struct {
	ID     string                 `json:"id"`
	Use    string                 `json:"use"`
	Config map[string]interface{} `json:"config"`
}

type ActionConfig map[string]interface{}

type ActionFactory func(config ActionConfig) (Action, error)

type ActionRepository struct {
	Factories map[string]ActionFactory
}
