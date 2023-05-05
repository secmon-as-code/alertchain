package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type AlertPolicyResult struct {
	Alerts []AlertMetaData `json:"alert"`
}

type ActionRunRequest struct {
	Called []Process `json:"called"`

	Alert   Alert         `json:"alert"`
	EnvVars types.EnvVars `json:"env"`
	Seq     int           `json:"seq"`
}

type ActionRunResponse struct {
	Runs []Process `json:"run"`
}

type Process struct {
	ID      types.ProcessID  `json:"id"`
	Uses    types.ActionName `json:"uses"`
	Args    ActionArgs       `json:"args"`
	Secrets ActionSecrets    `json:"secrets"`

	Result any `json:"result"`
}

type Exit struct {
	Abort  bool        `json:"abort"`
	Params []Parameter `json:"params"`

	// Set by runAction
	Proc Process `json:"-"`
}

type ActionExitRequest struct {
	Proc   Process   `json:"proc"`
	Called []Process `json:"called"`

	Alert   Alert         `json:"alert"`
	EnvVars types.EnvVars `json:"env"`
	Seq     int           `json:"seq"`
}

type ActionExitResponse struct {
	Exit []Exit `json:"exit"`
}
