package model

import (
	"encoding/json"

	"github.com/google/go-jsonnet"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
)

type Config struct {
	Policy PolicyConfig `json:"policy"`
}

func (x *Config) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "enable-print",
			Category:    "logging",
			Aliases:     []string{"p"},
			EnvVars:     []string{"ALERTCHAIN_ENABLE_PRINT"},
			Usage:       "Enable print feature in Rego. The cli option is priority than config file.",
			Value:       false,
			Destination: &x.Policy.Print,
		},
		&cli.StringFlag{
			Name:        "policy-dir",
			Aliases:     []string{"d"},
			Usage:       "directory path of policy files",
			EnvVars:     []string{"ALERTCHAIN_POLICY_DIR"},
			Required:    true,
			Destination: &x.Policy.Path,
		},
	}
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
