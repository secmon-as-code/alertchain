package alertchain

import (
	"context"
	"sync"
	"time"

	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/zlog"
)

type Chain struct {
	jobs    Jobs
	sources []Source
	actions Actions

	config types.Config

	db     db.Interface
	api    *apiServer
	logger *zlog.Logger
}

func New(options ...Option) (*Chain, error) {
	chain := &Chain{
		logger: zlog.New(),
		config: types.Config{
			DB: types.DBConfig{
				Type:   "sqlite3",
				Config: "file:alertchain?mode=memory&cache=shared&_fk=1",
			},
		},
	}

	for _, opt := range options {
		opt(chain)
	}

	if chain.db == nil {
		dbClient, err := db.New(chain.config.DB.Type, chain.config.DB.Config)
		if err != nil {
			return nil, goerr.Wrap(types.ErrChainIsNotConfigured, "DB is not set")
		}
		chain.db = dbClient
	}

	return chain, nil
}

type Action interface {
	Name() string
	Executable(attr *Attribute) bool
	Execute(ctx *types.Context, attr *Attribute) error
}

type Actions []Action

func (x *Chain) Execute(ctx context.Context, alert *Alert) (*Alert, error) {
	if err := alert.validate(); err != nil {
		return nil, err
	}

	c, ok := ctx.(*types.Context)
	if !ok {
		c = types.NewContextWith(ctx, x.logger)
	}

	x.logger.With("alert", alert).Trace("Starting Chain.Execute")
	alertID, err := insertAlert(c, alert, x.db)
	if err != nil {
		return nil, err
	}

	x.logger.With("alert", alert).Trace("Exiting Chain.Execute")

	if err := x.jobs.Execute(c, x.db, alertID); err != nil {
		return nil, err
	}

	created, err := x.db.GetAlert(c, alertID)
	if err != nil {
		return nil, err
	}

	return newAlert(created), nil
}

func (x *Chain) handleError(err error) {
	x.logger.Err(err).Error("failed run")
}

func (x *Chain) Start() error {
	x.logger.With("config", x.config).Info("Starting AlertChain")

	handler := func(ctx context.Context, alert *Alert) error {
		_, err := x.Execute(ctx, alert)
		return err
	}

	var wg sync.WaitGroup
	for i := range x.sources {
		wg.Add(1)

		go func(src Source) {
			defer wg.Done()

			for {
				if err := src.Run(handler); err != nil {
					x.handleError(err)
				} else {
					break
				}

				time.Sleep(time.Second * 3)
			}
		}(x.sources[i])
	}

	if x.api != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if err := x.api.Run(); err != nil {
					x.handleError(err)
				} else {
					break
				}
			}
		}()
	}

	wg.Wait()

	return nil
}

func (x *Chain) Logger() *zlog.Logger {
	return x.logger
}
