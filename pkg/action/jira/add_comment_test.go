package jira_test

import (
	"testing"

	go_jira "github.com/andygrunwald/go-jira"

	"github.com/m-mizutani/alertchain/pkg/action/jira"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/gt"
)

func TestAddComment(t *testing.T) {
	var (
		accountID string
		userName  string
		token     string
		baseURL   string
		issueID   string
	)

	if err := utils.LoadEnv(
		utils.Env("TEST_JIRA_ACCOUNT_ID", &accountID),
		utils.Env("TEST_JIRA_USER", &userName),
		utils.Env("TEST_JIRA_TOKEN", &token),
		utils.Env("TEST_JIRA_BASE_URL", &baseURL),
		utils.Env("TEST_JIRA_ISSUE_ID", &issueID),
	); err != nil {
		t.Skipf("Skip test due to missing env: %v", err)
	}

	args := model.ActionArgs{
		"account_id":   accountID,
		"user":         userName,
		"secret_token": token,
		"base_url":     baseURL,
		"issue_id":     issueID,
		"body":         "test comment",
	}

	ctx := model.NewContext()
	alert := model.NewAlert(model.AlertMetaData{
		Title:       "test_alert",
		Description: "test_description",
		Source:      "test_source",
	}, "test_alert", struct{}{})

	ret := gt.R1(jira.AddComment(ctx, alert, args)).NoError(t)
	comment := gt.Cast[*go_jira.Comment](t, ret)
	gt.Equal(t, comment.Author.EmailAddress, userName)
}
