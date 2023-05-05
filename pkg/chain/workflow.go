package chain

import (
	"errors"
	"os"
	"strings"

	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
	"golang.org/x/exp/slog"
)

type workflow struct {
	chain   *Chain
	logger  interfaces.AlertLogger
	alert   model.Alert
	options []policy.QueryOption

	// mutable variables
	calledProc []model.Process
}

func (x *Chain) newWorkflow(alert model.Alert, options []policy.QueryOption) (*workflow, error) {
	copied, err := alert.Copy()
	if err != nil {
		return nil, err
	}

	logger := x.scenarioLogger.NewAlertLogger(&model.AlertLog{
		Alert:     copied,
		CreatedAt: x.now().Nanosecond(),
	})

	hdlr := &workflow{
		chain:   x,
		alert:   alert,
		logger:  logger,
		options: options,
	}

	return hdlr, nil
}

func (x *workflow) run(ctx *model.Context) error {
	ctx = ctx.New(model.WithAlert(x.alert))

	envVars := buildEnvVars()

	for i := 0; i < x.chain.maxStackDepth; i++ {
		runReq := &model.ActionRunRequest{
			Alert:   x.alert,
			EnvVars: envVars,
			Called:  x.calledProc,
		}

		resp, err := x.runAction(ctx, runReq)
		if err != nil {
			return err
		}
		x.calledProc = append(x.calledProc, resp.Called...)

		paramAll := x.alert.Params[:]
		for _, e := range resp.Exits {
			paramAll = append(paramAll, e.Params...)
		}
		x.alert.Params = model.TidyParameters(paramAll)

		if len(resp.Called) == 0 {
			break
		}
	}

	return nil
}

func (x *workflow) alreadyCalled(id types.ProcessID) bool {
	for _, p := range x.calledProc {
		if p.ID == id {
			return true
		}
	}
	return false
}

type runActionResponse struct {
	Exits  []model.Exit
	Called []model.Process
}

func (x *workflow) runAction(ctx *model.Context, runReq *model.ActionRunRequest) (*runActionResponse, error) {
	var runResp model.ActionRunResponse
	var resp runActionResponse

	ctx.Logger().Debug("request action.run policy", slog.Any("request", runReq))
	if err := x.chain.actionPolicy.Query(ctx, runReq, &runResp, x.options...); err != nil && !errors.Is(err, types.ErrNoPolicyResult) {
		return nil, goerr.Wrap(err, "failed to evaluate action.run").With("request", runReq)
	}
	ctx.Logger().Debug("response action.run policy", slog.Any("response", runResp))

	for _, p := range runResp.Runs {
		if p.ID == "" {
			p.ID = types.NewProcessID()
		} else if x.alreadyCalled(p.ID) {
			continue
		}

		result, err := x.runProc(ctx, p, runReq.Alert)
		if err != nil {
			return nil, err
		}

		p.Result = result
		resp.Called = append(resp.Called, p)

		exitReq := model.ActionExitRequest{
			Alert:  runReq.Alert,
			Proc:   p,
			Called: x.calledProc,
		}
		var exitResp model.ActionExitResponse

		ctx.Logger().Debug("request action.exit policy", slog.Any("request", exitReq))
		if err := x.chain.actionPolicy.Query(ctx, exitReq, &exitResp, x.options...); err != nil && !errors.Is(err, types.ErrNoPolicyResult) {
			return nil, goerr.Wrap(err, "failed to evaluate action.exit").With("request", exitReq)
		}
		ctx.Logger().Debug("response action.exit policy", slog.Any("response", exitResp))

		resp.Exits = append(resp.Exits, exitResp.Exit...)
	}

	return &resp, nil
}

func (x *workflow) runProc(ctx *model.Context, p model.Process, alert model.Alert) (any, error) {
	run, ok := x.chain.actionMap[p.Uses]
	if !ok {
		return nil, goerr.Wrap(types.ErrActionNotFound).With("uses", p.Uses)
	}
	utils.Logger().Info("run action", slog.Any("proc", p))

	// Run action. If actionMock is set, use it instead of action.Run()
	if x.chain.actionMock != nil {
		return x.chain.actionMock.GetResult(p.Uses), nil
	} else if !x.chain.disableAction {
		return run(ctx, alert, p.Args)
	}
	return nil, nil
}

func buildEnvVars() types.EnvVars {
	vars := types.EnvVars{}
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		vars[types.EnvVarName(pair[0])] = types.EnvVarValue(pair[1])
	}
	return vars
}
