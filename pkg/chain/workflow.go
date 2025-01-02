package chain

import (
	"context"
	"errors"
	"log/slog"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/interfaces"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/service"
)

func (x *Chain) runWorkflow(ctx context.Context, alert model.Alert, svc *service.Services) error {
	wfSvc, err := svc.Workflow.Create(ctx, alert)
	if err != nil {
		return err
	}

	copied := alert.Copy()
	AlertRecorder := x.recorder.NewAlertRecorder(&copied)
	logger := ctxutil.Logger(ctx)

	ctx = ctxutil.InjectAlert(ctx, &alert)

	if alert.Namespace != "" {
		timeoutAt := x.now().Add(x.timeout)
		if err := x.dbClient.Lock(ctx, alert.Namespace, timeoutAt); err != nil {
			return goerr.Wrap(err, "failed to lock namespace")
		}
		defer func() {
			if err := x.dbClient.Unlock(ctx, alert.Namespace); err != nil {
				logger.Error("failed to unlock", slog.Any("alert", alert))
			}
		}()

		persistent, err := x.dbClient.GetAttrs(ctx, alert.Namespace)
		if err != nil {
			return goerr.Wrap(err, "failed to get persistent attrs")
		}

		logger.Info("loaded persistent attributes", slog.Any("attrs", persistent))

		alert.Attrs = append(alert.Attrs, persistent...).Tidy()
	}

	var history actionHistory

	for i := 0; i < x.maxSequences; i++ {
		seq := &sequence{
			idx:               i,
			alert:             alert,
			envVars:           x.env(),
			history:           history,
			queryActionPolicy: x.queryActionPolicy,
			actionMap:         x.actionMap,
			actionMock:        x.actionMock,
		}

		results, err := seq.evaluateAndRunActions(ctx)
		if err != nil {
			return err
		}

		if len(results) > 0 {
			ActionRecorder := AlertRecorder.NewActionRecorder()

			for _, r := range results {
				ActionRecorder.Add(r.Action)
				history.add(*r)
			}
		}

		finalized := alert.Attrs.Copy()
		for _, r := range results {
			for _, c := range r.Commit {
				finalized = append(finalized, c.Attribute)
			}
		}
		alert.Attrs = finalized.Tidy()

		if len(results) == 0 || isAborted(results) {
			break
		}

	}

	if alert.Namespace != "" {
		var persistent model.Attributes
		for i := range alert.Attrs {
			if alert.Attrs[i].Persist {
				persistent = append(persistent, alert.Attrs[i])
			}
		}

		if err := x.dbClient.PutAttrs(ctx, alert.Namespace, persistent); err != nil {
			return goerr.Wrap(err, "failed to put persistent attrs")
		}

		logger.Info("saved persistent attributes", slog.Any("attrs", persistent))
	}

	if err := wfSvc.UpdateLastAttrs(ctx, alert.Attrs); err != nil {
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

func isAborted(actions []*model.ActionResult) bool {
	for _, r := range actions {
		if r.Action.Abort {
			return true
		}
	}
	return false
}

type sequence struct {
	idx               int
	alert             model.Alert
	history           actionHistory
	envVars           types.EnvVars
	queryActionPolicy func(ctx context.Context, in, out any) error
	actionMock        interfaces.ActionMock
	actionMap         map[types.ActionName]model.RunAction
}

func (x *sequence) evaluateAndRunActions(ctx context.Context) ([]*model.ActionResult, error) {
	// Evaluate `run` rules
	runReq := &model.ActionRunRequest{
		Alert:   x.alert,
		EnvVars: x.envVars,
		Called:  x.history.called,
		Seq:     x.idx,
	}

	var runResp model.ActionRunResponse
	if err := x.queryActionPolicy(ctx, runReq, &runResp); err != nil {
		return nil, err
	}

	var runActions []*model.ActionResult
	for _, p := range runResp.Runs {
		result, err := x.runAction(ctx, p)
		if result != nil {
			runActions = append(runActions, result)
		}
		if err != nil {
			// Even if action is aborted, continue to next action. The workflow will be stopped before the next iteration.
			if errors.Is(err, errActionAbort) {
				continue
			}
			return nil, err
		}
	}

	return runActions, nil
}

var errActionAbort = goerr.New("action aborted")

// runAction runs an action and returns the result. If the action is already called, it returns nil.
func (x *sequence) runAction(ctx context.Context, baseAction model.Action) (*model.ActionResult, error) {
	copied := baseAction.Copy()

	if copied.ID == "" {
		copied.ID = types.NewActionID()
	} else if x.history.alreadyCalled(copied.ID) {
		return nil, nil
	}

	logger := ctxutil.Logger(ctx)
	if copied.Abort {
		logger.Info("abort action", slog.Any("action", copied))
		return nil, errActionAbort
	}

	var result any
	if copied.Uses != "" {
		run, ok := x.actionMap[copied.Uses]
		if !ok {
			return nil, goerr.Wrap(types.ErrActionNotFound).With("uses", copied.Uses)
		}

		logger.Info("run action", slog.Any("proc", copied))

		// Run action. If actionMock is set, use it instead of action.Run()
		if x.actionMock != nil {
			result = x.actionMock.GetResult(copied.Uses)
		} else {
			resp, err := run(ctx, x.alert, copied.Args)
			if err != nil && !copied.Force {
				return nil, types.AsActionErr(goerr.Wrap(err))
			}
			result = resp
		}
	}

	// Resolve commit attributes and refresh commit list
	copied.Commit = nil
	for _, c := range baseAction.Commit {
		resolved, err := c.ToAttr(result)
		if err != nil {
			return nil, err
		}
		if resolved == nil {
			continue
		}

		newCommit := c.Copy()
		newCommit.Attribute = *resolved
		copied.Commit = append(copied.Commit, newCommit)
	}

	actionResult := model.ActionResult{
		Action: copied,
		Result: result,
	}

	return &actionResult, nil
}
