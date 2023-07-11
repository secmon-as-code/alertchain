package chain

import (
	"time"

	"github.com/m-mizutani/alertchain/pkg/action"
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/memory"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/goerr"
	"golang.org/x/exp/slog"
)

type Chain struct {
	alertPolicy  *policy.Client
	actionPolicy *policy.Client
	dbClient     interfaces.Database
	timeout      time.Duration

	scenarioLogger interfaces.ScenarioLogger
	actionMock     interfaces.ActionMock
	actionMap      map[types.ActionName]interfaces.RunAction

	disableAction bool
	enablePrint   bool
	maxStackDepth int

	now func() time.Time
	env interfaces.Env
}

type Option func(c *Chain)

func New(options ...Option) (*Chain, error) {
	c := &Chain{
		dbClient:       memory.New(),
		timeout:        5 * time.Minute,
		actionMap:      action.Map(),
		scenarioLogger: &dummyScenarioLogger{},
		maxStackDepth:  types.DefaultMaxStackDepth,
		now:            time.Now,
		env:            Env,
	}

	for _, opt := range options {
		opt(c)
	}

	return c, nil
}

// HandleAlert is main function of alert chain. It receives alert data and execute actions according to the Rego policies.
func (x *Chain) HandleAlert(ctx *model.Context, schema types.Schema, data any) error {
	ctx.Logger().Info("[input] detect alert", slog.Any("data", data), slog.Any("schema", schema))
	alerts, err := x.detectAlert(ctx, schema, data)
	if err != nil {
		return goerr.Wrap(err)
	}
	ctx.Logger().Info("[output] detect alert", slog.Any("alerts", alerts))

	if x.actionPolicy == nil {
		return nil
	}

	for _, alert := range alerts {
		baseOpt := []policy.QueryOption{}
		if x.enablePrint {
			baseOpt = append(baseOpt, policy.WithRegoPrint(makeRegoPrint(ctx)))
		}

		w, err := x.newWorkflow(alert, baseOpt)
		if err != nil {
			return err
		}
		if err := w.run(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (x *Chain) detectAlert(ctx *model.Context, schema types.Schema, data any) ([]model.Alert, error) {
	if x.alertPolicy == nil {
		return nil, nil
	}

	var alertResult model.AlertPolicyResult
	opt := []policy.QueryOption{
		policy.WithPackageSuffix(string(schema)),
	}
	if x.enablePrint {
		opt = append(opt, policy.WithRegoPrint(makeRegoPrint(ctx)))
	}

	if err := x.alertPolicy.Query(ctx, data, &alertResult, opt...); err != nil {
		return nil, goerr.Wrap(err)
	}

	if len(alertResult.Alerts) == 0 {
		return nil, nil
	}

	alerts := make([]model.Alert, len(alertResult.Alerts))
	for i, meta := range alertResult.Alerts {
		alerts[i] = model.NewAlert(meta, schema, data)
	}
	return alerts, nil
}

func makeRegoPrint(ctx *model.Context) policy.RegoPrint {
	return func(file string, row int, msg string) error {
		ctx.Logger().Info("rego print",
			slog.String("file", file),
			slog.Int("row", row),
			slog.String("msg", msg),
		)
		return nil
	}
}
