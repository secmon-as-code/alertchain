package config

import (
	"github.com/urfave/cli/v2"
)

type Policy struct {
	path  string
	print bool
}

func (x *Policy) Path() string { return x.path }
func (x *Policy) Print() bool  { return x.print }

func (x *Policy) Flags() []cli.Flag {
	category := "Policy"

	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "enable-print",
			Usage:       "Enable print feature in Rego. The cli option is priority than config file.",
			Category:    category,
			Aliases:     []string{"p"},
			EnvVars:     []string{"ALERTCHAIN_ENABLE_PRINT"},
			Value:       false,
			Destination: &x.print,
		},
		&cli.StringFlag{
			Name:        "policy-dir",
			Usage:       "directory path of policy files",
			Category:    category,
			Aliases:     []string{"d"},
			EnvVars:     []string{"ALERTCHAIN_POLICY_DIR"},
			Required:    true,
			Destination: &x.path,
		},
	}
}
