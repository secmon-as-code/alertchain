package policy

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

// Client is a policy engine client
type Client struct {
	dirs     []string
	files    []string
	policies map[string]string

	readFile readFile

	compiler *ast.Compiler
	query    string
}

type RegoPrint func(file string, row int, msg string) error
type readFile func(string) ([]byte, error)

// Option is a functional option for Client
type Option func(x *Client)

// WithDir specifies directory path of .rego policy. Import policy files recursively.
func WithDir(dirPath string) Option {
	return func(x *Client) {
		x.dirs = append(x.dirs, filepath.Clean(dirPath))
	}
}

// WithFile specifies file path of .rego policy. Import policy files recursively.
func WithFile(filePath string) Option {
	return func(x *Client) {
		x.files = append(x.files, filepath.Clean(filePath))
	}
}

// WithReadFile specifies file path of .rego policy. Import policy files recursively.
func WithReadFile(fn func(string) ([]byte, error)) Option {
	return func(x *Client) {
		x.readFile = fn
	}
}

// WithPolicyData specifies raw policy data with name. If the `name` conflicts with file path loaded by WithFile or WithDir, the policy overwrites data loaded by WithFile or WithDir.
func WithPolicyData(name, policy string) Option {
	return func(x *Client) {
		x.policies[name] = policy
	}
}

// WithPackage specifies using package name. e.g. "example.my_policy"
func WithPackage(pkg string) Option {
	return func(x *Client) {
		x.query = "data." + pkg
	}
}

// New creates a new Local client. It requires one or more WithFile, WithDir or WithPolicyData.
func New(options ...Option) (*Client, error) {
	client := &Client{
		query:    "data",
		policies: make(map[string]string),
	}
	for _, opt := range options {
		opt(client)
	}

	policies := make(map[string]string)
	var targetFiles []string
	for _, dirPath := range client.dirs {
		err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return goerr.Wrap(err, "Failed to walk directory").With("path", path)
			}
			if d.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".rego" {
				return nil
			}

			targetFiles = append(targetFiles, path)

			return nil
		})
		if err != nil {
			return nil, goerr.Wrap(err)
		}
	}
	targetFiles = append(targetFiles, client.files...)

	for _, filePath := range targetFiles {
		raw, err := os.ReadFile(filepath.Clean(filePath))
		if err != nil {
			return nil, goerr.Wrap(err, "Failed to read policy file").With("path", filePath)
		}

		policies[filePath] = string(raw)
	}

	for k, v := range client.policies {
		policies[k] = v
	}

	if len(policies) == 0 {
		return nil, goerr.Wrap(types.ErrNoPolicyData)
	}

	compiler, err := ast.CompileModulesWithOpt(policies, ast.CompileOpts{
		EnablePrintStatements: true,
	})
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	client.compiler = compiler

	return client, nil
}

type queryConfig struct {
	pkgSuffix []string
	regoPrint RegoPrint
}

func newQueryConfig(options ...QueryOption) *queryConfig {
	cfg := &queryConfig{}
	for _, opt := range options {
		opt(cfg)
	}
	return cfg
}

type QueryOption func(cfg *queryConfig)

// WithPackageSuffix specifies package suffix. e.g. "example.my_policy"
func WithPackageSuffix(suffix ...string) QueryOption {
	return func(cfg *queryConfig) {
		cfg.pkgSuffix = append(cfg.pkgSuffix, suffix...)
	}
}

func WithRegoPrint(callback RegoPrint) QueryOption {
	return func(cfg *queryConfig) {
		cfg.regoPrint = callback
	}
}

// Query evaluates policy with `input` data. The result will be written to `out`. `out` must be pointer of instance.
func (x *Client) Query(ctx context.Context, input interface{}, output interface{}, options ...QueryOption) error {
	cfg := newQueryConfig(options...)

	query := strings.Join(append([]string{x.query}, cfg.pkgSuffix...), ".")
	regoOpt := []func(r *rego.Rego){
		rego.Query(query),
		rego.Compiler(x.compiler),
		rego.Input(input),
	}
	if cfg.regoPrint != nil {
		regoOpt = append(regoOpt, rego.PrintHook(&regoPrintHook{
			callback: cfg.regoPrint,
		}))
	}

	rs, err := rego.New(regoOpt...).Eval(ctx)
	if err != nil {
		return goerr.Wrap(err, "fail to eval local policy").With("input", input)
	}
	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		return goerr.Wrap(types.ErrNoPolicyResult)
	}

	raw, err := json.Marshal(rs[0].Expressions[0].Value)
	if err != nil {
		return goerr.Wrap(err, "fail to marshal a result of rego.Eval").With("rs", rs)
	}
	if err := json.Unmarshal(raw, output); err != nil {
		return goerr.Wrap(err, "fail to unmarshal a result of rego.Eval to out").With("rs", rs)
	}

	return nil
}
