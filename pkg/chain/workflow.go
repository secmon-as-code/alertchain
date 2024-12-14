package chain

import (
	"errors"
	"log/slog"

	"github.com/m-mizutani/alertchain/pkg/chain/core"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/alertchain/pkg/service"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
)

type workflow struct {
	core    *core.Core
	alert   model.Alert
	options []policy.QueryOption
	service *service.Workflow
}

func newWorkflow(c *core.Core, alert model.Alert, svc *service.Workflow) *workflow {
	hdlr := &workflow{
		core:    c,
		alert:   alert,
		service: svc,
	}

	return hdlr
}

func (x *workflow) Run(ctx *model.Context) error {
	copied := x.alert.Copy()
	logger := x.core.ScenarioLogger().NewAlertLogger(&copied)
	envVars := x.core.Env()

	ctx = ctx.New(model.WithAlert(x.alert))

	if x.alert.Namespace != "" {
		timeoutAt := x.core.Now().Add(x.core.Timeout())
		if err := x.core.DBClient().Lock(ctx, x.alert.Namespace, timeoutAt); err != nil {
			return goerr.Wrap(err, "failed to lock namespace")
		}
		defer func() {
			if err := x.core.DBClient().Unlock(ctx, x.alert.Namespace); err != nil {
				ctx.Logger().Error("failed to unlock", slog.Any("alert", x.alert))
			}
		}()

		global, err := x.core.DBClient().GetAttrs(ctx, x.alert.Namespace)
		if err != nil {
			return goerr.Wrap(err, "failed to get global attrs")
		}

		ctx.Logger().Info("loaded global attributes", slog.Any("attrs", global))

		x.alert.Attrs = append(x.alert.Attrs, global...).Tidy()
	}

	var history actionHistory

	for i := 0; i < x.core.MaxSequences(); i++ {
		p := &proc{
			seq:     i,
			alert:   x.alert,
			core:    x.core,
			options: x.options,
			envVars: envVars,
			history: &history,
		}

		if err := p.evaluate(ctx); err != nil {
			return err
		}

		if len(p.run) > 0 {
			actionLogger := logger.NewActionLogger()
			actionLogger.LogRun(p.run)
		}

		x.alert.Attrs = p.finalized

		if len(p.run) == 0 || p.aborted() {
			break
		}
	}

	if x.alert.Namespace != "" {
		var global model.Attributes
		for i := range x.alert.Attrs {
			if x.alert.Attrs[i].Global {
				global = append(global, x.alert.Attrs[i])
			}
		}

		if err := x.core.DBClient().PutAttrs(ctx, x.alert.Namespace, global); err != nil {
			return goerr.Wrap(err, "failed to put global attrs")
		}

		ctx.Logger().Info("saved global attributes", slog.Any("attrs", global))
	}

	if err := x.service.UpdateLastAttrs(ctx, x.alert.Attrs); err != nil {
		return err
	}

	return nil
}

type actionHistory struct {
	called []model.ActionResult
}

func (x *actionHistory) add(result model.ActionResult) {
	x.called = append(x.called, result)
}

func (x *actionHistory) alreadyCalled(id types.ActionID) bool {
	for _, p := range x.called {
		if p.ID == id {
			return true
		}
	}
	return false
}

type proc struct {
	seq int

	alert   model.Alert
	core    *core.Core
	options []policy.QueryOption
	envVars types.EnvVars

	// logs
	run []model.Action

	history   *actionHistory
	finalized model.Attributes
}

func (x *proc) aborted() bool {
	for _, r := range x.run {
		if r.Abort {
			return true
		}
	}

	return false
}

func (x *proc) evaluate(ctx *model.Context) error {
	// Evaluate `run` rules
	runReq := &model.ActionRunRequest{
		Alert:   x.alert,
		EnvVars: x.envVars,
		Called:  x.history.called,
		Seq:     x.seq,
	}

	var runResp model.ActionRunResponse
	if err := x.core.QueryActionPolicy(ctx, runReq, &runResp); err != nil {
		return err
	}

	x.finalized = x.alert.Attrs.Copy()

	for _, p := range runResp.Runs {
		if err := x.runAction(ctx, p); err != nil {
			if errors.Is(err, errActionAbort) {
				break
			}
			return err
		}
	}

	return nil
}

var errActionAbort = goerr.New("action aborted")

func (x *proc) runAction(ctx *model.Context, p model.Action) error {
	if p.ID == "" {
		p.ID = types.NewActionID()
	} else if x.history.alreadyCalled(p.ID) {
		return nil
	}

	var newAttrs model.Attributes

	defer func() {
		copied := p.Copy()
		copied.Commit = make([]model.Commit, len(newAttrs))
		for i := range newAttrs {
			copied.Commit[i].Attribute = newAttrs[i]
		}
		x.run = append(x.run, copied)
	}()

	if p.Abort {
		utils.Logger().Info("abort action", slog.Any("action", p))
		return errActionAbort
	}

	var result any
	if p.Uses != "" {
		r, err := x.executeAction(ctx, p, x.alert)
		if err != nil {
			if !p.Force {
				return err
			}

			utils.Logger().Warn("got error, but force run action",
				slog.Any("action", p),
				utils.ErrLog(err),
			)
		}
		result = r
	}

	for _, c := range p.Commit {
		attr, err := c.ToAttr(result)
		if err != nil {
			return err
		}
		if attr == nil {
			continue
		}
		newAttrs = append(newAttrs, *attr)
	}

	actionResult := model.ActionResult{
		Action: p,
		Result: result,
	}
	x.history.add(actionResult)

	x.finalized = append(x.finalized, newAttrs...).Tidy()

	return nil
}

func (x *proc) executeAction(ctx *model.Context, p model.Action, alert model.Alert) (any, error) {
	run, ok := x.core.GetAction(p.Uses)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionNotFound).With("uses", p.Uses)
	}

	utils.Logger().Info("run action", slog.Any("proc", p))

	// Run action. If actionMock is set, use it instead of action.Run()
	var result any
	if x.core.ActionMock() != nil {
		result = x.core.ActionMock().GetResult(p.Uses)
	} else if !x.core.DisableAction() {
		resp, err := run(ctx, alert, p.Args)
		if err != nil {
			return nil, types.AsActionErr(goerr.Wrap(err))
		}
		result = resp
	}

	return result, nil
}
