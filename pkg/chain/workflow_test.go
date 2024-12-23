package chain_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/alertchain/pkg/chain"
	"github.com/secmon-lab/alertchain/pkg/chain/core"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/infra/logging"
	"github.com/secmon-lab/alertchain/pkg/infra/memory"
	"github.com/secmon-lab/alertchain/pkg/infra/policy"
	"github.com/secmon-lab/alertchain/pkg/service"
)

type buffer struct {
	bytes.Buffer
}

func (x *buffer) Close() error {
	return nil
}

func TestWorkflow(t *testing.T) {
	actionPolicy := gt.R1(policy.New(
		policy.WithPackage("action"),
		policy.WithFile("testdata/play_workflow/action.rego"),
		policy.WithReadFile(read),
	)).NoError(t)

	var playbook model.Playbook
	gt.NoError(t, model.ParsePlaybook("testdata/play_workflow/playbook.jsonnet", read, &playbook))
	gt.A(t, playbook.Scenarios).Length(1)

	var calledMock int
	mock := func(ctx context.Context, alert model.Alert, _ model.ActionArgs) (any, error) {
		calledMock++
		return nil, nil
	}

	buf := &buffer{}
	recorder := logging.NewJSONLogger(buf, playbook.Scenarios[0])
	c := core.New(
		core.WithPolicyAction(actionPolicy),
		core.WithExtraAction("mock", mock),
		core.WithActionMock(&playbook.Scenarios[0].Events[0]),
		core.WithScenarioLogger(recorder),
		core.WithEnv(func() types.EnvVars { return types.EnvVars{} }),
		core.WithEnablePrint(),
	)

	ctx := context.Background()
	alert := model.NewAlert(model.AlertMetaData{
		Title: "test-alert",
	}, "test-alert", "test-data")

	w := gt.R1(service.New(memory.New()).Workflow.Create(ctx, alert)).NoError(t)
	wf := chain.NewWorkflow(c, alert, w)
	gt.NoError(t, wf.Run(ctx))
	recorder.Flush()

	var log model.ScenarioLog
	gt.NoError(t, json.NewDecoder(bytes.NewReader(buf.Bytes())).Decode(&log))

	gt.V(t, log.ID).Equal("s1")
	gt.V(t, log.Title).Equal("Scenario 1")
	gt.A(t, log.Results).Length(1).At(0, func(t testing.TB, v *model.PlayLog) {
		gt.V(t, v.Alert.Title).Equal("test-alert")
		gt.A(t, v.Actions).Length(2).At(0, func(t testing.TB, v *model.ActionLog) {
			gt.V(t, v.Seq).Equal(0)
			gt.V(t, v.Uses).Equal("mock")
			gt.V(t, v.ID).Equal("1st")
			gt.M(t, v.Args).EqualAt("tick", float64(1))
		}).At(1, func(t testing.TB, v *model.ActionLog) {
			gt.V(t, v.Seq).Equal(1)
			gt.V(t, v.Uses).Equal("mock")
			gt.V(t, v.ID).Equal("2nd")
			gt.M(t, v.Args).EqualAt("tick", float64(2))
		})
	})
}
