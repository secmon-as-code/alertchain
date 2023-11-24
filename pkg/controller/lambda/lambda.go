package lambda

import (
	"context"
	"encoding/json"
	"os"

	"log/slog"

	"github.com/m-mizutani/alertchain/pkg/chain"
	"github.com/m-mizutani/alertchain/pkg/controller/cli/flag"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
)

type config struct {
	alertPolicyDir  string
	actionPolicyDir string

	readFile func(string) ([]byte, error)
}

type Option func(*config)

func WithAlertPolicyDir(dir string) Option {
	return func(cfg *config) {
		cfg.alertPolicyDir = dir
	}
}

func WithActionPolicyDir(dir string) Option {
	return func(cfg *config) {
		cfg.actionPolicyDir = dir
	}
}

func WithReadFile(f func(string) ([]byte, error)) Option {
	return func(cfg *config) {
		cfg.readFile = f
	}
}

type Callback func(context.Context, types.Schema, any) error

type Handler func(ctx context.Context, event any, cb Callback) error

func New(handler Handler, options ...Option) func(context.Context, any) (any, error) {
	var cfg config

	for _, opt := range options {
		opt(&cfg)
	}

	c, err := chain.New()
	if err != nil {
		utils.Logger().Error("Fail to initialize chain: %+v", err)
		panic(err)
	}

	utils.ReconfigureLogger(os.Stdout, slog.LevelInfo, flag.LogFormatJSON)

	return func(ctx context.Context, data any) (any, error) {
		defer func() {
			if r := recover(); r != nil {
				utils.Logger().Error("Recovered: %v", r)
			}
		}()

		callback := func(ctx context.Context, schema types.Schema, data any) error {
			return c.HandleAlert(model.NewContext(model.WithBase(ctx)), schema, data)
		}

		if err := handler(ctx, data, callback); err != nil {
			return nil, goerr.Wrap(err, "fail to handle alert")
		}

		return "OK", nil
	}
}

func remapEvent(src, dst any) error {
	raw, err := json.Marshal(src)
	if err != nil {
		return goerr.Wrap(err, "fail to marshal event")
	}

	if err := json.Unmarshal(raw, dst); err != nil {
		return goerr.Wrap(err, "fail to unmarshal event")
	}
	return nil
}
