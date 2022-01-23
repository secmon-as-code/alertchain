package alert_test

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/service/alert"
)

func TestAlertService(t *testing.T) {
	clients := infra.New(db.NewDBMock(t), nil)
	newAlert := model.NewAlert(&model.Alert{Title: "testing", Detector: "blue"})
	svc := alert.New(newAlert, clients)

	t.Skip("TODO: implement")
	svc.HandleChangeRequest(types.NewContext(), &model.ChangeRequest{
		NewReferences: []*model.Reference{},
	})
}
