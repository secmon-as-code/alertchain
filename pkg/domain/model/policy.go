package model

type AlertPolicyResult struct {
	Alerts []AlertMetaData `json:"alert"`
}

type ActionPolicyRequest struct {
	Alert  Alert `json:"alert"`
	Result any   `json:"result"`
}

type ActionPolicyResponse struct {
	Actions []Action `json:"action"`
}
