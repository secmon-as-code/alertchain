package model

import (
	"encoding/json"

	"github.com/google/go-jsonnet"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

type Config struct {
	Policy  PolicyConfig   `json:"policy"`
	Actions []ActionConfig `json:"actions"`
}

type PolicyConfig struct {
	Path    string              `json:"path"`
	Package PolicyPackageConfig `json:"package"`
}

type PolicyPackageConfig struct {
	Alert  string `json:"alert"`
	Action string `json:"action"`
}

type ActionConfig struct {
	ID     types.ActionID     `json:"id"`
	Uses   types.ActionName   `json:"uses"`
	Config ActionConfigValues `json:"config"`
}

func (x ActionConfig) Validate() error {
	if x.ID == "" {
		return goerr.Wrap(types.ErrConfigNoActionID)
	}
	if x.Uses == "" {
		return goerr.Wrap(types.ErrConfigNoActionName).With("id", x.ID)
	}
	return nil
}

type ActionConfigValues map[string]any

type EnvVar struct {
	Key   string
	Value string
}

func ParseConfig(filename string, data string, envVars []EnvVar, cfg *Config) error {
	vm := jsonnet.MakeVM()

	for _, v := range envVars {
		vm.ExtVar(v.Key, v.Value)
	}

	raw, err := vm.EvaluateAnonymousSnippet(filename, data)
	if err != nil {
		return goerr.Wrap(err, "evaluating config jsonnet")
	}

	if err := json.Unmarshal([]byte(raw), cfg); err != nil {
		return goerr.Wrap(err, "unmarshal config by jsonnet")
	}

	return nil
}
