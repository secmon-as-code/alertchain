package chatgpt_test

import (
	_ "embed"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/action/chatgpt"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/gt"
	"github.com/sashabaranov/go-openai"
)

//go:embed testdata/aws_guardduty_example.json
var alertData []byte

func TestAnalystInquiry(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_CHATGPT_ANALYST"); !ok {
		t.Skip("Skipping test because TEST_CHATGPT_ANALYST is not set")
	}

	cfg := model.ActionArgs{
		"api_key": os.Getenv("TEST_CHATGPT_API_KEY"),
	}

	requiredVars := []string{"api_key"}
	for _, key := range requiredVars {
		gt.V(t, cfg[key]).NotEqual("")
	}

	var body any
	gt.NoError(t, json.Unmarshal(alertData, &body))

	ctx := model.NewContext()
	alert := model.Alert{
		AlertMetaData: model.AlertMetaData{
			Title:       "test",
			Description: "test",
		},
		Data:      body,
		CreatedAt: time.Now(),
	}

	resp := gt.R1(chatgpt.CommentAlert(ctx, alert, model.ActionArgs{})).NoError(t)
	data := gt.Cast[openai.ChatCompletionResponse](t, resp)
	gt.A(t, data.Choices).Length(1)
	t.Log(data.Choices[0].Message.Content)
}
