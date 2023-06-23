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
	calledProc []model.Action
}

func (x *Chain) newWorkflow(alert model.Alert, options []policy.QueryOption) (*workflow, error) {
	copied, err := alert.Copy()
	if err != nil {
		return nil, err
	}

	logger := x.scenarioLogger.NewAlertLogger(&model.AlertLog{
		Alert:     copied,
		CreatedAt: x.now().UnixNano(),
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

	envVars := x.chain.Env()

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

		attrAll := x.alert.Attrs[:]
		for _, e := range resp.Exits {
			attrAll = append(attrAll, e.Attrs...)
		}
		x.alert.Attrs = model.TidyAttributes(attrAll)

		if len(resp.Called) == 0 {
			break
		}
	}

	return nil
}

func (x *workflow) alreadyCalled(id types.ActionID) bool {
	for _, p := range x.calledProc {
		if p.ID == id {
			return true
		}
	}
	return false
}

type runActionResponse struct {
	Exits  []model.Exit
	Called []model.Action
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
			p.ID = types.NewActionID()
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
			Action: p,
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

func (x *workflow) runProc(ctx *model.Context, p model.Action, alert model.Alert) (any, error) {
	run, ok := x.chain.actionMap[p.Uses]
	if !ok {
		return nil, goerr.Wrap(types.ErrActionNotFound).With("uses", p.Uses)
	}
	log := &model.ActionLog{
		Action:    p,
		StartedAt: x.chain.now().UnixNano(),
	}
	defer x.logger.Log(log)

	utils.Logger().Info("run action", slog.Any("proc", p))

	// Run action. If actionMock is set, use it instead of action.Run()
	var result any
	if x.chain.actionMock != nil {
		result = x.chain.actionMock.GetResult(p.Uses)
	} else if !x.chain.disableAction {
		resp, err := run(ctx, alert, p.Args)
		if err != nil {
			return nil, err
		}
		result = resp
	}
	log.EndedAt = x.chain.now().UnixNano()

	return result, nil
}

func Env() types.EnvVars {
	vars := types.EnvVars{}
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		vars[types.EnvVarName(pair[0])] = types.EnvVarValue(pair[1])
	}
	return vars
}
