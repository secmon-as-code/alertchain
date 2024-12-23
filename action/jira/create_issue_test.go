package jira_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/alertchain/action/jira"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/utils"
)

func TestCreateIssue(t *testing.T) {
	var (
		accountID string
		userName  string
		token     string
		baseURL   string
		project   string
	)

	if err := utils.LoadEnv(
		utils.EnvDef("TEST_JIRA_ACCOUNT_ID", &accountID),
		utils.EnvDef("TEST_JIRA_USER", &userName),
		utils.EnvDef("TEST_JIRA_TOKEN", &token),
		utils.EnvDef("TEST_JIRA_BASE_URL", &baseURL),
		utils.EnvDef("TEST_JIRA_PROJECT", &project),
	); err != nil {
		t.Skipf("Skip test due to missing env: %v", err)
	}

	args := model.ActionArgs{
		"account_id":   accountID,
		"user":         userName,
		"secret_token": token,
		"base_url":     baseURL,
		"project":      project,
		"issue_type":   "Task",
		"labels":       []string{"test"},
		"assignee":     accountID,
	}

	ctx := context.Background()
	alert := model.NewAlert(model.AlertMetaData{
		Title:       "Alert testing",
		Description: "This is test alert",
		Source:      "test_source",
		Attrs: model.Attributes{
			{
				ID:    "my_id",
				Key:   "my_key",
				Value: "my_value",
			},
		},
		Refs: model.References{
			{
				Title: "my_ref_title",
				URL:   "https://example.com",
			},
		},
	}, "test_alert", struct {
		MyRecord string
	}{
		MyRecord: "my_record",
	})

	issue := gt.R1(jira.CreateIssue(ctx, alert, args)).NoError(t)
	gt.V(t, issue).NotNil()
}

func TestTemplate(t *testing.T) {
	alert := model.NewAlert(model.AlertMetaData{
		Title:       "my test alert",
		Description: "test_description",
		Source:      "test_source",
		Attrs: model.Attributes{
			{
				ID:    "my_id",
				Key:   "my_key",
				Value: "my_value",
			},
		},
		Refs: model.References{
			{
				Title: "my_ref_title",
				URL:   "my_ref_url",
			},
		},
	}, "test_alert", struct {
		MyRecord string
	}{
		MyRecord: "my_record",
	})

	buf := gt.R1(jira.ExecTemplate(alert)).NoError(t)
	gt.S(t, buf).
		NotContains("my test alert"). // should not contain title
		Contains("test_description").
		NotContains("my_id").
		Contains("my_key"). // should contain attribute key
		Contains("my_value").
		Contains("my_ref_title").
		Contains("my_ref_url")
}
