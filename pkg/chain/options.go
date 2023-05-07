package chain

import (
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
)

func WithPolicyAlert(p *policy.Client) Option {
	return func(c *Chain) {
		c.alertPolicy = p
	}
}

func WithPolicyAction(p *policy.Client) Option {
	return func(c *Chain) {
		c.actionPolicy = p
	}
}

func WithDisableAction() Option {
	return func(c *Chain) {
		c.disableAction = true
	}
}

func WithEnablePrint() Option {
	return func(c *Chain) {
		c.enablePrint = true
	}
}

func WithExtraAction(name types.ActionName, action interfaces.RunAction) Option {
	return func(c *Chain) {
		if _, ok := c.actionMap[name]; ok {
			panic("action name is already registered: " + name)
		}
		c.actionMap[name] = action
	}
}

func WithActionMock(mock interfaces.ActionMock) Option {
	return func(c *Chain) {
		c.actionMock = mock
	}
}

func WithScenarioLogger(logger interfaces.ScenarioLogger) Option {
	return func(c *Chain) {
		c.scenarioLogger = logger
	}
}

func WithEnv(f interfaces.Env) Option {
	return func(c *Chain) {
		c.Env = f
	}
}
