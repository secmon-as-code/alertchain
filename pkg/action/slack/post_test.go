package slack_test

import (
	"os"
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/action/slack"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/gt"
)

func TestPost(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_SLACK_POST"); !ok {
		t.Skip("TEST_SLACK_POST is not set")
	}

	url, ok := os.LookupEnv("TEST_SLACK_WEBHOOK_URL")
	if !ok {
		t.Skip("TEST_SLACK_WEBHOOK_URL is not set")
	}
	channel, ok := os.LookupEnv("TEST_SLACK_CHANNEL")
	if !ok {
		t.Skip("TEST_SLACK_CHANNEL is not set")
	}

	alert := model.Alert{
		AlertMetaData: model.AlertMetaData{
			Title:       "test_title",
			Description: "test_description",
			Source:      "test_source",
			Attrs: []model.Attribute{
				{
					Name:  "test_pattr",
					Value: "test_value",
				},
			},
		},
		CreatedAt: time.Now(),
		Schema:    "test_schema",
		Raw:       "test_raw",
	}

	args := model.ActionArgs{
		"secret_url": url,
		"channel":    channel,
	}

	ctx := model.NewContext()
	any, err := slack.Post(ctx, alert, args)
	gt.NoError(t, err)
	gt.V(t, any).Nil()
}
