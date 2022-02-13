package chain

import (
	"sync"
	"time"

	"github.com/m-mizutani/alertchain/pkg/actions"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/service/alert"
	"github.com/m-mizutani/goerr"
)

type Chain struct {
	clients *infra.Clients
	actions map[string]model.Action
	jobs    []*model.Job
}

func New(clients *infra.Clients, actionDefs []model.ActionDefinition, jobDefs []model.JobDefinition) (*Chain, error) {
	if clients == nil {
		return nil, goerr.Wrap(types.ErrInvalidChainConfig, "clients is not set")
	}

	chain := &Chain{
		clients: clients,
	}

	for _, def := range actionDefs {
		if exist, ok := chain.actions[def.ID]; ok {
			return nil, goerr.Wrap(types.ErrDuplicatedActionID).With("exist", exist).With("dupped", def)
		}

		action, err := actions.New(def.Use, def.Config)
		if err != nil {
			return nil, err
		}
		chain.actions[def.ID] = action
	}

	for _, def := range jobDefs {
		job := def.Job
		if def.Timeout != "" {
			timeout, err := time.ParseDuration(def.Timeout)
			if err != nil {
				return nil, types.ErrInvalidChainConfig.Wrap(err).With("job", def).With("invalid_value", def.Timeout)
			}

			job.Timeout = timeout
		}

		for _, actionID := range def.Actions {
			action, ok := chain.actions[actionID]
			if !ok {
				return nil, goerr.Wrap(types.ErrActionNotDefined).With("action", action)
			}
			job.Actions = append(job.Actions, action)
		}

		chain.jobs = append(chain.jobs, &job)
	}

	return chain, nil
}

func invokeAction(
	ctx *types.Context,
	action model.Action,
	alert *model.Alert,
	args []*model.Attribute,
	wg *sync.WaitGroup,
	reqCh chan *model.ChangeRequest,
	errCh chan error,
) {
	defer wg.Done()
	req, err := action.Run(ctx, alert, args...)

	if err != nil {
		errCh <- err
		return
	}

	if req != nil {
		reqCh <- req
	}
}

func (x *Chain) Invoke(ctx *types.Context, target *model.Alert) error {
	alertSvc := alert.New(target, x.clients)

	for _, job := range x.jobs {
		wg := &sync.WaitGroup{}
		errCh := make(chan error)
		reqCh := make(chan *model.ChangeRequest)

		for i := range job.Actions {
			input := &model.ActionInquiryInput{
				Alert:  alertSvc.Alert(),
				Job:    job,
				Action: &job.Actions[i],
			}
			var result model.ActionInquiryResult
			if err := x.clients.Policy().Eval(ctx, input, &result); err != nil {
				return goerr.Wrap(err).With("input", input)
			}

			if result.Cancel {
				continue
			}

			wg.Add(1)
			go invokeAction(ctx, job.Actions[i], alertSvc.Alert(), result.Args, wg, reqCh, errCh)
		}

		wg.Wait()
		close(errCh)
		close(reqCh)

		var errs []error
		for err := range errCh {
			errs = append(errs, err)
		}

		if len(errs) > 0 && job.ExitOnErr {
			return errs[0]
		}

		for req := range reqCh {
			if err := alertSvc.HandleChangeRequest(ctx, req); err != nil {
				return goerr.Wrap(err).With("changeRequest", req)
			}
		}
		if err := alertSvc.Refresh(ctx); err != nil {
			return err
		}

		input := &model.AlertInquiryInput{
			Alert: alertSvc.Alert(),
		}
		var result model.AlertInquiryResult
		if err := x.clients.Policy().Eval(ctx, input, &result); err != nil {
			return err
		}

		var req model.ChangeRequest
		if result.Severity != "" {
			req.UpdateSeverity(types.Severity(result.Severity))
		}
		if err := alertSvc.HandleChangeRequest(ctx, &req); err != nil {
			return goerr.Wrap(err).With("changeRequest", req)
		}
	}

	return nil
}
