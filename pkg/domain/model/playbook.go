package model

import (
	"encoding/json"
	"path/filepath"

	"github.com/google/go-jsonnet"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

type Playbook struct {
	Scenarios []*Scenario   `json:"scenarios"`
	Env       types.EnvVars `json:"env"`
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

type Event struct {
	Input   any                        `json:"input"`
	Schema  types.Schema               `json:"schema"`
	Results map[types.ActionName][]any `json:"results"`

	actionIndex map[types.ActionName]int
}

type Scenario struct {
	ID     types.ScenarioID    `json:"id"`
	Title  types.ScenarioTitle `json:"title"`
	Events []Event             `json:"events"`
}

func (x *Scenario) Validate() error {
	if x.ID == "" {
		return goerr.Wrap(types.ErrInvalidScenario, "scenario ID is empty")
	}
	if x.Title == "" {
		return goerr.Wrap(types.ErrInvalidScenario, "scenario title is empty")
	}

	// check event schema
	for _, e := range x.Events {
		if err := e.Validate(); err != nil {
			return goerr.Wrap(err, "invalid event")
		}
	}

	return nil
}

func (x *Event) Validate() error {
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

func (x *Event) GetResult(actionName types.ActionName) any {
	if x.actionIndex == nil {
		x.actionIndex = map[types.ActionName]int{}
	}

	idx, ok := x.actionIndex[actionName]
	if !ok {
		idx = 0
	}
	if len(x.Results[actionName]) <= idx {
		return nil
	}
	x.actionIndex[actionName] = idx + 1

	return x.Results[actionName][idx]
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

	return book.Validate()
}
