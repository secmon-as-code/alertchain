package model_test

import (
	"embed"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/gt"
)

//go:embed testdata/playbook/*.jsonnet
//go:embed testdata/playbook/*.json
var playbooks embed.FS

func TestParsePlaybook(t *testing.T) {
	var playbook model.Playbook
	gt.NoError(t, model.ParsePlaybook("testdata/playbook/playbook1.jsonnet", playbooks.ReadFile, &playbook))
	gt.Array(t, playbook.Scenarios).Length(1).At(0, func(t testing.TB, v *model.Scenario) {
		gt.V(t, v.Name).Equal("test1")

		alert := gt.Cast[map[string]any](t, v.Alert)
		gt.M(t, alert).EqualAt("color", "blue")

		gt.V(t, v.Schema).Equal("scc")
		gt.M(t, v.Results).At("ticket", func(t testing.TB, v []any) {
			gt.Array(t, v).Length(1).At(0, func(t testing.TB, v any) {
				r := gt.Cast[map[string]any](t, v)
				gt.Map(t, r).EqualAt("name", "test1")
			})
		})
	})
}
