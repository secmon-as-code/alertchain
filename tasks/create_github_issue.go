package tasks

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/goerr"
)

var logger = utils.Logger

type CreateGitHubIssue struct {
	AppID      int64
	InstallID  int64
	PrivateKey []byte
	Owner      string
	Repo       string

	AlertChainURL string

	// For test
	TestRoundTripper http.RoundTripper
}

func (x *CreateGitHubIssue) Name() string { return "Create GitHub issue" }

func (x *CreateGitHubIssue) Execute(ctx *types.Context, alert *alertchain.Alert) error {
	logger.Trace("starting CreateGitHubIssue")
	rt := http.DefaultTransport
	if x.TestRoundTripper != nil {
		rt = x.TestRoundTripper
	}

	itr, err := ghinstallation.New(rt, x.AppID, x.InstallID, x.PrivateKey)
	if err != nil {
		return goerr.Wrap(err)
	}

	client := github.NewClient(&http.Client{Transport: itr})

	issue, resp, err := client.Issues.Create(ctx, x.Owner, x.Repo, alert2issue(x.AlertChainURL, alert))
	if err != nil {
		return goerr.Wrap(err)
	}
	if resp.StatusCode != http.StatusCreated {
		return goerr.Wrap(err).With("resp", resp)
	}

	alert.AddReference(&ent.Reference{
		Source: "CreateGitHubIssue",
		Title:  "alert issue",
		URL:    *issue.HTMLURL,
	})

	logger.Trace("exiting CreateGitHubIssue")
	return nil
}

type issueBody struct {
	lines []string
}

func (x *issueBody) add(s ...string) {
	x.lines = append(x.lines, s...)
}
func (x *issueBody) str() *string {
	d := strings.Join(x.lines, "\n")
	return &d
}

func alert2issue(url string, alert *alertchain.Alert) *github.IssueRequest {
	var b issueBody

	b.add(
		// TODO: re-enable alertchain link
		// fmt.Sprintf("[alertchain](%s/alert/%s)", url, alert.ID),
		// "",
		"## Description",
		alert.Description+"",
		"",
		"- - - - - - - -",
		"",
		"- Created at: "+time.Unix(alert.CreatedAt, 0).Format("2006-01-02 15:04:05"),
		"- Detected by: "+alert.Detector,
		"- Severity: "+string(alert.Severity),
		"",
		"- - - - - - - -",
		"",
		"## Attributes",
	)

	if len(alert.Attributes) > 0 {
		b.add([]string{
			"| Key | Value | Type | Context |",
			"|:----|:----|:---|:----|",
		}...)
		for _, attr := range alert.Attributes {
			b.add(fmt.Sprintf("| `%s` | %s | %s| %s|", attr.Key, attr.Value, attr.Type, strings.Join(attr.Context, ", ")))
		}
	} else {
		b.add("n/a")
	}

	return &github.IssueRequest{
		Title: &alert.Title,
		Body:  b.str(),
	}
}
