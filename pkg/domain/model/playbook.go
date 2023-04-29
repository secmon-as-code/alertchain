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

type Scenario struct {
	Name    string                   `json:"name"`
	Alert   any                      `json:"alert"`
	Schema  types.Schema             `json:"schema"`
	Results map[types.ActionID][]any `json:"results"`
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
