package github

import (
	"bytes"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"text/template"

	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/github"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/gots/ptr"
	"golang.org/x/exp/slog"
)

//go:embed issuer_template.md
var issueTemplateData string

var issueTemplate *template.Template

func init() {
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}

	issueTemplate = template.Must(template.New("issue").Funcs(funcMap).Parse(issueTemplateData))
}

func CreateIssue(ctx *model.Context, alert model.Alert, args model.ActionArgs) (any, error) {
	// Create a new issue body from template
	var buf bytes.Buffer
	if err := issueTemplate.Execute(&buf, alert); err != nil {
		return nil, goerr.Wrap(err, "Failed to render issue template")
	}
	req := &github.IssueRequest{
		Title: &alert.Title,
		Body:  github.String(buf.String()),
	}

	// Required arguments
	appID, ok := args["app_id"].(int)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidConfig, "app_id is required")
	}

	installID, ok := args["install_id"].(int)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidConfig, "install_id is required")
	}

	privateKey, ok := args["private_key"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidConfig, "private_key is required")
	} else if !isRSAPrivateKey(privateKey) {
		return nil, goerr.Wrap(types.ErrActionInvalidConfig, "private_key must be RSA private key")
	}

	owner, ok := args["owner"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidConfig, "owner is required")
	}

	repo, ok := args["repo"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidConfig, "repo is required")
	}

	// Optional arguments
	if v, ok := args["assignee"].(string); ok && v != "" {
		req.Assignee = github.String(v)
	}
	if v, ok := args["labels"].([]string); ok && len(v) > 0 {
		req.Labels = &v
	}

	if ctx.DryRun() {
		return nil, nil
	}

	rt := http.DefaultTransport

	itr, err := ghinstallation.New(rt, int64(appID), int64(installID), []byte(privateKey))
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to create GitHub App installation transport")
	}

	client := github.NewClient(&http.Client{Transport: itr})

	issue, resp, err := client.Issues.Create(ctx, owner, repo, req)
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, goerr.Wrap(err).With("resp", resp)
	}

	utils.Logger().Info("Created GitHub issue",
		slog.Any("issue_number", ptr.From(issue.Number)),
		slog.Any("title", ptr.From(issue.Title)),
	)

	return issue, nil
}

func isRSAPrivateKey(s string) bool {
	block, _ := pem.Decode([]byte(s))
	if block == nil {
		return false
	}

	_, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	return err == nil
}
