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

func WithEnricher(enrichers ...interfaces.Enricher) Option {
	return func(c *Chain) {
		c.enrichers = append(c.enrichers, enrichers...)
	}
}

func WithPolicyAlert(policy opac.Client) Option {
	return func(c *Chain) {
		c.alertPolicy = policy
	}
}

func WithPolicyEnrich(policy opac.Client) Option {
	return func(c *Chain) {
		c.enrichPolicy = policy
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
