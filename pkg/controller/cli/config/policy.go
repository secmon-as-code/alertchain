package config

import (
	"context"
	"log/slog"

	"github.com/secmon-lab/alertchain/pkg/chain"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/infra/policy"
	"github.com/urfave/cli/v3"
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
			Sources:     cli.EnvVars("ALERTCHAIN_ENABLE_PRINT"),
			Value:       false,
			Destination: &x.print,
		},
		&cli.StringFlag{
			Name:        "policy-dir",
			Usage:       "directory path of policy files",
			Category:    category,
			Aliases:     []string{"d"},
			Sources:     cli.EnvVars("ALERTCHAIN_POLICY_DIR"),
			Required:    true,
			Destination: &x.path,
		},
	}
}

func (x *Policy) Load(ctx context.Context, pkgName string) (*policy.Client, error) {
	ctxutil.Logger(ctx).Info("loading policy",
		slog.String("package", pkgName),
		slog.String("path", x.path),
	)
	return policy.New(policy.WithDir(x.path), policy.WithPackage(pkgName))
}

func (x *Policy) CoreOption(ctx context.Context) ([]chain.Option, error) {
	var options []chain.Option

	if x.Print() {
		ctxutil.Logger(ctx).Info("enable print mode")
		options = append(options, chain.WithEnablePrint())
	}

	alertPolicy, err := x.Load(ctx, "alert")
	if err != nil {
		return nil, err
	}
	options = append(options, chain.WithPolicyAlert(alertPolicy))

	actionPolicy, err := x.Load(ctx, "action")
	if err != nil {
		return nil, err
	}
	options = append(options, chain.WithPolicyAction(actionPolicy))

	return options, nil
}
