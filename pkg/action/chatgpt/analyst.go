package chatgpt

import (
	_ "embed"
	"encoding/json"

	openai "github.com/sashabaranov/go-openai"

	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/goerr"
)

type Analyst struct {
	id     types.ActionID
	client *openai.Client
	prompt string
}

type AnalystFactory struct{}

func (x *AnalystFactory) Name() types.ActionName {
	return "chatgpt-analyst"
}

func (x *AnalystFactory) New(id types.ActionID, cfg model.ActionConfigValues) (interfaces.Action, error) {
	apiKey, ok := cfg["api_key"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidConfig, "api_key is required")
	}

	client := openai.NewClient(apiKey)
	action := &Analyst{
		id:     id,
		client: client,
		prompt: "Summarize the following json formatted data of security alert and propose security administrator's action: ",
	}

	if v, ok := cfg["prompt"].(string); ok {
		action.prompt = v
	}

	return action, nil
}

func (x *Analyst) ID() types.ActionID { return x.id }

func (x *Analyst) Run(ctx *model.Context, alert model.Alert, params model.ActionArgs) (any, error) {
	data, err := json.Marshal(alert.Data)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to marshal alert data")
	}

	resp, err := x.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: x.prompt + string(data),
				},
			},
		},
	)
	if err != nil {
		return nil, goerr.Wrap(err, "Failed to call OpenAI API")
	}

	return resp, nil
}
