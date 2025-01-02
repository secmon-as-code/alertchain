package chain

import (
	"context"

	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/service"
)

func (x *Chain) RunWorkflow(ctx context.Context, alert model.Alert, svc *service.Services) error {
	return x.runWorkflow(ctx, alert, svc)
}
