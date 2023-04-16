package chain

import (
	"errors"

	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/opac"
	"golang.org/x/exp/slog"
)

type Chain struct {
	actions   []interfaces.Action
	actionMap map[types.ActionID]interfaces.Action

	alertPolicy   opac.Client
	inspectPolicy opac.Client
	actionPolicy  opac.Client

	disableAction bool
}

type Option func(c *Chain)

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
	alerts, err := x.detectAlert(ctx, schema, data)
	if err != nil {
		return goerr.Wrap(err)
	}

	if x.actionPolicy == nil {
		return nil
	}

	for _, alert := range alerts {
		var actions model.ActionPolicyResult
		initOpt := []opac.QueryOption{
			opac.WithPackageSuffix(".main"),
		}
		if err := x.actionPolicy.Query(ctx, alert, &actions, initOpt...); err != nil {
			return goerr.Wrap(err, "failed to evaluate alert for action").With("alert", alert)
		}

		for _, tgt := range actions.Actions {
			if err := x.runAction(ctx, alert, tgt); err != nil {
				return err
			}
		}
	}

	return nil
}

func (x *Chain) detectAlert(ctx *types.Context, schema types.Schema, data any) ([]model.Alert, error) {
	if x.alertPolicy == nil {
		return nil, nil
	}

	var alertResult model.AlertPolicyResult
	opt := []opac.QueryOption{
		opac.WithPackageSuffix("." + string(schema)),
	}
	if err := x.alertPolicy.Query(ctx, data, &alertResult, opt...); err != nil {
		return nil, goerr.Wrap(err)
	}

	if len(alertResult.Alerts) == 0 {
		return nil, nil
	}

	alerts := make([]model.Alert, len(alertResult.Alerts))
	for i, meta := range alertResult.Alerts {
		alerts[i] = model.NewAlert(meta, schema, data)
	}
	return alerts, nil
}

func (x *Chain) runAction(ctx *types.Context, base model.Alert, tgt model.Action) error {
	alert := base.Clone(tgt.Params...)

	action, ok := x.actionMap[tgt.ID]
	if !ok {
		return goerr.Wrap(types.ErrNoSuchActionID).With("ID", tgt.ID)
	}

	utils.Logger().Info("action triggered", slog.Any("id", action.ID()))
	if x.disableAction {
		utils.Logger().Info("disable-action option is true, skip action")
		return nil
	}

	resp, err := action.Run(ctx, alert, tgt.Args)
	if err != nil {
		return err
	}

	var result model.ActionPolicyResult
	opt := []opac.QueryOption{
		opac.WithPackageSuffix("." + string(action.ID())),
	}
	if err := x.actionPolicy.Query(ctx, resp, &result, opt...); err != nil && !errors.Is(err, opac.ErrNoEvalResult) {
		return goerr.Wrap(err, "failed to evaluate action response").With("response", resp)
	}

	for _, newTgt := range result.Actions {
		if err := x.runAction(ctx, alert, newTgt); err != nil {
			return err
		}
	}

	return nil
}
