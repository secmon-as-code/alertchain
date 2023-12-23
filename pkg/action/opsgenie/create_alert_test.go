package opsgenie_test

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/action/opsgenie"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/gt"
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
)

func TestOpsgenie(t *testing.T) {
	var (
		apiKey        string
		responderID   string
		responderType string
	)

	if err := utils.LoadEnv(
		utils.Env("TEST_OPSGENIE_API_KEY", &apiKey),
		utils.Env("TEST_OPSGENIE_RESPONDER_ID", &responderID),
		utils.Env("TEST_OPSGENIE_RESPONDER_TYPE", &responderType),
	); err != nil {
		t.Skipf("Skip test due to missing env: %v", err)
	}

	t.Run("Create alert", func(t *testing.T) {
		input := model.NewAlert(model.AlertMetaData{
			Title:       "test_alert",
			Description: "test_description",
			Source:      "test_source",
			Attrs: model.Attributes{
				{
					Key:   "key1",
					Value: "val1",
				},
			},
		}, "test_alert", struct{}{})
		ctx := model.NewContext()
		args := model.ActionArgs{
			"secret_api_key": apiKey,
			"responders": []opsgenie.Responder{
				{
					ID:   responderID,
					Type: responderType,
				},
			},
		}

		ret := gt.R1(opsgenie.CreateAlert(ctx, input, args)).NoError(t)
		resp := gt.Cast[*alert.AsyncAlertResult](t, ret)
		gt.NotEqual(t, resp.RequestId, "")
	})
}
