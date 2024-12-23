package cli

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/alertchain/pkg/controller/cli/config"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/infra/gemini"
	"github.com/urfave/cli/v3"
)

func cmdNew() *cli.Command {
	return &cli.Command{
		Name:    "new",
		Aliases: []string{"n"},
		Usage:   "Create new alertchain policy",
		Commands: []*cli.Command{
			cmdNewIgnore(),
		},
	}
}

//go:embed prompt/new_ignore.md
var newIgnorePrompt string

//go:embed prompt/alert_slug.md
var alertSlugPrompt string

func cmdNewIgnore() *cli.Command {
	var (
		alertID          types.AlertID
		basePolicyFile   string
		testDataDir      string
		testDataRegoPath string
		geminiProjectID  string
		geminiLocation   string
		overWrite        bool

		dbCfg config.Database
	)

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "alert-id",
			Aliases:     []string{"i"},
			Usage:       "Alert ID to ignore",
			Sources:     cli.EnvVars("ALERTCHAIN_ALERT_ID"),
			Required:    true,
			Destination: (*string)(&alertID),
		},
		&cli.StringFlag{
			Name:        "base-policy-file",
			Aliases:     []string{"b"},
			Usage:       "Base policy file. It will be used as a template",
			Sources:     cli.EnvVars("ALERTCHAIN_BASE_POLICY"),
			Required:    true,
			Destination: &basePolicyFile,
		},
		&cli.StringFlag{
			Name:        "test-data-dir",
			Aliases:     []string{"d"},
			Usage:       "Directory path to store test data",
			Sources:     cli.EnvVars("ALERTCHAIN_TEST_DATA_DIR"),
			Required:    true,
			Destination: &testDataDir,
		},
		&cli.StringFlag{
			Name:        "test-data-rego-path",
			Aliases:     []string{"r"},
			Usage:       "Path to store test data in rego format",
			Sources:     cli.EnvVars("ALERTCHAIN_TEST_DATA_REGO_PATH"),
			Required:    true,
			Destination: &testDataRegoPath,
		},
		&cli.StringFlag{
			Name:        "gemini-project-id",
			Usage:       "Google Cloud Project ID for Gemini",
			Sources:     cli.EnvVars("ALERTCHAIN_GEMINI_PROJECT_ID"),
			Required:    true,
			Destination: &geminiProjectID,
		},
		&cli.StringFlag{
			Name:        "gemini-location",
			Usage:       "Google Cloud Location for Gemini",
			Sources:     cli.EnvVars("ALERTCHAIN_GEMINI_LOCATION"),
			Required:    true,
			Destination: &geminiLocation,
		},
		&cli.BoolFlag{
			Name:        "overwrite",
			Aliases:     []string{"w"},
			Usage:       "Overwrite existing base policy file",
			Sources:     cli.EnvVars("ALERTCHAIN_OVERWRITE"),
			Destination: &overWrite,
		},
	}

	flags = append(flags, dbCfg.Flags()...)

	return &cli.Command{
		Name:  "ignore",
		Usage: "Create new ignore policy based on the alert",
		Flags: flags,

		Action: func(ctx context.Context, cmd *cli.Command) error {
			logger := ctxutil.Logger(ctx)

			geminiClient, err := gemini.New(ctx, geminiProjectID, geminiLocation)
			if err != nil {
				return err
			}

			dbClient, dbClose, err := dbCfg.New(ctx)
			if err != nil {
				return err
			}
			defer dbClose()

			alert, err := dbClient.GetAlert(ctx, alertID)
			if err != nil {
				return err
			}
			alertData, err := json.Marshal(alert.Data)
			if err != nil {
				return goerr.Wrap(err, "failed to marshal alert data")
			}
			logger.Info("Got alert data", "alertID", alertID)

			basePolicy, err := os.ReadFile(basePolicyFile)
			if err != nil {
				return goerr.Wrap(err, "failed to read base policy file")
			}
			logger.Info("Got base policy", "file", basePolicyFile)

			testFiles, err := os.ReadDir(testDataDir)
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

			slugResp, err := geminiClient.Generate(ctx, alertSlugPrompt, string(alertData))
			if err != nil {
				return err
			}
			alertSlug := strings.TrimSpace(slugResp[0])
			logger.Info("Generated slug", "slug", alertSlug)

			policyResp, err := geminiClient.Generate(ctx, newIgnorePrompt, string(alertData), string(basePolicy))
			if err != nil {
				return err
			}

			if overWrite {
				if err := os.WriteFile(basePolicyFile, []byte(policyResp[0]), 0644); err != nil {
					return goerr.Wrap(err, "failed to write base policy file")
				}
			} else {
				println(policyResp[0])
			}

			testDataPath := filepath.Join(testDataDir, alertSlug, "data.json")
			if err := os.MkdirAll(filepath.Dir(testDataPath), 0755); err != nil {
				return goerr.Wrap(err, "failed to create test data directory")
			}
			if err := os.WriteFile(testDataPath, alertData, 0644); err != nil {
				return goerr.Wrap(err, "failed to write test data file")
			}
			logger.Info("Wrote test data", "file", testDataPath)

			testFilePath := genTestFilePath(basePolicyFile)
			file, err := os.OpenFile(testFilePath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return goerr.Wrap(err, "failed to open test file")
			}
			defer file.Close()

			newTest := []string{
				"test_" + alertSlug + " {",
				`  resp := alert with input as ` + testDataRegoPath + "." + alertSlug,
				`  count(resp) == 0`,
				`}`,
			}
			if _, err := file.WriteString(strings.Join(newTest, "\n")); err != nil {
				return goerr.Wrap(err, "failed to write test file")
			}
			logger.Info("Wrote test file", "file", testFilePath)

			return nil
		},
	}
}

func genTestFilePath(filePath string) string {
	ext := filepath.Ext(filePath)
	base := strings.TrimSuffix(filePath, ext)
	newPath := base + "_test" + ext
	return newPath
}
