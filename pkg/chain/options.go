package chain

import (
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/opac"
)

func WithAction(actions ...interfaces.Action) Option {
	return func(c *Chain) {
		c.actions = append(c.actions, actions...)
	}
}

func WithPolicyAlert(policy opac.Client) Option {
	return func(c *Chain) {
		c.alertPolicy = policy
	}
}

func WithPolicyEnrich(policy opac.Client) Option {
	return func(c *Chain) {
		c.inspectPolicy = policy
	}
}

func WithPolicyAction(policy opac.Client) Option {
	return func(c *Chain) {
		c.actionPolicy = policy
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
