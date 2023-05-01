package model

import (
	"encoding/json"
	"path/filepath"

	"github.com/google/go-jsonnet"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

type Playbook struct {
	Scenarios []*Scenario `json:"scenarios"`
}

func (x *Playbook) Validate() error {
	for _, s := range x.Scenarios {
		if err := s.Validate(); err != nil {
			return goerr.Wrap(err, "invalid scenario")
		}
	}

	// check scenario ID uniqueness
	ids := map[types.ScenarioID]struct{}{}
	for _, s := range x.Scenarios {
		if err := s.Validate(); err != nil {
			return goerr.Wrap(err, "invalid scenario")
		}

		if _, ok := ids[s.ID]; ok {
			return goerr.Wrap(types.ErrInvalidScenario, "scenario ID is not unique")
		}
		ids[s.ID] = struct{}{}
	}

	return nil
}

type Scenario struct {
	ID      types.ScenarioID         `json:"id"`
	Title   types.ScenarioTitle      `json:"title"`
	Alert   any                      `json:"alert"`
	Schema  types.Schema             `json:"schema"`
	Results map[types.ActionID][]any `json:"results"`

	actionIndex map[types.ActionID]int
}

func (x *Scenario) Validate() error {
	if x.ID == "" {
		return goerr.Wrap(types.ErrInvalidScenario, "scenario ID is empty")
	}
	if x.Title == "" {
		return goerr.Wrap(types.ErrInvalidScenario, "scenario title is empty")
	}
	if x.Schema == "" {
		return goerr.Wrap(types.ErrInvalidScenario, "schema is empty")
	}

	return nil
}

func (x *Scenario) ToLog() ScenarioLog {
	return ScenarioLog{
		ID:    x.ID,
		Title: x.Title,
	}
}

func (x *Scenario) GetResult(actionID types.ActionID) any {
	if x.actionIndex == nil {
		x.actionIndex = map[types.ActionID]int{}
	}

	idx, ok := x.actionIndex[actionID]
	if !ok {
		idx = 0
	}
	if len(x.Results[actionID]) <= idx {
		return nil
	}
	x.actionIndex[actionID] = idx + 1

	return x.Results[actionID][idx]
}

type embedImporter struct {
	readFile ReadFile
}

type ReadFile func(name string) ([]byte, error)

func (x *embedImporter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	dirName := filepath.Dir(importedFrom)
	path := filepath.Join(dirName, importedPath)
	fileContent, err := x.readFile(path)
	if err != nil {
		// Fail and return custom error message
		return jsonnet.MakeContents(""), "", goerr.Wrap(err, "failed to read file")
	}
	return jsonnet.MakeContents(string(fileContent)), path, nil
}

func ParsePlaybook(entryFile string, readFile ReadFile, book *Playbook) error {
	vm := jsonnet.MakeVM()
	vm.Importer(&embedImporter{readFile: readFile})

	raw, err := vm.EvaluateFile(entryFile)
	if err != nil {
		return goerr.Wrap(err, "evaluating playbook jsonnet")
	}

	if err := json.Unmarshal([]byte(raw), book); err != nil {
		return goerr.Wrap(err, "unmarshal playbook by jsonnet")
	}

	return nil
}
