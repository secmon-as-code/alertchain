package chain_test

import (
	"embed"
	"encoding/json"
	"path/filepath"

	"testing"

	"github.com/m-mizutani/alertchain/pkg/chain"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/infra/logger"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/gt"
)

//go:embed testdata/**
var testDataFS embed.FS

func read(path string) ([]byte, error) {
	return testDataFS.ReadFile(filepath.Clean(path))
}

func TestBasic(t *testing.T) {
	var alertData any
	sccData := gt.R1(read("testdata/basic/input/scc.json")).NoError(t)
	gt.NoError(t, json.Unmarshal([]byte(sccData), &alertData))

	alertPolicy := gt.R1(policy.New(
		policy.WithPackage("alert"),
		policy.WithFile("testdata/basic/alert.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	actionPolicy := gt.R1(policy.New(
		policy.WithPackage("action"),
		policy.WithFile("testdata/basic/action.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	var called int
	mock := func(ctx *model.Context, _ model.Alert, _ model.ActionArgs) (any, error) {
		called++
		return nil, nil
	}

	c := gt.R1(chain.New(
		chain.WithPolicyAlert(alertPolicy),
		chain.WithPolicyAction(actionPolicy),
		chain.WithExtraAction("mock", mock),
	)).NoError(t)

	ctx := model.NewContext()
	gt.NoError(t, c.HandleAlert(ctx, "scc", alertData))
	gt.N(t, called).Equal(1)
}

func TestDisableAction(t *testing.T) {
	var alertData any
	sccData := gt.R1(read("testdata/basic/input/scc.json")).NoError(t)
	gt.NoError(t, json.Unmarshal([]byte(sccData), &alertData))

	alertPolicy := gt.R1(policy.New(
		policy.WithPackage("alert"),
		policy.WithFile("testdata/basic/alert.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	actionPolicy := gt.R1(policy.New(
		policy.WithPackage("action"),
		policy.WithFile("testdata/basic/action.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	var called int
	mock := func(ctx *model.Context, _ model.Alert, _ model.ActionArgs) (any, error) {
		called++
		return nil, nil
	}

	c := gt.R1(chain.New(
		chain.WithPolicyAlert(alertPolicy),
		chain.WithPolicyAction(actionPolicy),
		chain.WithExtraAction("mock", mock),
		chain.WithDisableAction(),
	)).NoError(t)

	ctx := model.NewContext()
	gt.NoError(t, c.HandleAlert(ctx, "scc", alertData))
	gt.N(t, called).Equal(0) // Action should not be called
}

func TestChainControl(t *testing.T) {
	var alertData any

	alertPolicy := gt.R1(policy.New(
		policy.WithPackage("alert"),
		policy.WithFile("testdata/control/alert.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	actionPolicy := gt.R1(policy.New(
		policy.WithPackage("action"),
		policy.WithFile("testdata/control/action.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	var calledMock, calledMockAfter int
	mock := func(ctx *model.Context, alert model.Alert, _ model.ActionArgs) (any, error) {
		gt.A(t, alert.Params).Length(2).
			At(0, func(t testing.TB, v model.Parameter) {
				gt.V(t, v.Name).Equal("k1")
				gt.V(t, v.Value).Equal("v1")
			}).
			At(1, func(t testing.TB, v model.Parameter) {
				gt.V(t, v.Name).Equal("k2")
				gt.V(t, v.Value).Equal("v2")
			})
		calledMock++

		return nil, nil
	}

	mockAfter := func(ctx *model.Context, alert model.Alert, _ model.ActionArgs) (any, error) {
		gt.A(t, alert.Params).Length(3).
			At(0, func(t testing.TB, v model.Parameter) {
				gt.V(t, v.Name).Equal("k1")
				gt.V(t, v.Value).Equal("v1")
			}).
			At(1, func(t testing.TB, v model.Parameter) {
				gt.V(t, v.Name).Equal("k2")
				gt.V(t, v.Value).Equal("v2a")
			}).
			At(2, func(t testing.TB, v model.Parameter) {
				gt.V(t, v.Name).Equal("k3")
				gt.V(t, v.Value).Equal("v3")
			})

		calledMockAfter++
		return nil, nil
	}

	c := gt.R1(chain.New(
		chain.WithPolicyAlert(alertPolicy),
		chain.WithPolicyAction(actionPolicy),
		chain.WithExtraAction("mock", mock),
		chain.WithExtraAction("mock.after", mockAfter),
	)).NoError(t)

	ctx := model.NewContext()
	gt.NoError(t, c.HandleAlert(ctx, "my_test", alertData))
	gt.N(t, calledMock).Equal(1)
	gt.N(t, calledMockAfter).Equal(1)
}

func TestChainLoop(t *testing.T) {
	var alertData any

	alertPolicy := gt.R1(policy.New(
		policy.WithPackage("alert"),
		policy.WithFile("testdata/loop/alert.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	actionPolicy := gt.R1(policy.New(
		policy.WithPackage("action"),
		policy.WithFile("testdata/loop/action.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	var calledMock int
	mock := func(ctx *model.Context, alert model.Alert, _ model.ActionArgs) (any, error) {
		gt.A(t, alert.Params).Length(1)
		calledMock++
		return nil, nil
	}

	c := gt.R1(chain.New(
		chain.WithPolicyAlert(alertPolicy),
		chain.WithPolicyAction(actionPolicy),
		chain.WithExtraAction("mock", mock),
	)).NoError(t)

	ctx := model.NewContext()
	gt.NoError(t, c.HandleAlert(ctx, "my_test", alertData))
	gt.N(t, calledMock).Equal(9)
}

func TestPlaybook(t *testing.T) {
	alertPolicy := gt.R1(policy.New(
		policy.WithPackage("alert"),
		policy.WithFile("testdata/play/alert.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	actionPolicy := gt.R1(policy.New(
		policy.WithPackage("action"),
		policy.WithFile("testdata/play/action.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	var playbook model.Playbook
	gt.NoError(t, model.ParsePlaybook("testdata/play/playbook.jsonnet", read, &playbook))
	gt.A(t, playbook.Scenarios).Length(1).At(0, func(t testing.TB, v *model.Scenario) {
		gt.V(t, v.ID).Equal("s1")
		event := gt.Cast[map[string]any](t, v.Event)
		gt.V(t, event).Equal(map[string]any{
			"class": "threat",
		})
	})

	var calledMock int
	mock := func(ctx *model.Context, alert model.Alert, _ model.ActionArgs) (any, error) {
		gt.A(t, alert.Params).Length(1)
		calledMock++
		return nil, nil
	}

	recorder := logger.NewMemory(playbook.Scenarios[0])
	c := gt.R1(chain.New(
		chain.WithPolicyAlert(alertPolicy),
		chain.WithPolicyAction(actionPolicy),
		chain.WithExtraAction("mock", mock),
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
			gt.V(t, v.Name).Equal("c")

			// Value has been converted to float64 by encoding/decoding json
			gt.V(t, v.Value).Equal(float64(1))
		})
	})
}
