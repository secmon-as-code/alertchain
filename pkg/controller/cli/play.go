package cli

import (
	"io"
	"os"
	"path/filepath"

	"log/slog"

	"github.com/m-mizutani/alertchain/pkg/chain/core"
	"github.com/m-mizutani/alertchain/pkg/controller/cli/config"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/logging"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
)

func cmdPlay() *cli.Command {
	var (
		playbookPath string
		outDir       string
		scenarioIDs  cli.StringSlice

		policyCfg config.Policy
	)

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "playbook",
			Aliases:     []string{"b"},
			Usage:       "playbook file",
			EnvVars:     []string{"ALERTCHAIN_PLAYBOOK"},
			Required:    true,
			Destination: &playbookPath,
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       "output directory",
			EnvVars:     []string{"ALERTCHAIN_OUTPUT"},
			Destination: &outDir,
			Value:       "./output",
		},
		&cli.StringSliceFlag{
			Name:        "scenario",
			Aliases:     []string{"s"},
			Usage:       "scenario ID to play. If not specified, all scenarios are played",
			EnvVars:     []string{"ALERTCHAIN_SCENARIO"},
			Destination: &scenarioIDs,
		},
	}
	flags = append(flags, policyCfg.Flags()...)

	return &cli.Command{
		Name:    "play",
		Aliases: []string{"p"},
		Usage:   "Simulate alertchain policy",
		Flags:   flags,

		Action: func(c *cli.Context) error {
			// Load playbook
			var playbook model.Playbook
			if err := model.ParsePlaybook(playbookPath, os.ReadFile, &playbook); err != nil {
				return goerr.Wrap(err, "failed to parse playbook")
			}

			ctx := model.NewContext(
				model.WithBase(c.Context),
				model.WithCLI(),
			)
			ctx.Logger().Info("starting alertchain with play mode", slog.Any("playbook", playbookPath))

			targets := make(map[types.ScenarioID]struct{})
			for _, id := range scenarioIDs.Value() {
				targets[types.ScenarioID(id)] = struct{}{}
			}

			for _, s := range playbook.Scenarios {
				if _, ok := targets[s.ID]; len(targets) > 0 && !ok {
					continue
				}

				if err := playScenario(ctx, s, &policyCfg, outDir, playbook.Env); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

type actionMockWrapper struct {
	ev *model.Event
}

func (x *actionMockWrapper) GetResult(name types.ActionName) any {
	return x.ev.GetResult(name)
}

func playScenario(ctx *model.Context, scenario *model.Scenario, cfg *config.Policy, outDir string, envVars types.EnvVars) error {
	ctx.Logger().Debug("Start scenario", slog.Any("scenario", scenario))

	w, err := openLogFile(outDir, string(scenario.ID))
	if err != nil {
		return err
	}
	defer func() {
		if err := w.Close(); err != nil {
			ctx.Logger().Warn("Failed to close log file", slog.String("err", err.Error()))
		}
	}()
	lg := logging.NewJSONLogger(w, scenario)

	mockWrapper := &actionMockWrapper{}
	options := []core.Option{
		core.WithScenarioLogger(lg),
		core.WithActionMock(mockWrapper),
	}

	if envVars != nil {
		options = append(options, core.WithEnv(func() types.EnvVars {
			return envVars
		}))
	}

	chain, err := buildChain(cfg, options...)
	if err != nil {
		return err
	}

	for i, ev := range scenario.Events {
		mockWrapper.ev = &scenario.Events[i]
		if _, err := chain.HandleAlert(ctx, ev.Schema, ev.Input); err != nil {
			lg.LogError(err)
			break
		}
	}

	if err := lg.Flush(); err != nil {
		ctx.Logger().Error("Failed to close scenario logger", "err", err)
	}

	return nil
}

func openLogFile(dir, name string) (io.WriteCloser, error) {
	dirName := filepath.Clean(filepath.Join(dir, name))
	// #nosec G301
	if err := os.MkdirAll(dirName, 0755); err != nil {
		return nil, goerr.Wrap(err, "Failed to create scenario logging directory")
	}

	path := filepath.Join(dirName, "data.json")
	fd, err := os.Create(filepath.Clean(path))
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to create scenario logging file")
	}

	return fd, nil
}
