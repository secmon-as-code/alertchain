package github

import (
	"bytes"
	_ "embed"
	"text/template"

	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/goerr"
	"golang.org/x/exp/slog"
)

type Issuer struct {
	id     types.ActionID
	client *github.Client

	owner string
	repo  string
}

//go:embed issuer_template.md
var issueTemplateData string

var issueTemplate *template.Template

func init() {
	issueTemplate = template.Must(template.New("issue").Parse(issueTemplateData))
}

type IssuerFactory struct{}

func (x *IssuerFactory) Name() types.ActionName {
	return "github-issuer"
}

func (x *IssuerFactory) New(id types.ActionID, cfg model.ActionConfigValues) (interfaces.Action, error) {
	appID, ok := cfg["app_id"].(int)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidConfig, "app_id is required")
	}

	installID, ok := cfg["install_id"].(int)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidConfig, "install_id is required")
	}

	privateKey, ok := cfg["private_key"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidConfig, "private_key is required")
	}

	owner, ok := cfg["owner"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidConfig, "owner is required")
	}

	repo, ok := cfg["repo"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidConfig, "repo is required")
	}

	rt := http.DefaultTransport

	itr, err := ghinstallation.New(rt, int64(appID), int64(installID), []byte(privateKey))
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to create GitHub App installation transport")
	}

	client := github.NewClient(&http.Client{Transport: itr})

	return &Issuer{
		id:     id,
		client: client,
		owner:  owner,
		repo:   repo,
	}, nil
}

func (x *Issuer) ID() types.ActionID { return x.id }

func (x *Issuer) Run(ctx *types.Context, alert model.Alert, params model.ActionParams) error {
	var buf bytes.Buffer
	if err := issueTemplate.Execute(&buf, alert); err != nil {
		return goerr.Wrap(err, "Failed to render issue template")
	}
	req := &github.IssueRequest{
		Title: &alert.Title,
		Body:  github.String(buf.String()),
	}

	issue, resp, err := x.client.Issues.Create(ctx, x.owner, x.repo, req)
	if err != nil {
		return goerr.Wrap(err)
	}
	if resp.StatusCode != http.StatusCreated {
		return goerr.Wrap(err).With("resp", resp)
	}

	utils.Logger().Info("Created issue", slog.Any("issue", issue))

	return nil
}
