package github

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

type CreateIssue struct {
	appID      int64
	installID  int64
	privateKey []byte
	owner      string
	repo       string

	rt http.RoundTripper
}

const CreateIssueID = "github-create-issue"

func NewCreateIssue(config model.ActionConfig) (model.Action, error) {
	newErr := func(msg string) error {
		return goerr.Wrap(types.ErrInvalidActionConfig, fmt.Sprintf("%s for %s", msg, CreateIssueID))
	}

	var ok bool
	action := &CreateIssue{}

	if v, ok := config["app_id"].(int); !ok {
		return nil, newErr("app_id is not set or invalid type")
	} else {
		action.appID = int64(v)
	}

	if v, ok := config["install_id"].(int); !ok {
		return nil, newErr("install_id is not set or invalid type")
	} else {
		action.installID = int64(v)
	}

	if envVar, ok := config["private_key_env"].(string); !ok {
		return nil, newErr("private_key_env is not set or invalid type")
	} else if key, ok := os.LookupEnv(envVar); !ok {
		return nil, newErr("env var '" + envVar + "' is not set")
	} else {
		action.privateKey = []byte(key)
	}

	if action.owner, ok = config["owner"].(string); !ok {
		return nil, newErr("owner is not set or invalid type")
	}
	if action.repo, ok = config["repo"].(string); !ok {
		return nil, newErr("repo is not set or invalid type")
	}

	return action, nil
}

func (x *CreateIssue) Run(ctx *types.Context, alert *model.Alert, args ...*model.Attribute) (*model.ChangeRequest, error) {
	ctx.Log().Trace("starting CreateGitHubIssue")

	rt := http.DefaultTransport
	if x.rt != nil {
		rt = x.rt
	}

	itr, err := ghinstallation.New(rt, x.appID, x.installID, x.privateKey)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	client := github.NewClient(&http.Client{Transport: itr})

	issue, resp, err := client.Issues.Create(ctx, x.owner, x.repo, alert2issue(alert))
	if err != nil {
		return nil, goerr.Wrap(err)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, goerr.Wrap(err).With("resp", resp)
	}

	var req model.ChangeRequest
	req.AddReference(&model.Reference{
		Source: "CreateGitHubIssue",
		Title:  "alert issue",
		URL:    *issue.HTMLURL,
	})

	ctx.Log().Trace("exiting CreateGitHubIssue")
	return &req, nil
}

type issueBody struct {
	lines []string
}

func (x *issueBody) add(s ...string) {
	x.lines = append(x.lines, s...)
}
func (x *issueBody) fmt(format string, args ...interface{}) {
	x.lines = append(x.lines, fmt.Sprintf(format, args...))
}
func (x *issueBody) str() *string {
	d := strings.Join(x.lines, "\n")
	return &d
}

func alert2issue(alert *model.Alert) *github.IssueRequest {
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
		"- Created at: "+alert.CreatedAt.Format("2006-01-02 15:04:05"),
		"- Detected by: "+alert.Detector,
		"- Severity: "+string(alert.Severity),
		"- ID: "+string(alert.ID()),
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

		sort.Slice(alert.Attributes, func(i, j int) bool {
			return alert.Attributes[i].Key < alert.Attributes[j].Key
		})
		for _, attr := range alert.Attributes {
			b.fmt("| `%s` | %s | %s | %s |", attr.Key, attr.Value, attr.Type, attr.Contexts.String())
		}
	} else {
		b.add("n/a")
	}

	b.add("", "## References")
	if len(alert.References) > 0 {
		for _, ref := range alert.References {
			b.fmt("- %s: [%s](%s) (%s)", ref.Source, ref.Title, ref.URL, ref.Comment)
		}
	} else {
		b.add("n/a")
	}

	return &github.IssueRequest{
		Title: &alert.Title,
		Body:  b.str(),
	}
}
