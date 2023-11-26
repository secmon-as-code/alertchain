package cli

import (
	"log/slog"

	"github.com/m-mizutani/alertchain/pkg/chain"
	"github.com/m-mizutani/alertchain/pkg/chain/core"
	"github.com/m-mizutani/alertchain/pkg/controller/cli/config"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
)

func setupPolicy(cfg *config.Policy) ([]core.Option, error) {
	if cfg.Path() == "" {
		return nil, goerr.Wrap(types.ErrConfigNoPolicyPath)
	}

	configs := []struct {
		pkgName string
		f       func(*policy.Client) core.Option
	}{
		{
			pkgName: "alert",
			f:       core.WithPolicyAlert,
		},
		{
			pkgName: "action",
			f:       core.WithPolicyAction,
		},
	}

	var options []core.Option

	for _, c := range configs {
		utils.Logger().Info("loading policy",
			slog.String("package", c.pkgName),
			slog.String("path", cfg.Path()),
		)
		client, err := policy.New(policy.WithDir(cfg.Path()), policy.WithPackage(c.pkgName))
		if err != nil {
			return nil, goerr.Wrap(err, "creating new policy.Client")
		}

		options = append(options, c.f(client))
	}

	return options, nil
}

func buildChain(policy *config.Policy, options ...core.Option) (*chain.Chain, error) {
	policyOptions, err := setupPolicy(policy)
	if err != nil {
		return nil, err
	}

	if policy.Print() {
		utils.Logger().Info("enable print mode")
		options = append(options, core.WithEnablePrint())
	}

	options = append(options, policyOptions...)

	return chain.New(options...)
}
