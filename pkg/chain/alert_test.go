package chain_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/chain"
	"github.com/m-mizutani/alertchain/pkg/chain/core"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/gt"
)

func TestAlertRaw(t *testing.T) {
	alertData := `{"color": "blue"}`
	var alertDataPP bytes.Buffer
	enc := json.NewEncoder(&alertDataPP)
	enc.SetIndent("", "  ")
	gt.NoError(t, enc.Encode(alertData))

	alertPolicy := gt.R1(policy.New(
		policy.WithPackage("alert"),
		policy.WithFile("testdata/alert_feature/alert.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	actionPolicy := gt.R1(policy.New(
		policy.WithPackage("action"),
		policy.WithFile("testdata/alert_feature/action.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	var calledMock int
	mock := func(ctx *model.Context, alert model.Alert, args model.ActionArgs) (any, error) {
		s := gt.Cast[string](t, args["raw"])
		gt.V(t, s).Equal(alertDataPP.String())
		calledMock++
		return nil, nil
	}

	c := gt.R1(chain.New(
		core.WithPolicyAlert(alertPolicy),
		core.WithPolicyAction(actionPolicy),
		core.WithExtraAction("test.output_raw", mock),
	)).NoError(t)

	ctx := model.NewContext()
	gt.R1(c.HandleAlert(ctx, "amber", alertData)).NoError(t)
	gt.N(t, calledMock).Equal(1)
}
