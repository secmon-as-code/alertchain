package lambda

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/m-mizutani/alertchain/pkg/chain"
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

func New(options ...Option) func(context.Context, any) (any, error) {
	cfg := &config{}

	for _, opt := range options {
		opt(cfg)
	}

	c, err := chain.New()
	if err != nil {
		utils.Logger().Error("Fail to initialize chain: %+v", err)
		panic(err)
	}
	if err := utils.ReconfigureLogger("json", "info", "-"); err != nil {
		utils.Logger().Error("Fail to initialize logger: %+v", err)
		parse(err)
	}

	return func(ctx context.Context, data any) (any, error) {
		defer func() {
			if r := recover(); r != nil {
				utils.Logger().Error("Recovered: %v", r)
			}
		}()

		event, err := parse(data)
		if err != nil {
			return "internal server error", goerr.Wrap(err, "fail to parse data")
		}

		switch v := event.(type) {
		case *events.LambdaFunctionURLRequest:
			return handleFunctionalURL(ctx, c, v)
		default:
			return "internal server error", goerr.Wrap(err, "unsupported event type")
		}
	}
}

func parse(data any) (any, error) {
	// try to parse as LambdaFunctionURLRequest
	{
		var event events.LambdaFunctionURLRequest
		if err := remapEvent(data, &event); err != nil {
			return nil, goerr.Wrap(err, "fail to remap event")
		}

		if event.RawPath != "" && event.RequestContext.HTTP.Method != "" {
			return &event, nil
		}
	}

	return data, nil
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
