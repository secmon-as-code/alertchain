package chain_test

import (
	"embed"
	_ "embed"
	"encoding/json"
	"os"

	"testing"

	"github.com/m-mizutani/alertchain/pkg/chain"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/logger"
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
	callback func(ctx *model.Context, alert model.Alert, params model.ActionArgs) (any, error)
}

func (x *mockAction) ID() types.ActionID {
	return x.id
}

func (x *mockAction) Run(ctx *model.Context, alert model.Alert, params model.ActionArgs) (any, error) {
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
		callback: func(ctx *model.Context, alert model.Alert, params model.ActionArgs) (any, error) {
			called++
			return nil, nil
		},
	}
	c := gt.R1(chain.New(
		chain.WithPolicyAlert(alertPolicy),
		chain.WithPolicyAction(actionPolicy),
		chain.WithAction(mock),
	)).NoError(t)

	ctx := model.NewContext()
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
		callback: func(ctx *model.Context, alert model.Alert, params model.ActionArgs) (any, error) {
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

	ctx := model.NewContext()
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
		callback: func(ctx *model.Context, alert model.Alert, args model.ActionArgs) (any, error) {
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
		callback: func(ctx *model.Context, alert model.Alert, args model.ActionArgs) (any, error) {
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

	ctx := model.NewContext()
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
		callback: func(ctx *model.Context, alert model.Alert, args model.ActionArgs) (any, error) {
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

	ctx := model.NewContext()
	gt.NoError(t, c.HandleAlert(ctx, "my_test", alertData))
	gt.N(t, calledMock).Equal(10)
}

//go:embed testdata/play/alert.rego
var testPlayAlertRego string

//go:embed testdata/play/action.main.rego
var testPlayActionMainRego string

//go:embed testdata/play/action.mock.rego
var testPlayActionMockRego string

//go:embed testdata/play/*.jsonnet
var testPlaybookFS embed.FS

func TestPlaybook(t *testing.T) {
	alertPolicy := gt.R1(opac.NewLocal(
		opac.WithPackage("alert"),
		opac.WithPolicyData("alert.rego", testPlayAlertRego),
	)).NoError(t)

	actionPolicy := gt.R1(opac.NewLocal(
		opac.WithPackage("action"),
		opac.WithPolicyData("action.main.rego", testPlayActionMainRego),
		opac.WithPolicyData("action.mock.rego", testPlayActionMockRego),
	)).NoError(t)

	var playbook model.Playbook
	gt.NoError(t, model.ParsePlaybook("testdata/play/playbook.jsonnet", testPlaybookFS.ReadFile, &playbook))
	gt.A(t, playbook.Scenarios).Length(1).At(0, func(t testing.TB, v *model.Scenario) {
		gt.V(t, v.ID).Equal("s1")
		alert := gt.Cast[map[string]any](t, v.Alert)
		gt.V(t, alert).Equal(map[string]any{
			"class": "threat",
		})
	})

	var calledMock int
	mock := &mockAction{
		id: "my_action",
		callback: func(ctx *model.Context, alert model.Alert, args model.ActionArgs) (any, error) {
			gt.A(t, alert.Params).Length(1)
			calledMock++
			return nil, nil
		},
	}

	recorder := logger.NewMemory(playbook.Scenarios[0])
	c := gt.R1(chain.New(
		chain.WithPolicyAlert(alertPolicy),
		chain.WithPolicyAction(actionPolicy),
		chain.WithAction(mock),
		chain.WithActionMock(playbook.Scenarios[0]),
		chain.WithScenarioLogger(recorder),
	)).NoError(t)

	var alertData any
	ctx := model.NewContext()
	gt.NoError(t, c.HandleAlert(ctx, "my_test", alertData))
	gt.N(t, calledMock).Equal(0)

	gt.V(t, recorder.Log.ID).Equal("s1")
	gt.V(t, recorder.Log.Title).Equal("Scenario 1")
	gt.A(t, recorder.Log.AlertLog).Length(1).At(0, func(t testing.TB, v *model.AlertLog) {
		gt.V(t, v.Alert.Title).Equal("test alert")
		gt.A(t, v.Alert.Params).Length(1).At(0, func(t testing.TB, v model.Parameter) {
			gt.V(t, v.Key).Equal("c")
			gt.V(t, v.Value).Equal(float64(1))
		})
	})
}
