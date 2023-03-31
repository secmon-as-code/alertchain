package chain_test

import (
	_ "embed"
	"encoding/json"
	"os"

	"testing"

	"github.com/m-mizutani/alertchain/pkg/chain"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/opac"
)

//go:embed testdata/test1/alert.rego
var alertRego string

//go:embed testdata/test1/action.rego
var actionRego string

//go:embed testdata/test1/scc.json
var sccData string

type mockAction struct {
	callback func(ctx *types.Context, alert model.Alert, params model.ActionParams) error
}

func (x *mockAction) ID() types.ActionID {
	return "mock"
}

func (x *mockAction) Run(ctx *types.Context, alert model.Alert, params model.ActionParams) error {
	return x.callback(ctx, alert, params)
}

func TestBasic(t *testing.T) {
	var alertData any
	gt.NoError(t, json.Unmarshal([]byte(sccData), &alertData))

	alertPolicy := gt.R1(opac.NewLocal(
		opac.WithPackage("alert"),
		opac.WithRegoPrint(os.Stdout),
		opac.WithPolicyData("alert.rego", alertRego),
	)).NoError(t)

	actionPolicy := gt.R1(opac.NewLocal(
		opac.WithPackage("action"),
		opac.WithPolicyData("action.rego", actionRego),
	)).NoError(t)

	var called int
	mock := &mockAction{
		callback: func(ctx *types.Context, alert model.Alert, params model.ActionParams) error {
			called++
			return nil
		},
	}
	c := gt.R1(chain.New(
		chain.WithPolicyAlert(alertPolicy),
		chain.WithPolicyAction(actionPolicy),
		chain.WithAction(mock),
	)).NoError(t)

	ctx := types.NewContext()
	gt.NoError(t, c.HandleAlert(ctx, "test", alertData))
	gt.N(t, called).Equal(1)
}
