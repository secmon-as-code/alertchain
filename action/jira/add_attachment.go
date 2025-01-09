package jira

import (
	"context"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
)

func AddAttachment(ctx context.Context, alert model.Alert, args model.ActionArgs) (any, error) {
	var (
		accountID string
		userName  string
		token     string
		baseURL   string
		issueID   string
		fileName  string
		data      string
	)

	if err := args.Parse(
		model.ArgDef("account_id", &accountID),
		model.ArgDef("user", &userName),
		model.ArgDef("secret_token", &token),
		model.ArgDef("base_url", &baseURL),
		model.ArgDef("issue_id", &issueID),
		model.ArgDef("file_name", &fileName),
		model.ArgDef("data", &data),
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

	body := strings.NewReader(data)
	attach, _, err := jiraClient.Issue.PostAttachmentWithContext(ctx, issueID, body, fileName)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to post attachment")
	}

	return attach, nil
}
