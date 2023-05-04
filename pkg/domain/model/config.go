package model

import (
	"encoding/json"

	"github.com/google/go-jsonnet"
	"github.com/m-mizutani/goerr"
)

type Config struct {
	Policy PolicyConfig `json:"policy"`
}

type PolicyConfig struct {
	Path    string              `json:"path"`
	Package PolicyPackageConfig `json:"package"`
	Print   bool                `json:"print"`
}

type PolicyPackageConfig struct {
	Alert  string `json:"alert"`
	Action string `json:"action"`
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
