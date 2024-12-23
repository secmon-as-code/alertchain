package gemini_test

import (
	"context"
	"os"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/alertchain/pkg/infra/gemini"
	"github.com/secmon-lab/alertchain/pkg/utils"
)

func TestGenerateRule(t *testing.T) {
	var (
		projectID string
		location  string
	)
	if err := utils.LoadEnv(
		utils.EnvDef("TEST_GEMINI_PROJECT_ID", &projectID),
		utils.EnvDef("TEST_GEMINI_LOCATION", &location),
	); err != nil {
		t.Skipf("Skip test due to missing env: %v", err)
	}

	ctx := context.Background()
	client, err := gemini.New(ctx, projectID, location)
	gt.NoError(t, err)

	policy := gt.R1(os.ReadFile("scc.rego")).NoError(t)
	alert := gt.R1(os.ReadFile("alert.json")).NoError(t)
	prompt := `
Instructions:

The initial JSON data provided contains information about false positive alerts. Based on the code given thereafter, generate a new Rego policy file to ignore these alerts.


Constraints:


The new Rego policy file must include the content of all existing rules.
Integrate rules if possible.
The output should be in Rego code format only, not Markdown.
Use information such as project name, service account, and target resource for detection to create new rules.
Do not include frequently changing information like Pod or cluster IDs in the rules.
	`

	t.Log("Generate new policy")
	resp, err := client.Generate(ctx, prompt, string(alert), string(policy))
	gt.NoError(t, err)
	for _, line := range resp {
		t.Log(line)
	}
}
