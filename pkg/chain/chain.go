package chain

import (
	"log/slog"

	"github.com/m-mizutani/alertchain/pkg/chain/core"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/service"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
)

type Chain struct {
	core *core.Core
}

func New(options ...core.Option) (*Chain, error) {
	c := &Chain{
		core: core.New(options...),
	}

	return c, nil
}

// HandleAlert is main function of alert chain. It receives alert data and execute actions according to the Rego policies.
func (x *Chain) HandleAlert(ctx *model.Context, schema types.Schema, data any) ([]*model.Alert, error) {
	ctx.Logger().Info("[input] detect alert", slog.Any("data", data), slog.Any("schema", schema))
	alerts, err := x.detectAlert(ctx, schema, data)
	if err != nil {
		return nil, types.AsPolicyErr(goerr.Wrap(err))
	}
	ctx.Logger().Info("[output] detect alert", slog.Any("alerts", alerts))

	svc := service.New(x.core.DBClient())

	for _, alert := range alerts {
		record, err := svc.Workflow.Create(ctx, alert)
		if err != nil {
			return nil, err
		}

		w := newWorkflow(x.core, alert, record)

		if err := w.Run(ctx); err != nil {
			return nil, err
		}
	}

	return utils.ToPtrSlice(alerts), nil
}

func (x *Chain) detectAlert(ctx *model.Context, schema types.Schema, data any) ([]model.Alert, error) {
	var alertResult model.AlertPolicyResult
	if err := x.core.QueryAlertPolicy(ctx, schema, data, &alertResult); err != nil {
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
