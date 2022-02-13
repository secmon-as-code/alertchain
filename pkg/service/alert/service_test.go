package alert_test

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/db/spanner"
	"github.com/m-mizutani/alertchain/pkg/service/alert"
)

func TestAlertService(t *testing.T) {
	clients := infra.New(spanner.NewTestDB(t), nil)
	newAlert := model.NewAlert(&model.Alert{Title: "testing", Detector: "blue"})
	_ = alert.New(newAlert, clients)

	t.Skip("TODO: implement")
}
