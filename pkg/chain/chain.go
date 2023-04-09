package chain

import (
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/opac"
	"golang.org/x/exp/slog"
)

type Chain struct {
	enrichers []interfaces.Enricher
	actions   []interfaces.Action
	actionMap map[types.ActionID]interfaces.Action

	alertPolicy  opac.Client
	enrichPolicy opac.Client
	actionPolicy opac.Client
}

type Option func(c *Chain)

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

func New(options ...Option) (*Chain, error) {
	c := &Chain{
		actionMap: make(map[types.ActionID]interfaces.Action),
	}

	for _, opt := range options {
		opt(c)
	}

	for _, action := range c.actions {
		if _, exists := c.actionMap[action.ID()]; exists {
			return nil, goerr.Wrap(types.ErrConfigConflictActionID).With("id", action.ID())
		}

		c.actionMap[action.ID()] = action
	}

	return c, nil
}

func (x *Chain) HandleAlert(ctx *types.Context, schema types.Schema, data any) error {
	// Step 1: Detect alert
	queryOpt := []opac.QueryOption{
		opac.WithPackageSuffix("." + string(schema)),
	}

	var alertResult model.AlertPolicyResult
	if x.alertPolicy != nil {
		if err := x.alertPolicy.Query(ctx, data, &alertResult, queryOpt...); err != nil {
			return goerr.Wrap(err)
		}
	} else {
		alertResult.Alerts = []model.AlertMetaData{
			{
				Title:  "N/A",
				Params: []types.Parameter{},
			},
		}
	}

	for _, meta := range alertResult.Alerts {
		// Step 2: Enrich indicators in the alert
		alert := model.NewAlert(meta, schema, data)

		utils.Logger().Debug("alert detected", slog.Any("alert", alert))

		if x.enrichPolicy != nil {
			var enrich model.EnrichPolicyResult
			if err := x.enrichPolicy.Query(ctx, alert, &enrich); err != nil {
				return goerr.Wrap(err, "failed to evaluate alert for enrich").With("alert", alert)
			}

			for _, tgt := range enrich.Targets {
				for _, enricher := range x.enrichers {
					ref, err := enricher.Enrich(ctx, tgt)
					if err != nil {
						return err
					}
					alert.References = append(alert.References, ref...)
				}
			}
		}

		// Step 3: Do action(s) if required
		var actions model.ActionPolicyResult
		if x.actionPolicy != nil {
			if err := x.actionPolicy.Query(ctx, alert, &actions); err != nil {
				return goerr.Wrap(err, "failed to evaluate alert for action").With("alert", alert)
			}

			for _, tgt := range actions.Actions {
				action, ok := x.actionMap[tgt.ID]
				if !ok {
					return goerr.Wrap(types.ErrNoSuchActionID).With("ID", tgt.ID)
				}

				if err := action.Run(ctx, alert, tgt.Params); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
