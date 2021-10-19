package alertchain

import (
	"context"

	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

var logger = utils.Logger

type Chain struct {
	Jobs    Jobs
	Sources []Source
	Actions Actions

	db db.Interface
}

func New(dbClient db.Interface) *Chain {
	return &Chain{
		db: dbClient,
	}
}

type Action interface {
	Name() string
	Executable(attr *Attribute) bool
	Execute(ctx *types.Context, attr *Attribute) error
}

type Actions []Action

func (x *Chain) diagnosis() error {
	if x.db == nil {
		return goerr.Wrap(types.ErrChainIsNotConfigured, "DB is not set")
	}
	return nil
}

func (x *Chain) Execute(ctx context.Context, alert *Alert) (*Alert, error) {
	if err := x.diagnosis(); err != nil {
		return nil, err
	}
	if err := alert.validate(); err != nil {
		return nil, err
	}

	c, ok := ctx.(*types.Context)
	if !ok {
		c = types.NewContextWith(ctx)
	}

	logger.With("alert", alert).Trace("Starting Chain.Execute")
	alertID, err := insertAlert(c, alert, x.db)
	if err != nil {
		return nil, err
	}

	logger.With("alert", alert).Trace("Exiting Chain.Execute")

	if err := x.Jobs.Execute(c, x.db, alertID); err != nil {
		return nil, err
	}

	created, err := x.db.GetAlert(c, alertID)
	if err != nil {
		return nil, err
	}

	return newAlert(created), nil
}

func (x *Chain) AddJob(job *Job) {
	x.Jobs = append(x.Jobs, job)
}

func (x *Chain) AddJobs(jobs ...*Job) {
	x.Jobs = append(x.Jobs, jobs...)
}
