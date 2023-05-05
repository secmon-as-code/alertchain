package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type AlertPolicyResult struct {
	Alerts []AlertMetaData `json:"alert"`
}

type ActionRunRequest struct {
	Alert   Alert         `json:"alert"`
	EnvVars types.EnvVars `json:"env"`
	Seq     int           `json:"seq"`
	Called  []Action      `json:"called"`
}

type ActionRunResponse struct {
	Runs []Action `json:"run"`
}

type Action struct {
	ID   types.ActionID   `json:"id"`
	Uses types.ActionName `json:"uses"`
	Args ActionArgs       `json:"args"`

	Result any `json:"result"`
}

type Exit struct {
	Abort  bool        `json:"abort"`
	Params []Parameter `json:"params"`

	// Set by runAction
	Proc Action `json:"-"`
}

type ActionExitRequest struct {
	Action Action   `json:"action"`
	Called []Action `json:"called"`

	Alert   Alert         `json:"alert"`
	EnvVars types.EnvVars `json:"env"`
	Seq     int           `json:"seq"`
}

type ActionExitResponse struct {
	Exit []Exit `json:"exit"`
}
