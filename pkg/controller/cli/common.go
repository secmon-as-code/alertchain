package cli

import (
	"github.com/m-mizutani/alertchain/pkg/chain"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
	"golang.org/x/exp/slog"
)

func setupPolicy(cfg model.PolicyConfig) ([]chain.Option, error) {
	if cfg.Path == "" {
		return nil, goerr.Wrap(types.ErrConfigNoPolicyPath)
	}

	configs := []struct {
		pkgName     string
		defaultName string
		f           func(*policy.Client) chain.Option
	}{
		{
			pkgName:     cfg.Package.Alert,
			defaultName: "alert",
			f:           chain.WithPolicyAlert,
		},
		{
			pkgName:     cfg.Package.Action,
			defaultName: "action",
			f:           chain.WithPolicyAction,
		},
	}

	var options []chain.Option

	for _, c := range configs {
		pkgName := c.defaultName
		if c.pkgName != "" {
			pkgName = c.pkgName
		}

		utils.Logger().Info("loading policy",
			slog.String("package", pkgName),
			slog.String("path", cfg.Path),
		)
		client, err := policy.New(policy.WithDir(cfg.Path), policy.WithPackage(pkgName))
		if err != nil {
			return nil, goerr.Wrap(err, "creating new policy.Client")
		}

		options = append(options, c.f(client))
	}

	return options, nil
}

func buildChain(cfg model.Config, options ...chain.Option) (*chain.Chain, error) {
	policyOptions, err := setupPolicy(cfg.Policy)
	if err != nil {
		return nil, err
	}

	if cfg.Policy.Print {
		utils.Logger().Info("enable print mode")
		options = append(options, chain.WithEnablePrint())
	}

	options = append(options, policyOptions...)

	return chain.New(options...)
}
