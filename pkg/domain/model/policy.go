package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type AlertPolicyQuery struct {
	Label   string `json:"label"`
	Message any    `json:"message"`
}

type AlertPolicyResult struct {
	Alerts []AlertMetaData `json:"alert"`
}

type EnrichPolicyResult struct {
	Targets []types.Parameter `json:"enrich"`
}

type ActionPolicyResult struct {
	Actions []Action `json:"action"`
}
