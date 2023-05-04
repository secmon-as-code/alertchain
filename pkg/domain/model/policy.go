package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type AlertPolicyResult struct {
	Alerts []AlertMetaData `json:"alert"`
}

type ActionRunRequest struct {
	Called []Proc `json:"called"`

	Alert   Alert         `json:"alert"`
	EnvVars types.EnvVars `json:"env"`
}

type ActionRunResponse struct {
	Runs []Proc `json:"run"`
}

type Proc struct {
	ID      types.ActionID   `json:"id"`
	Uses    types.ActionName `json:"uses"`
	Args    ActionArgs       `json:"args"`
	Secrets ActionSecrets    `json:"secrets"`

	Result any `json:"result"`
}

type Exit struct {
	Abort  bool        `json:"abort"`
	Params []Parameter `json:"params"`

	// Set by runAction
	Proc Proc `json:"-"`
}

type ActionExitRequest struct {
	Proc   Proc   `json:"proc"`
	Called []Proc `json:"called"`

	Alert   Alert         `json:"alert"`
	EnvVars types.EnvVars `json:"env"`
}

type ActionExitResponse struct {
	Exit []Exit `json:"exit"`
}
