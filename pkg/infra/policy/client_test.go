package policy_test

import (
	"context"
	"os"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/alertchain/pkg/infra/policy"
)

const examplePolicy = `
package test

default allow = false

allow {
	input.role == "admin"
}
`

type examplePolicyResult struct {
	Allow bool `json:"allow"`
}

func TestClient_Query(t *testing.T) {
	client, err := policy.New(policy.WithPolicyData("test.rego", examplePolicy), policy.WithPackage("test"))
	gt.NoError(t, err)

	ctx := context.Background()

	tests := []struct {
		name   string
		input  interface{}
		output *examplePolicyResult
		expect bool
	}{
		{
			name:   "admin role should be allowed",
			input:  map[string]interface{}{"role": "admin"},
			output: new(examplePolicyResult),
			expect: true,
		},
		{
			name:   "non-admin role should not be allowed",
			input:  map[string]interface{}{"role": "user"},
			output: new(examplePolicyResult),
			expect: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := client.Query(ctx, test.input, test.output)
			gt.NoError(t, err)
			gt.V(t, test.output.Allow).Equal(test.expect)
		})
	}
}

func TestClient_New_WithFile(t *testing.T) {
	policyFile := "test.rego"
	err := os.WriteFile(policyFile, []byte(examplePolicy), 0644)
	gt.NoError(t, err)
	defer os.Remove(policyFile)

	client, err := policy.New(policy.WithFile(policyFile), policy.WithPackage("test"))
	gt.NoError(t, err)

	ctx := context.Background()

	input := map[string]interface{}{"role": "admin"}
	var output examplePolicyResult
	err = client.Query(ctx, input, &output)
	gt.NoError(t, err)
	gt.B(t, output.Allow).True()
}

func TestClient_New_WithDir(t *testing.T) {
	policyDir := "policy_test"
	err := os.Mkdir(policyDir, 0755)
	gt.NoError(t, err)
	defer os.RemoveAll(policyDir)

	policyFile := policyDir + "/test.rego"
	err = os.WriteFile(policyFile, []byte(examplePolicy), 0644)
	gt.NoError(t, err)

	client, err := policy.New(policy.WithDir(policyDir), policy.WithPackage("test"))
	gt.NoError(t, err)

	ctx := context.Background()

	input := map[string]interface{}{"role": "admin"}
	var output examplePolicyResult
	err = client.Query(ctx, input, &output)
	gt.NoError(t, err)
	gt.B(t, output.Allow).True()
}

func TestClient_New_NoPolicy(t *testing.T) {
	_, err := policy.New()
	gt.Error(t, err)
}

func TestClient_Query_NoResult(t *testing.T) {
	client, err := policy.New(policy.WithPolicyData("test.rego", examplePolicy), policy.WithPackage("test"))
	gt.NoError(t, err)

	ctx := context.Background()

	input := map[string]interface{}{"unknown_key": "unknown_value"}
	var output bool
	err = client.Query(ctx, input, &output)
	gt.Error(t, err)
}
