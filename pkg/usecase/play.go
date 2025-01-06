package usecase

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/alertchain/pkg/chain"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/infra/recorder"
	"github.com/secmon-lab/alertchain/pkg/utils"
)

type PlayInput struct {
	ScenarioPath string
	OutDir       string
	Targets      []string
	CoreOptions  []chain.Option
}

func (x PlayInput) Validate() error {
	if x.ScenarioPath == "" {
		return goerr.New("scenario path is required")
	}
	if x.OutDir == "" {
		return goerr.New("output directory is required")
	}
	return nil
}

func Play(ctx context.Context, input PlayInput) error {
	if err := input.Validate(); err != nil {
		return goerr.Wrap(err, "invalid input")
	}

	scenarioFiles := make([]string, 0)
	err := filepath.Walk(input.ScenarioPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".jsonnet" {
			scenarioFiles = append(scenarioFiles, path)
		}
		return nil
	})
	if err != nil {
		return goerr.Wrap(err, "failed to walk through playbook directory")
	}

	logger := ctxutil.Logger(ctx)

	var playbook model.Playbook
	for _, scenarioFile := range scenarioFiles {
		logger.Debug("Load scenario", slog.String("file", scenarioFile))
		s, err := model.ParseScenario(scenarioFile, os.ReadFile)
		if err != nil {
			return goerr.Wrap(err, "failed to parse playbook")
		}

		playbook.Scenarios = append(playbook.Scenarios, s)
	}

	if err := playbook.Validate(); err != nil {
		return err
	}

	logger.Info("starting alertchain with play mode",
		"scenario dir", input.ScenarioPath,
		"output dir", input.OutDir,
		"targets", input.Targets,
	)

	targets := make(map[types.ScenarioID]struct{})
	for _, id := range input.Targets {
		targets[types.ScenarioID(id)] = struct{}{}
	}

	result := model.PlayResult{
		StartedAt: time.Now(),
	}

	var errs []error
	for _, s := range playbook.Scenarios {
		if _, ok := targets[s.ID]; len(targets) > 0 && !ok {
			continue
		}

		scenarioAttr := slog.Any("id", s.ID)
		err := playScenario(ctx, s, input.CoreOptions, input.OutDir)
		if err != nil {
			logger.Error("Failed to play scenario", scenarioAttr, "error", err)
		} else {
			logger.Info("Done scenario", scenarioAttr)
		}

		r := &model.ScenarioResult{
			ID:      s.ID,
			Output:  s.GetLogFilePath(input.OutDir),
			Success: err == nil,
		}
		if err != nil {
			r.Error = err.Error()
			errs = append(errs, err)
		}
		result.Scenarios = append(result.Scenarios, r)
	}
	result.FinishedAt = time.Now()

	resultFilePath := filepath.Clean(filepath.Join(input.OutDir, "result.json"))
	fd, err := os.Create(resultFilePath)
	if err != nil {
		return goerr.Wrap(err, "failed to create result file")
	}
	defer utils.SafeClose(ctx, fd)

	encoder := json.NewEncoder(fd)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return goerr.Wrap(err, "failed to write result file")
	}

	if len(errs) > 0 {
		return goerr.New("failed to play scenarios").With("errors", errs)
	}
	return nil
}

type actionMockWrapper struct {
	ev *model.Event
}

func (x *actionMockWrapper) GetResult(name types.ActionName) any {
	return x.ev.GetResult(name)
}

func playScenario(ctx context.Context, scenario *model.Scenario, baseOptions []chain.Option, outDir string) error {
	logger := ctxutil.Logger(ctx)
	logger.Debug("Start scenario", slog.Any("scenario", scenario))

	w, err := openLogFile(scenario.GetLogFilePath(outDir))
	if err != nil {
		return err
	}
	defer func() {
		if err := w.Close(); err != nil {
			logger.Warn("Failed to close log file", slog.String("err", err.Error()))
		}
	}()
	lg := recorder.NewJsonRecorder(w, scenario)

	mockWrapper := &actionMockWrapper{}
	options := append(baseOptions, []chain.Option{
		chain.WithScenarioRecorder(lg),
		chain.WithActionMock(mockWrapper),
	}...)

	if scenario.Env != nil {
		options = append(options, chain.WithEnv(func() types.EnvVars {
			return scenario.Env
		}))
	}

	chain, err := chain.New(options...)
	if err != nil {
		return err
	}

	for i, ev := range scenario.Events {
		mockWrapper.ev = &scenario.Events[i]
		if _, err := chain.HandleAlert(ctx, ev.Schema, ev.Input); err != nil {
			lg.LogError(err)
			attrs := []any{slog.Any("msg", err.Error())}
			if goErr := goerr.Unwrap(err); goErr != nil {
				attrs = append(attrs, slog.Any("value", goErr.Values()))
			}
			logger.Error("Failed to handle alert",
				"scenario.id", scenario.ID,
				"event.sequence", i,
				"event.schema", ev.Schema,
				slog.Group("error", attrs...),
			)
			break
		}
	}

	if err := lg.Flush(); err != nil {
		logger.Error("Failed to close scenario logger", "err", err)
	}

	return nil
}

func openLogFile(fpath string) (io.WriteCloser, error) {
	dirName := filepath.Dir(filepath.Clean(fpath))
	// #nosec G301
	if err := os.MkdirAll(dirName, 0755); err != nil {
		return nil, goerr.Wrap(err, "Failed to create scenario logging directory")
	}

	fd, err := os.Create(filepath.Clean(fpath))
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to create scenario logging file")
	}

	return fd, nil
}
