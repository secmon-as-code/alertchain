package jira

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/andygrunwald/go-jira"
	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/utils"
)

//go:embed issue_template.txt
var issueTemplateData string

var issueTemplate *template.Template

func init() {
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}

	issueTemplate = template.Must(template.New("issue").Funcs(funcMap).Parse(issueTemplateData))
}

func execTemplate(data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := issueTemplate.Execute(&buf, data); err != nil {
		return "", goerr.Wrap(err, "Failed to execute issue template")
	}
	return buf.String(), nil
}

func CreateIssue(ctx context.Context, alert model.Alert, args model.ActionArgs) (any, error) {
	var (
		accountID string
		userName  string
		token     string
		baseURL   string
		project   string
		issueType string
		labels    []string
		assignee  string
	)

	if err := args.Parse(
		model.ArgDef("account_id", &accountID),
		model.ArgDef("user", &userName),
		model.ArgDef("secret_token", &token),
		model.ArgDef("base_url", &baseURL),
		model.ArgDef("project", &project),
		model.ArgDef("issue_type", &issueType),
		model.ArgDef("labels", &labels, model.ArgOptional()),
		model.ArgDef("assignee", &assignee, model.ArgOptional()),
	); err != nil {
		return nil, err
	}

	tp := jira.BasicAuthTransport{
		Username: userName,
		Password: token,
	}

	jiraClient, err := jira.NewClient(tp.Client(), baseURL)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to create JIRA client")
	}

	buf, err := execTemplate(alert)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to execute issue template")
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Reporter: &jira.User{
				AccountID: accountID,
			},
			Type: jira.IssueType{
				Name: issueType,
			},
			Project: jira.Project{
				Key: project,
			},
			Summary:     alert.Title,
			Description: buf,
		},
	}
	if len(labels) > 0 {
		i.Fields.Labels = labels
	}
	if assignee != "" {
		i.Fields.Assignee = &jira.User{
			AccountID: assignee,
		}
	}

	issue, resp, err := jiraClient.Issue.CreateWithContext(ctx, &i)
	if err != nil {
		data, _ := io.ReadAll(resp.Body)
		return nil, goerr.Wrap(err, "Failed to create issue", goerr.V("body", string(data)))
	}

	fname := fmt.Sprintf("alert-%s.json", alert.ID)
	body := strings.NewReader(alert.Raw)
	if _, _, err := jiraClient.Issue.PostAttachmentWithContext(ctx, issue.ID, body, fname); err != nil {
		return nil, goerr.Wrap(err, "Failed to post attachment")
	}

	return utils.ToAny(issue)
}
