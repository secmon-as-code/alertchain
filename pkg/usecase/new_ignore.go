package usecase

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/interfaces"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

//go:embed prompt/new_ignore.md
var newIgnorePrompt string

//go:embed prompt/alert_slug.md
var alertSlugPrompt string

func genTestFilePath(filePath string) string {
	ext := filepath.Ext(filePath)
	base := strings.TrimSuffix(filePath, ext)
	newPath := base + "_test" + ext
	return newPath
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func getRegoPackageName(policyData string) (string, error) {
	module, err := ast.ParseModule("temp.rego", policyData)
	if err != nil {
		return "", goerr.Wrap(err, "failed to parse rego file")
	}

	if len(module.Package.Path) <= 1 {
		return "", goerr.New("package path is empty")
	}

	var terms []string
	for _, term := range module.Package.Path[1:] {
		terms = append(terms, strings.Trim(term.String(), `"`))
	}
	return strings.Join(terms, "."), nil
}

type NewIgnorePolicyInput struct {
	AlertID          types.AlertID
	BasePolicyFile   string
	TestDataDir      string
	TestDataRegoPath string
	OverWrite        bool
}

func (x NewIgnorePolicyInput) Validate() error {
	if x.AlertID == "" {
		return goerr.New("AlertID is empty")
	}
	if x.BasePolicyFile == "" {
		return goerr.New("BasePolicyFile is empty")
	}
	if x.TestDataDir == "" {
		return goerr.New("TestDataDir is empty")
	}
	if x.TestDataRegoPath == "" {
		return goerr.New("TestDataRegoPath is empty")
	}
	return nil
}

func NewIgnorePolicy(ctx context.Context,
	dbClient interfaces.Database,
	genAI interfaces.GenAI,
	input NewIgnorePolicyInput,
) error {
	if err := input.Validate(); err != nil {
		return err
	}

	logger := ctxutil.Logger(ctx)

	alert, err := dbClient.GetAlert(ctx, input.AlertID)
	if err != nil {
		return err
	}
	alertData, err := json.Marshal(alert.Data)
	if err != nil {
		return goerr.Wrap(err, "failed to marshal alert data")
	}
	logger.Info("Got alert data", "alertID", input.AlertID)

	basePolicy, err := os.ReadFile(input.BasePolicyFile)
	if err != nil {
		return goerr.Wrap(err, "failed to read base policy file")
	}
	logger.Info("Got base policy", "file", input.BasePolicyFile)

	pkgName, err := getRegoPackageName(string(basePolicy))
	if err != nil {
		return err
	}

	testFiles, err := os.ReadDir(input.TestDataDir)
	if err != nil {
		return goerr.Wrap(err, "failed to read test data directory")
	}
	if len(testFiles) > 0 {
		var testFileList []string
		for _, f := range testFiles {
			testFileList = append(testFileList, f.Name())
		}
		alertSlugPrompt += "\nExclude these list: \n" + strings.Join(testFileList, "\n") + "\n"
	}

	slugResp, err := genAI.Generate(ctx, alertSlugPrompt, string(alertData))
	alertSlug := strings.TrimSpace(slugResp[0])
	logger.Info("Generated slug", "slug", alertSlug)

	policyResp, err := genAI.Generate(ctx, newIgnorePrompt, string(alertData), string(basePolicy))
	if err != nil {
		return err
	}

	if input.OverWrite {
		logger.Info("Overwriting base policy file", "file", input.BasePolicyFile)
		if err := os.WriteFile(input.BasePolicyFile, []byte(policyResp[0]), 0644); err != nil {
			return goerr.Wrap(err, "failed to write base policy file")
		}
	} else {
		println(policyResp[0])
	}

	testDataPath := filepath.Join(input.TestDataDir, alertSlug, "data.json")
	if err := os.MkdirAll(filepath.Dir(testDataPath), 0755); err != nil {
		return goerr.Wrap(err, "failed to create test data directory")
	}
	if err := os.WriteFile(testDataPath, alertData, 0644); err != nil {
		return goerr.Wrap(err, "failed to write test data file")
	}
	logger.Info("Wrote test data", "file", testDataPath)

	testFilePath := genTestFilePath(input.BasePolicyFile)
	newTestFile := false
	if !isExist(testFilePath) {
		newTestFile = true
	}

	tf, err := os.OpenFile(testFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return goerr.Wrap(err, "failed to open test file")
	}
	defer tf.Close()

	header := []string{}
	if newTestFile {
		header = []string{
			"package " + pkgName,
			"",
			"import rego.v1",
			"",
		}
	}

	newTest := append(header, []string{
		"test_" + alertSlug + " if {",
		`  resp := alert with input as ` + input.TestDataRegoPath + "." + alertSlug,
		`  count(resp) == 0`,
		`}`,
		"",
	}...)
	newTestBody := strings.Join(newTest, "\n")

	if input.OverWrite {
		if _, err := tf.WriteString(newTestBody); err != nil {
			return goerr.Wrap(err, "failed to write test file")
		}
		logger.Info("Wrote test file", "file", testFilePath)
	} else {
		logger.Info("Skip writing test file", "file", testFilePath, "new_test", newTestBody)
	}

	return nil
}
