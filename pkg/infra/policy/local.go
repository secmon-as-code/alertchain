package policy

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/topdown/print"
)

type Local struct {
	compiler *ast.Compiler
}

func NewLocal(path string) (*Local, error) {
	policies := make(map[string]string)

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return goerr.Wrap(err, "fail to stat to read .rego").With("path", path)
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".rego" {
			return nil
		}

		raw, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return goerr.Wrap(err, "fail to read .rego file as local policy").With("path", path)
		}

		policies[path] = string(raw)
		utils.Logger.With("path", path).Debug("read .rego file")

		return nil
	})
	if err != nil {
		return nil, goerr.Wrap(err, "fail to walk rego files/directory")
	}

	compiler, err := ast.CompileModules(policies)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to compile local policies")
	}

	return &Local{
		compiler: compiler,
	}, nil
}

type printLogger struct{}

func (x *printLogger) Print(_ print.Context, msg string) error {
	utils.Logger.With("msg", msg).Debug("print in Rego")
	return nil
}

func (x *Local) Eval(ctx context.Context, in interface{}, out interface{}) error {
	utils.Logger.With("in", in).Trace("start Local.Eval")
	rego := rego.New(
		rego.Query("data"),
		rego.PrintHook(&printLogger{}),
		rego.Compiler(x.compiler),
		rego.Input(in),
	)

	rs, err := rego.Eval(ctx)

	if err != nil {
		return goerr.Wrap(err, "fail to eval local policy").With("input", in)
	}
	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		return goerr.Wrap(types.ErrNoEvalResult)
	}

	utils.Logger.With("rs", rs).Trace("got a result of rego.Eval")

	raw, err := json.Marshal(rs[0].Expressions[0].Value)
	if err != nil {
		return goerr.Wrap(err, "fail to marshal a result of rego.Eval").With("rs", rs)
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return goerr.Wrap(err, "fail to unmarshal a result of rego.Eval to out").With("rs", rs)
	}

	utils.Logger.With("rs", rs).Trace("done Local.Eval")

	return nil
}
