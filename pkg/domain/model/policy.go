package model

import (
	"errors"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

type AlertPolicyResult struct {
	Alerts []AlertMetaData `json:"alert"`
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
	ID     types.ActionID   `json:"id"`
	Name   string           `json:"name"`
	Uses   types.ActionName `json:"uses"`
	Args   ActionArgs       `json:"args"`
	Force  bool             `json:"force"`
	Abort  bool             `json:"abort"`
	Commit []Commit         `json:"commit"`
}

func (x Action) Copy() Action {
	copied := x
	return copied
}

type Commit struct {
	Attribute
	Path string `json:"path"`
}

func (x Commit) Copy() Commit {
	return Commit{
		Attribute: x.Attribute.Copy(),
		Path:      x.Path,
	}
}

func (x *Commit) ToAttr(data any) (*Attribute, error) {
	attr := x.Attribute

	if x.Path == "" {
		if attr.Value == nil {
			return nil, goerr.New("Path is empty and Value is nil", goerr.V("attr", attr))
		}
		return &attr, nil
	}

	if data == nil {
		return nil, goerr.New("Data is nil", goerr.V("commit", x))
	}

	builder := gval.Full(jsonpath.PlaceholderExtension())
	dst, err := builder.Evaluate(x.Path, data)
	if err != nil {
		if unwrapped := errors.Unwrap(err); unwrapped != nil && unwrapped.Error() == "unknown key invalid" {
			if attr.Value == nil {
				return nil, nil
			}
			return &attr, nil
		}

		return nil, goerr.Wrap(err, "failed to evaluate JSON path", goerr.V("path", x.Path), goerr.V("data", data))
	}
	attr.Value = dst

	return &attr, nil
}

type ActionResult struct {
	Action
	Result any `json:"result,omitempty"`
}
