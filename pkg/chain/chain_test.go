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
var test1AlertRego string

//go:embed testdata/test1/action.rego
var test1ActionRego string

//go:embed testdata/test1/input/scc.json
var sccData string

type mockAction struct {
	id       types.ActionID
	callback func(ctx *types.Context, alert model.Alert, params model.ActionArgs) (any, error)
}

func (x *mockAction) ID() types.ActionID {
	return x.id
}

func (x *mockAction) Run(ctx *types.Context, alert model.Alert, params model.ActionArgs) (any, error) {
	return x.callback(ctx, alert, params)
}

func TestBasic(t *testing.T) {
	var alertData any
	gt.NoError(t, json.Unmarshal([]byte(sccData), &alertData))

	alertPolicy := gt.R1(opac.NewLocal(
		opac.WithPackage("alert"),
		opac.WithRegoPrint(os.Stdout),
		opac.WithPolicyData("alert.rego", test1AlertRego),
	)).NoError(t)

	actionPolicy := gt.R1(opac.NewLocal(
		opac.WithPackage("action"),
		opac.WithPolicyData("action.rego", test1ActionRego),
	)).NoError(t)

	var called int
	mock := &mockAction{
		id: "mock",
		callback: func(ctx *types.Context, alert model.Alert, params model.ActionArgs) (any, error) {
			called++
			return nil, nil
		},
	}
	c := gt.R1(chain.New(
		chain.WithPolicyAlert(alertPolicy),
		chain.WithPolicyAction(actionPolicy),
		chain.WithAction(mock),
	)).NoError(t)

	ctx := types.NewContext()
	gt.NoError(t, c.HandleAlert(ctx, "scc", alertData))
	gt.N(t, called).Equal(1)
}

func TestDisableAction(t *testing.T) {
	var alertData any
	gt.NoError(t, json.Unmarshal([]byte(sccData), &alertData))

	alertPolicy := gt.R1(opac.NewLocal(
		opac.WithPackage("alert"),
		opac.WithPolicyData("alert.rego", test1AlertRego),
	)).NoError(t)

	actionPolicy := gt.R1(opac.NewLocal(
		opac.WithPackage("action"),
		opac.WithPolicyData("action.rego", test1ActionRego),
	)).NoError(t)

	var called int
	mock := &mockAction{
		id: "mock",
		callback: func(ctx *types.Context, alert model.Alert, params model.ActionArgs) (any, error) {
			called++
			return nil, nil
		},
	}

	c := gt.R1(chain.New(
		chain.WithPolicyAlert(alertPolicy),
		chain.WithPolicyAction(actionPolicy),
		chain.WithAction(mock),
		chain.WithDisableAction(),
	)).NoError(t)

	ctx := types.NewContext()
	gt.NoError(t, c.HandleAlert(ctx, "scc", alertData))
	gt.N(t, called).Equal(0) // Action should not be called
}

//go:embed testdata/test2/alert.rego
var test2AlertRego string

//go:embed testdata/test2/action.rego
var test2ActionRego string

//go:embed testdata/test2/action.mock.rego
var test2ActionMockRego string

func TestChainControl(t *testing.T) {
	var alertData any

	alertPolicy := gt.R1(opac.NewLocal(
		opac.WithPackage("alert"),
		opac.WithPolicyData("alert.rego", test2AlertRego),
	)).NoError(t)

	actionPolicy := gt.R1(opac.NewLocal(
		opac.WithPackage("action"),
		opac.WithPolicyData("action.rego", test2ActionRego),
		opac.WithPolicyData("action.mock.rego", test2ActionMockRego),
	)).NoError(t)

	var calledMock, calledMockAfter int
	mock := &mockAction{
		id: "mock",
		callback: func(ctx *types.Context, alert model.Alert, args model.ActionArgs) (any, error) {
			gt.A(t, alert.Params).Length(2).
				Have(model.Parameter{
					Key:   "k1",
					Value: "v1a",
				}).
				Have(model.Parameter{
					Key:   "k2",
					Value: "v2",
				})
			calledMock++
			return nil, nil
		},
	}

	mockAfter := &mockAction{
		id: "mock.after",
		callback: func(ctx *types.Context, alert model.Alert, args model.ActionArgs) (any, error) {
			gt.A(t, alert.Params).Length(3).
				Have(model.Parameter{
					Key:   "k1",
					Value: "v1a",
				}).
				Have(model.Parameter{
					Key:   "k2",
					Value: "v2a",
				}).
				Have(model.Parameter{
					Key:   "k3",
					Value: "v3",
				})

			calledMockAfter++
			return nil, nil
		},
	}

	c := gt.R1(chain.New(
		chain.WithPolicyAlert(alertPolicy),
		chain.WithPolicyAction(actionPolicy),
		chain.WithAction(mock),
		chain.WithAction(mockAfter),
	)).NoError(t)

	ctx := types.NewContext()
	gt.NoError(t, c.HandleAlert(ctx, "my_test", alertData))
	gt.N(t, calledMock).Equal(1)
	gt.N(t, calledMockAfter).Equal(1)
}

//go:embed testdata/test3/alert.rego
var test3AlertRego string

//go:embed testdata/test3/action.main.rego
var test3ActionMainRego string

//go:embed testdata/test3/action.mock.rego
var test3ActionMockRego string

func TestChainLoop(t *testing.T) {
	var alertData any

	alertPolicy := gt.R1(opac.NewLocal(
		opac.WithPackage("alert"),
		opac.WithPolicyData("alert.rego", test3AlertRego),
	)).NoError(t)

	actionPolicy := gt.R1(opac.NewLocal(
		opac.WithPackage("action"),
		opac.WithPolicyData("action.main.rego", test3ActionMainRego),
		opac.WithPolicyData("action.mock.rego", test3ActionMockRego),
	)).NoError(t)

	var calledMock int
	mock := &mockAction{
		id: "mock",
		callback: func(ctx *types.Context, alert model.Alert, args model.ActionArgs) (any, error) {
			gt.A(t, alert.Params).Length(1)
			calledMock++
			return nil, nil
		},
	}

	c := gt.R1(chain.New(
		chain.WithPolicyAlert(alertPolicy),
		chain.WithPolicyAction(actionPolicy),
		chain.WithAction(mock),
	)).NoError(t)

	ctx := types.NewContext()
	gt.NoError(t, c.HandleAlert(ctx, "my_test", alertData))
	gt.N(t, calledMock).Equal(10)
}
