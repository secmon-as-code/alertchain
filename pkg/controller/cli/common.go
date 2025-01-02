package cli

import (
	"github.com/secmon-lab/alertchain/pkg/chain"
	"github.com/secmon-lab/alertchain/pkg/controller/cli/config"
	"github.com/secmon-lab/alertchain/pkg/logging"
)

func buildChain(policy *config.Policy, options ...chain.Option) (*chain.Chain, error) {
	if policy.Print() {
		logging.Default().Info("enable print mode")
		options = append(options, chain.WithEnablePrint())
	}

	alertPolicy, err := policy.Load("alert")
	if err != nil {
		return nil, err
	}
	options = append(options, chain.WithPolicyAlert(alertPolicy))

	actionPolicy, err := policy.Load("action")
	if err != nil {
		return nil, err
	}
	options = append(options, chain.WithPolicyAction(actionPolicy))

	return chain.New(options...)
}
