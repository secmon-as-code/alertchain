package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type AlertPolicyResult struct {
	Alerts []AlertMetaData `json:"alert"`
}

type ActionInitRequest struct {
	Seq     int           `json:"seq"`
	Alert   Alert         `json:"alert"`
	EnvVars types.EnvVars `json:"env" masq:"secret"`
}

type ActionInitResponse struct {
	Init []Next `json:"init"`
}

func (x *ActionInitResponse) Abort() bool {
	for _, e := range x.Init {
		if e.Abort {
			return true
		}
	}
	return false
}

func (x *ActionInitResponse) Attrs() Attributes {
	var attrs Attributes
	for _, e := range x.Init {
		attrs = append(attrs, e.Attrs...)
	}
	return attrs
}

type ActionRunRequest struct {
	Alert   Alert          `json:"alert"`
	EnvVars types.EnvVars  `json:"env" masq:"secret"`
	Seq     int            `json:"seq"`
	Called  []ActionResult `json:"called"`
}

type ActionRunResponse struct {
	Runs []Action `json:"run"`
}

type Action struct {
	ID    types.ActionID   `json:"id"`
	Name  string           `json:"name"`
	Uses  types.ActionName `json:"uses"`
	Args  ActionArgs       `json:"args"`
	Force bool             `json:"force"`
}

type ActionResult struct {
	Action
	Result any `json:"result,omitempty"`
}

type Next struct {
	Abort bool        `json:"abort"`
	Attrs []Attribute `json:"attrs"`

	// Set by runAction
	Proc Action `json:"-"`
}

type ActionExitRequest struct {
	Action ActionResult   `json:"action"`
	Called []ActionResult `json:"called"`

	Alert   Alert         `json:"alert"`
	EnvVars types.EnvVars `json:"env" masq:"secret"`
	Seq     int           `json:"seq"`
}

type ActionExitResponse struct {
	Exit []Next `json:"exit"`
}

func (x *ActionExitResponse) Abort() bool {
	for _, e := range x.Exit {
		if e.Abort {
			return true
		}
	}
	return false
}

func (x *ActionExitResponse) Attrs() Attributes {
	var attrs Attributes
	for _, e := range x.Exit {
		attrs = append(attrs, e.Attrs...)
	}
	return attrs
}
