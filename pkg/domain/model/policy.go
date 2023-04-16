package model

type AlertPolicyResult struct {
	Alerts []AlertMetaData `json:"alert"`
}

type ActionPolicyResult struct {
	Actions []Action `json:"action"`
}
