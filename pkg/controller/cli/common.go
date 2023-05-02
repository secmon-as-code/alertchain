package cli

import (
	"github.com/m-mizutani/alertchain/pkg/action"
	"github.com/m-mizutani/alertchain/pkg/chain"
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/goerr"
)

func setupActions(configs []model.ActionConfig) (actions []interfaces.Action, err error) {
	// Configure actions
	for _, cfg := range configs {
		newAction, err := action.New(cfg)
		if err != nil {
			return nil, err
		}
		actions = append(actions, newAction)
	}

	return
}

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

		client, err := policy.New(policy.WithDir(cfg.Path), policy.WithPackage(pkgName))
		if err != nil {
			return nil, goerr.Wrap(err, "creating new policy.Client")
		}

		options = append(options, c.f(client))
	}

	return options, nil
}

func buildChain(cfg model.Config, options ...chain.Option) (*chain.Chain, error) {
	actions, err := setupActions(cfg.Actions)
	if err != nil {
		return nil, err
	}

	policyOptions, err := setupPolicy(cfg.Policy)
	if err != nil {
		return nil, err
	}

	if cfg.Policy.Print {
		options = append(options, chain.WithEnablePrint())
	}

	options = append(options, chain.WithAction(actions...))
	options = append(options, policyOptions...)

	return chain.New(options...)
}
