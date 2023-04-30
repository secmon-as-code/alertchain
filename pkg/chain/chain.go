package chain

import (
	"errors"
	"time"

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

	scenarioLogger interfaces.ScenarioLogger

	disableAction bool
	enablePrint   bool
	maxStackDepth int

	now func() time.Time
}

type Option func(c *Chain)

func New(options ...Option) (*Chain, error) {
	c := &Chain{
		actionMap:      make(map[types.ActionID]interfaces.Action),
		scenarioLogger: &dummyScenarioLogger{},
		maxStackDepth:  types.DefaultMaxStackDepth,
		now:            time.Now,
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

// HandleAlert is main function of alert chain. It receives alert data and execute actions according to the Rego policies.
func (x *Chain) HandleAlert(ctx *model.Context, schema types.Schema, data any) error {
	defer func() {
		if err := x.scenarioLogger.Flush(); err != nil {
			ctx.Logger().Error("Failed to close scenario logger", "err", err)
		}
	}()

	alerts, err := x.detectAlert(ctx, schema, data)
	if err != nil {
		return goerr.Wrap(err)
	}

	if x.actionPolicy == nil {
		return nil
	}

	for _, alert := range alerts {
		alertLogger := x.scenarioLogger.NewAlertLogger(&model.AlertLog{
			Alert:     alert,
			CreatedAt: x.now().Nanosecond(),
		})
		ctx = ctx.New(model.WithAlert(alert))

		var actions model.ActionPolicyResponse
		initOpt := []opac.QueryOption{
			opac.WithPackageSuffix(".main"),
		}
		if x.enablePrint {
			initOpt = append(initOpt, opac.WithPrintWriter(newPrintHook(ctx)))
		}

		if err := x.actionPolicy.Query(ctx, alert, &actions, initOpt...); err != nil {
			return goerr.Wrap(err, "failed to evaluate alert for action").With("alert", alert)
		}

		for _, tgt := range actions.Actions {
			if err := x.runAction(ctx, alert, tgt, alertLogger.Log); err != nil {
				return err
			}
		}
	}

	return nil
}

func (x *Chain) detectAlert(ctx *model.Context, schema types.Schema, data any) ([]model.Alert, error) {
	if x.alertPolicy == nil {
		return nil, nil
	}

	var alertResult model.AlertPolicyResult
	opt := []opac.QueryOption{
		opac.WithPackageSuffix("." + string(schema)),
	}

	if x.enablePrint {
		opt = append(opt, opac.WithPrintWriter(newPrintHook(ctx)))
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

func (x *Chain) runAction(ctx *model.Context, base model.Alert, tgt model.Action, log func(log *model.ActionLog)) error {
	if ctx.Stack() > x.maxStackDepth {
		return goerr.Wrap(types.ErrMaxStackDepth).With("stack", ctx.Stack())
	}

	startAt := x.now()

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

	result, err := action.Run(ctx, alert, tgt.Args)
	if err != nil {
		return err
	}

	// query action policy with action result
	opt := []opac.QueryOption{
		opac.WithPackageSuffix("." + string(action.ID())),
	}
	if x.enablePrint {
		opt = append(opt, opac.WithPrintWriter(newPrintHook(ctx)))
	}

	request := model.ActionPolicyRequest{
		Alert:  alert,
		Result: result,
	}
	var response model.ActionPolicyResponse
	if err := x.actionPolicy.Query(ctx, request, &response, opt...); err != nil && !errors.Is(err, opac.ErrNoEvalResult) {
		return goerr.Wrap(err, "failed to evaluate action response").With("request", request)
	}

	log(&model.ActionLog{
		Action: model.Action{
			ID:     action.ID(),
			Params: alert.Params,
			Args:   tgt.Args,
		},
		Next:      response.Actions,
		StartedAt: startAt.Nanosecond(),
		EndedAt:   x.now().Nanosecond(),
	})

	newCtx := ctx.New(model.WithStackIncrement())
	for _, newTgt := range response.Actions {
		if err := x.runAction(newCtx, alert, newTgt, log); err != nil {
			return err
		}
	}

	return nil
}
