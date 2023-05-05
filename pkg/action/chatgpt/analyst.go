package chatgpt

import (
	_ "embed"
	"encoding/json"

	openai "github.com/sashabaranov/go-openai"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

func CommentAlert(ctx *model.Context, alert model.Alert, args model.ActionArgs) (any, error) {
	apiKey, ok := args["api_key"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidConfig, "api_key is required")
	}

	client := openai.NewClient(apiKey)

	prompt := "Summarize the following json formatted data of security alert and propose security administrator's action: "
	if v, ok := args["prompt"].(string); ok {
		prompt = v
	}

	data, err := json.Marshal(alert.Data)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to marshal alert data")
	}

	if ctx.DryRun() {
		return nil, nil
	}

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt + string(data),
				},
			},
		},
	)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to call OpenAI API")
	}

	return resp, nil
}
