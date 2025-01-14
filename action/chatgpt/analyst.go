package chatgpt

import (
	"context"
	_ "embed"
	"encoding/json"

	openai "github.com/sashabaranov/go-openai"

	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/utils"
)

func Query(ctx context.Context, alert model.Alert, args model.ActionArgs) (any, error) {
	apiKey, ok := args["secret_api_key"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "secret_api_key is required")
	}

	client := openai.NewClient(apiKey)

	data, err := json.Marshal(alert.Data)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to marshal alert data")
	}

	prompt := "Please analyze and summarize the given JSON-formatted security alert data, and suggest appropriate actions for the security administrator to respond to the alert: " + string(data)

	if v, ok := args["prompt"].(string); ok {
		prompt = v
	}

	if ctxutil.IsDryRun(ctx) {
		return nil, nil
	}

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to call OpenAI API")
	}

	return utils.ToAny(resp)
}
