package github

import (
	"bytes"
	"context"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"text/template"

	"net/http"

	"log/slog"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/github"
	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/gots/ptr"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
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

func CreateIssue(ctx context.Context, alert model.Alert, args model.ActionArgs) (any, error) {
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
	appID, ok := args["app_id"].(float64)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "app_id is required")
	}

	installID, ok := args["install_id"].(float64)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "install_id is required")
	}

	privateKey, ok := args["secret_private_key"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "private_key is required")
	} else if !isRSAPrivateKey(privateKey) {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "private_key must be RSA private key")
	}

	owner, ok := args["owner"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "owner is required")
	}

	repo, ok := args["repo"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "repo is required")
	}

	// Optional arguments
	if v, ok := args["assignee"].(string); ok && v != "" {
		req.Assignee = github.String(v)
	}
	if v, ok := args["labels"].([]string); ok && len(v) > 0 {
		req.Labels = &v
	}

	if ctxutil.IsDryRun(ctx) {
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
		return nil, goerr.Wrap(err, "Failed to create GitHub issue")
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, goerr.New("unexpected status code in creating GitHub issue", goerr.V("status", resp.StatusCode))
	}

	ctxutil.Logger(ctx).Debug("Created GitHub issue",
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
