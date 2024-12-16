package model_test

import (
	"embed"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
)

//go:embed testdata/playbook/*.jsonnet
//go:embed testdata/playbook/*.json
var playbooks embed.FS

func TestParseScenario(t *testing.T) {
	s, err := model.ParseScenario("testdata/playbook/base.jsonnet", playbooks.ReadFile)
	gt.NoError(t, err)
	gt.Equal(t, s.ID, "test1")
	gt.V(t, s.Events[0].Actions["chatgpt.query"][0]).NotNil()
}
