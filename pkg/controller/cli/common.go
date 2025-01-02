package cli

import (
	"context"

	"github.com/secmon-lab/alertchain/pkg/chain"
	"github.com/secmon-lab/alertchain/pkg/controller/cli/config"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
)

func buildChain(ctx context.Context, policy *config.Policy, options ...chain.Option) (*chain.Chain, error) {
	if policy.Print() {
		ctxutil.Logger(ctx).Info("enable print mode")
		options = append(options, chain.WithEnablePrint())
	}

	alertPolicy, err := policy.Load(ctx, "alert")
	if err != nil {
		return nil, err
	}
	options = append(options, chain.WithPolicyAlert(alertPolicy))

	actionPolicy, err := policy.Load(ctx, "action")
	if err != nil {
		return nil, err
	}
	options = append(options, chain.WithPolicyAction(actionPolicy))

	return chain.New(options...)
}
