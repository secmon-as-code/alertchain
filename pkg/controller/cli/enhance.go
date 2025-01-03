package cli

import (
	"context"
	_ "embed"

	"github.com/secmon-lab/alertchain/pkg/controller/cli/config"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/infra/gemini"
	"github.com/secmon-lab/alertchain/pkg/usecase"
	"github.com/urfave/cli/v3"
)

func cmdEnhance() *cli.Command {
	return &cli.Command{
		Name:    "enhance",
		Aliases: []string{"e"},
		Usage:   "Enhance alertchain policy with Generative AI",
		Commands: []*cli.Command{
			cmdEnhanceIgnore(),
		},
	}
}

func cmdEnhanceIgnore() *cli.Command {
	var (
		input           usecase.EnhanceIgnorePolicyInput
		alertIDSet      []string
		geminiProjectID string
		geminiLocation  string

		dbCfg config.Database
	)

	flags := []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "alert-id",
			Aliases:     []string{"i"},
			Usage:       "Alert ID to ignore",
			Sources:     cli.EnvVars("ALERTCHAIN_ALERT_ID"),
			Required:    true,
			Destination: (*[]string)(&alertIDSet),
		},
		&cli.StringFlag{
			Name:        "base-policy-file",
			Aliases:     []string{"b"},
			Usage:       "Base policy file. It will be used as a template",
			Sources:     cli.EnvVars("ALERTCHAIN_BASE_POLICY"),
			Required:    true,
			Destination: &input.BasePolicyFile,
		},
		&cli.StringFlag{
			Name:        "test-data-dir",
			Aliases:     []string{"d"},
			Usage:       "Directory path to store test data (e.g. alert/testdata/your_rule)",
			Sources:     cli.EnvVars("ALERTCHAIN_TEST_DATA_DIR"),
			Required:    true,
			Destination: &input.TestDataDir,
		},
		&cli.StringFlag{
			Name:        "test-data-rego-path",
			Aliases:     []string{"r"},
			Usage:       "Path to store test data in rego format (e.g. data.alert.testdata.your_rule)",
			Sources:     cli.EnvVars("ALERTCHAIN_TEST_DATA_REGO_PATH"),
			Required:    true,
			Destination: &input.TestDataRegoPath,
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
			Usage:       "Google Cloud Location for Gemini (e.g. us-central1)",
			Sources:     cli.EnvVars("ALERTCHAIN_GEMINI_LOCATION"),
			Required:    true,
			Destination: &geminiLocation,
		},
		&cli.BoolFlag{
			Name:        "overwrite",
			Aliases:     []string{"w"},
			Usage:       "Overwrite existing base policy file",
			Sources:     cli.EnvVars("ALERTCHAIN_OVERWRITE"),
			Destination: &input.OverWrite,
		},
	}

	flags = append(flags, dbCfg.Flags()...)

	return &cli.Command{
		Name:  "ignore",
		Usage: "Create new ignore policy based on the alert with Gemini",
		Flags: flags,

		Action: func(ctx context.Context, cmd *cli.Command) error {
			for _, id := range alertIDSet {
				input.AlertIDs = append(input.AlertIDs, types.AlertID(id))
			}

			geminiClient, err := gemini.New(ctx, geminiProjectID, geminiLocation)
			if err != nil {
				return err
			}

			dbClient, dbClose, err := dbCfg.New(ctx)
			if err != nil {
				return err
			}
			defer dbClose()

			if err := usecase.EnhanceIgnorePolicy(ctx, dbClient, geminiClient, input); err != nil {
				return err
			}

			return nil
		},
	}
}
