package chain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/alertchain/action"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/interfaces"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/infra/memory"
	"github.com/secmon-lab/alertchain/pkg/infra/policy"
	"github.com/secmon-lab/alertchain/pkg/service"
	"github.com/secmon-lab/alertchain/pkg/utils"
)

type Chain struct {
	alertPolicy  *policy.Client
	actionPolicy *policy.Client
	dbClient     interfaces.Database

	recorder   interfaces.ScenarioRecorder
	actionMock interfaces.ActionMock
	actionMap  map[types.ActionName]model.RunAction

	timeout      time.Duration
	enablePrint  bool
	maxSequences int

	now func() time.Time
	env interfaces.Env
}

func New(options ...Option) (*Chain, error) {
	c := &Chain{
		dbClient:     memory.New(),
		timeout:      5 * time.Minute,
		actionMap:    action.Map(),
		actionMock:   nil,
		recorder:     &dummyScenarioRecorder{},
		maxSequences: types.DefaultMaxSequences,
		now:          time.Now,
		env:          utils.Env,
	}

	for _, opt := range options {
		opt(c)
	}

	return c, nil
}

type Option func(c *Chain)

func WithPolicyAlert(p *policy.Client) Option {
	return func(c *Chain) {
		c.alertPolicy = p
	}
}

func WithPolicyAction(p *policy.Client) Option {
	return func(c *Chain) {
		c.actionPolicy = p
	}
}

func WithEnablePrint() Option {
	return func(c *Chain) {
		c.enablePrint = true
	}
}

func WithExtraAction(name types.ActionName, action model.RunAction) Option {
	return func(c *Chain) {
		if _, ok := c.actionMap[name]; ok {
			panic("action name is already registered: " + name)
		}
		c.actionMap[name] = action
	}
}

func WithActionMock(mock interfaces.ActionMock) Option {
	return func(c *Chain) {
		c.actionMock = mock
	}
}

func WithScenarioRecorder(logger interfaces.ScenarioRecorder) Option {
	return func(c *Chain) {
		c.recorder = logger
	}
}

func WithEnv(f interfaces.Env) Option {
	return func(c *Chain) {
		c.env = f
	}
}

func WithDatabase(db interfaces.Database) Option {
	return func(c *Chain) {
		c.dbClient = db
	}
}

// HandleAlert is main function of alert chain. It receives alert data and execute actions according to the Rego policies.
func (x *Chain) HandleAlert(ctx context.Context, schema types.Schema, data any) ([]*model.Alert, error) {
	logger := ctxutil.Logger(ctx)
	logger.Debug("[input] detect alert", slog.Any("data", data), slog.Any("schema", schema))

	var alertResult model.AlertPolicyResult
	if err := x.queryAlertPolicy(ctx, schema, data, &alertResult); err != nil {
		return nil, goerr.Wrap(err)
	}

	if len(alertResult.Alerts) == 0 {
		return nil, nil
	}

	alerts := make([]model.Alert, len(alertResult.Alerts))
	for i, meta := range alertResult.Alerts {
		alerts[i] = model.NewAlert(meta, schema, data)
	}

	logger.Debug("[output] detect alert", slog.Any("alerts", alerts))

	svc := service.New(x.dbClient)

	for _, alert := range alerts {
		newCtx := ctxutil.InjectLogger(ctx, logger.With("alert_id", alert.ID))
		if err := x.runWorkflow(newCtx, alert, svc); err != nil {
			return nil, err
		}
	}

	return utils.ToPtrSlice(alerts), nil
}

func (x *Chain) queryAlertPolicy(ctx context.Context, schema types.Schema, in, out any) error {
	if x.alertPolicy == nil {
		return nil
	}

	options := []policy.QueryOption{
		policy.WithPackageSuffix(string(schema)),
	}
	if x.enablePrint {
		options = append(options, policy.WithRegoPrint(makeRegoPrint(ctx)))
	}

	if err := x.alertPolicy.Query(ctx, in, out, options...); err != nil && !errors.Is(err, types.ErrNoPolicyResult) {
		return types.AsPolicyErr(goerr.Wrap(err, "failed to evaluate alert policy").With("request", in))
	}
	ctxutil.Logger(ctx).Debug("queried action policy", slog.Any("in", in), slog.Any("out", out))

	return nil
}

func (x *Chain) queryActionPolicy(ctx context.Context, in, out any) error {
	if x.actionPolicy == nil {
		return nil
	}

	var options []policy.QueryOption
	if x.enablePrint {
		options = append(options, policy.WithRegoPrint(makeRegoPrint(ctx)))
	}

	if err := x.actionPolicy.Query(ctx, in, out, options...); err != nil && !errors.Is(err, types.ErrNoPolicyResult) {
		return types.AsPolicyErr(goerr.Wrap(err, "failed to evaluate action policy").With("request", in))
	}
	ctxutil.Logger(ctx).Debug("queried action policy", slog.Any("in", in), slog.Any("out", out))

	return nil
}

func makeRegoPrint(ctx context.Context) policy.RegoPrint {
	return func(file string, row int, msg string) error {
		if ctxutil.IsCLI(ctx) {
			fmt.Printf("	%s:%d: %s\n", file, row, msg)
		} else {
			ctxutil.Logger(ctx).Info("rego print",
				slog.String("file", file),
				slog.Int("row", row),
				slog.String("msg", msg),
			)
		}
		return nil
	}
}
