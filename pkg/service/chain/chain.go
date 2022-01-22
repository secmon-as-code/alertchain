package chain

import (
	"sync"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra"
)

type Chain struct {
	clients *infra.Clients
	jobs    []*model.Job
}

func New(clients *infra.Clients, actions []model.ActionDefinition, jobs []model.JobDefinition) (*Chain, error) {
	return &Chain{}, nil
}

func invokeAction(
	ctx *types.Context,
	action model.Action,
	alert *model.Alert,
	wg *sync.WaitGroup,
	reqCh chan *model.ChangeRequest,
	errCh chan error,
) {
	defer wg.Done()
	req, err := action.Run(ctx, alert)

	if err != nil {
		errCh <- err
		return
	}

	if req != nil {
		reqCh <- req
	}
}

func (x *Chain) Invoke(ctx *types.Context, alert *model.Alert) error {
	for _, job := range x.jobs {
		wg := &sync.WaitGroup{}
		errCh := make(chan error)
		reqCh := make(chan *model.ChangeRequest)

		for i := range job.Actions {
			wg.Add(1)
			go invokeAction(ctx, job.Actions[i], alert, wg, reqCh, errCh)
		}

		wg.Wait()
		close(errCh)

		var errs []error
		for err := range errCh {
			errs = append(errs, err)
		}

		if len(errs) > 0 && job.ExitOnErr {
			return errs[0]
		}
	}

	return nil
}
