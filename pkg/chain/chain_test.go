package chain_test

import (
	"encoding/json"
	"sync"

	"testing"

	"github.com/m-mizutani/alertchain/pkg/chain"
	"github.com/m-mizutani/alertchain/pkg/chain/core"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/infra/logging"
	"github.com/m-mizutani/alertchain/pkg/infra/memory"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/gt"
)

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
	mock := func(ctx *model.Context, _ model.Alert, args model.ActionArgs) (any, error) {
		s, ok := args["s"].(string)
		gt.B(t, ok).True()
		gt.V(t, s).Equal("blue")

		n, ok := args["n"].(float64)
		gt.B(t, ok).True()
		gt.V(t, n).Equal(5)
		called++
		return nil, nil
	}

	c := gt.R1(chain.New(
		core.WithPolicyAlert(alertPolicy),
		core.WithPolicyAction(actionPolicy),
		core.WithExtraAction("mock", mock),
	)).NoError(t)

	ctx := model.NewContext()
	gt.R1(c.HandleAlert(ctx, "scc", alertData)).NoError(t)
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
		core.WithPolicyAlert(alertPolicy),
		core.WithPolicyAction(actionPolicy),
		core.WithExtraAction("mock", mock),
		core.WithDisableAction(),
	)).NoError(t)

	ctx := model.NewContext()
	gt.R1(c.HandleAlert(ctx, "scc", alertData)).NoError(t)
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
		gt.A(t, alert.Attrs).Length(2).
			At(0, func(t testing.TB, v model.Attribute) {
				gt.V(t, v.Key).Equal("k1")
				gt.V(t, v.Value).Equal("v1")
			}).
			At(1, func(t testing.TB, v model.Attribute) {
				gt.V(t, v.Key).Equal("k2")
				gt.V(t, v.Value).Equal("v2")
			})
		calledMock++

		return nil, nil
	}

	mockAfter := func(ctx *model.Context, alert model.Alert, _ model.ActionArgs) (any, error) {
		gt.A(t, alert.Attrs).Length(3).
			At(0, func(t testing.TB, v model.Attribute) {
				gt.V(t, v.Key).Equal("k1")
				gt.V(t, v.Value).Equal("v1")
			}).
			At(1, func(t testing.TB, v model.Attribute) {
				gt.V(t, v.Key).Equal("k2")
				gt.V(t, v.Value).Equal("v2a")
			}).
			At(2, func(t testing.TB, v model.Attribute) {
				gt.V(t, v.Key).Equal("k3")
				gt.V(t, v.Value).Equal("v3")
			})

		calledMockAfter++
		return nil, nil
	}

	c := gt.R1(chain.New(
		core.WithPolicyAlert(alertPolicy),
		core.WithPolicyAction(actionPolicy),
		core.WithExtraAction("mock", mock),
		core.WithExtraAction("mock.after", mockAfter),
	)).NoError(t)

	ctx := model.NewContext()
	gt.R1(c.HandleAlert(ctx, "my_test", alertData)).NoError(t)
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
		gt.A(t, alert.Attrs).Length(1)
		calledMock++
		return nil, nil
	}

	c := gt.R1(chain.New(
		core.WithPolicyAlert(alertPolicy),
		core.WithPolicyAction(actionPolicy),
		core.WithExtraAction("mock", mock),
	)).NoError(t)

	ctx := model.NewContext()
	gt.R1(c.HandleAlert(ctx, "my_test", alertData)).NoError(t)
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

		for _, event := range v.Events {
			gt.V(t, event.Schema).Equal("my_test")
			gt.V(t, event.Input).Equal(map[string]any{
				"class": "threat",
			})
		}
	})

	var calledMock int
	mock := func(ctx *model.Context, alert model.Alert, _ model.ActionArgs) (any, error) {
		gt.A(t, alert.Attrs).Length(1)
		calledMock++
		return nil, nil
	}

	recorder := logging.NewMemory(playbook.Scenarios[0])
	c := gt.R1(chain.New(
		core.WithPolicyAlert(alertPolicy),
		core.WithPolicyAction(actionPolicy),
		core.WithExtraAction("mock", mock),
		core.WithActionMock(&playbook.Scenarios[0].Events[0]),
		core.WithScenarioLogger(recorder),
	)).NoError(t)

	var alertData any
	ctx := model.NewContext()
	gt.R1(c.HandleAlert(ctx, "my_test", alertData)).NoError(t)
	gt.N(t, calledMock).Equal(0)

	gt.V(t, recorder.Log.ID).Equal("s1")
	gt.V(t, recorder.Log.Title).Equal("Scenario 1")
	gt.A(t, recorder.Log.Results).Length(1).At(0, func(t testing.TB, v *model.PlayLog) {
		gt.V(t, v.Alert.Title).Equal("test alert")
		gt.A(t, v.Alert.Attrs).Length(1).At(0, func(t testing.TB, v model.Attribute) {
			gt.V(t, v.Key).Equal("c")

			// Value has been converted to float64 by encoding/decoding json
			gt.V(t, v.Value).Equal(float64(1))
		})
	})
}

func TestGlobalAttr(t *testing.T) {
	var alertData any

	alertPolicy := gt.R1(policy.New(
		policy.WithPackage("alert"),
		policy.WithFile("testdata/global_attr/alert.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	actionPolicy := gt.R1(policy.New(
		policy.WithPackage("action"),
		policy.WithFile("testdata/global_attr/action.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	var calledMock int
	mock := func(ctx *model.Context, alert model.Alert, _ model.ActionArgs) (any, error) {
		calledMock++
		return nil, nil
	}

	c := gt.R1(chain.New(
		core.WithPolicyAlert(alertPolicy),
		core.WithPolicyAction(actionPolicy),
		core.WithExtraAction("mock", mock),
	)).NoError(t)

	ctx := model.NewContext()

	// call HandleAlert twice, but mock action should be called only once
	gt.R1(c.HandleAlert(ctx, "my_alert", alertData)).NoError(t)
	gt.R1(c.HandleAlert(ctx, "my_alert", alertData)).NoError(t)
	gt.N(t, calledMock).Equal(1)
}

func TestGlobalAttrRaceCondition(t *testing.T) {
	var alertData any

	alertPolicy := gt.R1(policy.New(
		policy.WithPackage("alert"),
		policy.WithFile("testdata/countup/alert.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	actionPolicy := gt.R1(policy.New(
		policy.WithPackage("action"),
		policy.WithFile("testdata/countup/action.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	cache := memory.New()

	var calledMock int
	mock := func(ctx *model.Context, alert model.Alert, _ model.ActionArgs) (any, error) {
		calledMock++
		return nil, nil
	}

	threadNum := 64

	c := gt.R1(chain.New(
		core.WithPolicyAlert(alertPolicy),
		core.WithPolicyAction(actionPolicy),
		core.WithExtraAction("mock", mock),
		core.WithDatabase(cache),
	)).NoError(t)

	ctx := model.NewContext()
	wg := sync.WaitGroup{}
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gt.R1(c.HandleAlert(ctx, "my_alert", alertData)).NoError(t)
		}()
	}
	wg.Wait()

	gt.N(t, calledMock).Equal(threadNum)
	attrs := gt.R1(cache.GetAttrs(ctx, "default")).NoError(t)
	gt.A(t, attrs).Length(1).At(0, func(t testing.TB, v model.Attribute) {
		gt.V(t, v.Key).Equal("counter")
		gt.V(t, v.Value).Equal(float64(threadNum))
	})
}
