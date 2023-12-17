package jira

import (
	_ "embed"

	"github.com/andygrunwald/go-jira"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/goerr"
)

func AddComment(ctx *model.Context, alert model.Alert, args model.ActionArgs) (any, error) {
	var (
		accountID string
		userName  string
		token     string
		baseURL   string
		issueID   string
		body      string
	)

	if err := args.Parse(
		model.ArgDef("account_id", &accountID),
		model.ArgDef("user", &userName),
		model.ArgDef("secret_token", &token),
		model.ArgDef("base_url", &baseURL),
		model.ArgDef("issue_id", &issueID),
		model.ArgDef("body", &body),
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

	input := &jira.Comment{
		Author: jira.User{
			AccountID: accountID,
		},
		Body: body,
	}
	comment, _, err := jiraClient.Issue.AddCommentWithContext(ctx, issueID, input)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to add comment")
	}

	return comment, nil
}
