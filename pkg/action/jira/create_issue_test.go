package jira_test

import (
	_ "embed"
	"os"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/action/jira"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/gt"
)

func env(key string, dst *string) func() error {
	return func() error {
		v, ok := os.LookupEnv(key)
		if !ok {
			return goerr.New("No such env: %s", key)
		}
		*dst = v
		return nil
	}
}

func loadEnv(envs ...func() error) error {
	for _, env := range envs {
		if err := env(); err != nil {
			return err
		}
	}
	return nil
}

func TestCreateIssue(t *testing.T) {
	var (
		accountID string
		userName  string
		token     string
		baseURL   string
		project   string
	)

	if err := loadEnv(
		env("TEST_JIRA_ACCOUNT_ID", &accountID),
		env("TEST_JIRA_USER", &userName),
		env("TEST_JIRA_TOKEN", &token),
		env("TEST_JIRA_BASE_URL", &baseURL),
		env("TEST_JIRA_PROJECT", &project),
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
	}

	ctx := model.NewContext()
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
