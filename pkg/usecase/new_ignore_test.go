package usecase_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/mock"
	"github.com/secmon-lab/alertchain/pkg/usecase"
)

//go:embed testdata/example_policy.rego
var examplePolicy []byte

func TestNewIgnorePolicy(t *testing.T) {
	ctx := context.Background()

	alertID := types.AlertID("test-alert-id")
	alertData := model.Alert{
		Data: map[string]interface{}{
			"key": "value",
		},
	}
	alert := model.Alert{Data: alertData}
	alertDataJSON, _ := json.Marshal(alertData)
	slug := "testing_slug"

	dbClient := &mock.DatabaseMock{
		GetAlertFunc: func(ctx context.Context, id types.AlertID) (*model.Alert, error) {
			return &alert, nil
		},
	}
	genSeq := 0
	genAI := &mock.GenAIMock{
		GenerateFunc: func(ctx context.Context, prompts ...string) ([]string, error) {
			genSeq++
			switch genSeq {
			case 1:
				return []string{slug}, nil
			case 2:
				return []string{"updated_policy"}, nil
			default:
				t.FailNow()
				return nil, nil
			}
		},
	}

	fd, err := os.CreateTemp("", "*_base_policy.rego")
	gt.NoError(t, err)
	gt.R1(fd.Write(examplePolicy)).NoError(t)
	gt.NoError(t, fd.Close())
	// TODO: Enable this line
	// defer os.Remove(input.BasePolicyFile)

	dir, err := os.MkdirTemp("", "test_data")
	gt.NoError(t, err)
	// TODO: Enable this line
	// defer os.RemoveAll(input.TestDataDir)

	input := usecase.NewIgnorePolicyInput{
		AlertIDs:         []types.AlertID{alertID},
		BasePolicyFile:   fd.Name(),
		TestDataDir:      dir,
		TestDataRegoPath: "test_rego_path",
		OverWrite:        true,
	}

	gt.NoError(t, usecase.NewIgnorePolicy(ctx, dbClient, genAI, input))

	// Check if test data file is created
	testDataPath := filepath.Join(input.TestDataDir, slug, "data.json")
	testDataContent, err := os.ReadFile(testDataPath)
	gt.NoError(t, err)
	gt.S(t, string(testDataContent)).Contains(string(alertDataJSON))

	// Check if test file is updated
	testFilePath := usecase.GenTestFilePath(input.BasePolicyFile)
	testFileContent, err := os.ReadFile(testFilePath)
	gt.NoError(t, err)

	gt.S(t, string(testFileContent)).Contains("package my_alert")
	gt.S(t, string(testFileContent)).Contains("test_testing_slug if {")
}
