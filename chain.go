package alertchain

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/zlog"
)

type Chain struct {
	jobs    Jobs
	sources []Source
	actions Actions

	config      types.Config
	configMutex sync.Mutex
	configured  int32

	db     db.Interface
	logger *zlog.Logger
}

func New(options ...Option) *Chain {
	chain := newDefault()

	for _, opt := range options {
		opt(chain)
	}

	return chain
}

func newDefault() *Chain {
	return &Chain{
		logger: zlog.New(),
		config: types.Config{
			DB: types.DBConfig{
				Type:   "sqlite3",
				Config: "file:alertchain?mode=memory&cache=shared&_fk=1",
			},
		},
	}
}

type Action interface {
	Name() string
	Executable(attr *Attribute) bool
	Execute(ctx *types.Context, attr *Attribute) error
}

type Actions []Action

func (x *Chain) init() error {
	if atomic.LoadInt32(&(x.configured)) > 0 {
		return nil
	}

	x.configMutex.Lock()
	defer x.configMutex.Unlock()

	if atomic.LoadInt32(&(x.configured)) > 0 {
		return nil
	}

	if x.db == nil {
		dbClient, err := db.New(x.config.DB.Type, x.config.DB.Config)
		if err != nil {
			return goerr.Wrap(types.ErrChainIsNotConfigured, "DB is not set")
		}
		x.db = dbClient
	}

	atomic.StoreInt32(&(x.configured), int32(1))

	return nil
}

func (x *Chain) Execute(ctx context.Context, alert *Alert) (*Alert, error) {
	if err := x.init(); err != nil {
		return nil, err
	}

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

func (x *Chain) Start() error {
	x.logger.With("config", x.config).Info("Starting AlertChain")

	x.StartSources()
	return nil
}
