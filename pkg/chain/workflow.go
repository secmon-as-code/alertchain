package chain

import (
	"github.com/m-mizutani/alertchain/pkg/chain/core"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
	"golang.org/x/exp/slog"
)

type workflow struct {
	core    *core.Core
	alert   model.Alert
	options []policy.QueryOption
}

func newWorkflow(c *core.Core, alert model.Alert) (*workflow, error) {

	hdlr := &workflow{
		core:  c,
		alert: alert,
	}

	return hdlr, nil
}

func (x *workflow) Run(ctx *model.Context) error {
	copied := x.alert.Copy()
	logger := x.core.ScenarioLogger().NewAlertLogger(&copied)
	defer x.core.ScenarioLogger().Flush()

	envVars := x.core.Env()

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

	ctx = ctx.New(model.WithAlert(x.alert))
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

		if len(p.init) > 0 || len(p.run) > 0 || len(p.exit) > 0 {
			actionLogger := logger.NewActionLogger()
			actionLogger.LogInit(p.init)
			actionLogger.LogRun(p.run)
			actionLogger.LogExit(p.exit)
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
				ttl := types.DefaultAttributeTTL
				if x.alert.Attrs[i].TTL > 0 {
					ttl = x.alert.Attrs[i].TTL
				}
				x.alert.Attrs[i].ExpiresAt = x.core.Now().UTC().Unix() + ttl
				global = append(global, x.alert.Attrs[i])
			}
		}

		if err := x.core.DBClient().PutAttrs(ctx, x.alert.Namespace, global); err != nil {
			return goerr.Wrap(err, "failed to put global attrs")
		}

		ctx.Logger().Info("saved global attributes", slog.Any("attrs", global))
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
	init []model.Chore
	run  []model.Action
	exit []model.Chore

	history   *actionHistory
	finalized model.Attributes
}

func (x *proc) aborted() bool {
	for _, i := range x.init {
		if i.Abort {
			return true
		}
	}

	for _, e := range x.exit {
		if e.Abort {
			return true
		}
	}
	return false
}

func (x *proc) evaluate(ctx *model.Context) error {
	// Evaluate `init` rules
	initReq := model.ActionInitRequest{
		Seq:     x.seq,
		Alert:   x.alert,
		EnvVars: x.envVars,
	}
	var initResp model.ActionInitResponse
	if err := x.core.QueryActionPolicy(ctx, initReq, &initResp); err != nil {
		return err
	}

	x.init = initResp.Init
	x.alert.Attrs = append(x.alert.Attrs, initResp.Attrs()...).Tidy()
	x.finalized = x.alert.Attrs[:]

	if initResp.Abort() {
		return nil
	}

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

	for _, p := range runResp.Runs {
		if p.ID == "" {
			p.ID = types.NewActionID()
		} else if x.history.alreadyCalled(p.ID) {
			continue
		}

		x.run = append(x.run, p)
		result, err := x.executeAction(ctx, p, x.alert)
		if err != nil {
			return err
		}

		actionResult := model.ActionResult{
			Action: p,
			Result: result,
		}
		x.history.add(actionResult)

		exitReq := model.ActionExitRequest{
			Seq:     x.seq,
			Alert:   x.alert,
			Action:  actionResult,
			EnvVars: x.envVars,
			Called:  x.history.called,
		}
		var exitResp model.ActionExitResponse
		if err := x.core.QueryActionPolicy(ctx, exitReq, &exitResp); err != nil {
			return err
		}

		x.exit = append(x.exit, exitResp.Exit...)
		x.finalized = append(x.finalized, exitResp.Attrs()...).Tidy()
	}

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
			return nil, err
		}
		result = resp
	}

	return result, nil
}
