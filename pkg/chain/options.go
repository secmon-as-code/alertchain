package chain

import (
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
)

func WithAction(actions ...interfaces.Action) Option {
	return func(c *Chain) {
		c.actions = append(c.actions, actions...)
	}
}

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
