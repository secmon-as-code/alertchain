package slack

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/alertchain/pkg/ctxutil"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/slack-go/slack"
)

type notifyContents struct {
	Text   string
	Body   string
	Color  string
	Fields []*notifyField
	Raw    string
}

type notifyField struct {
	Name  string
	Value string
	URL   string
}

// Post is a function to post message to Slack via incoming webhook
func Post(ctx context.Context, alert model.Alert, args model.ActionArgs) (any, error) {
	notify := &notifyContents{
		Text: "Notification from AlertChain",
		Body: fmt.Sprintf("*%s*\n%s", alert.Title, alert.Description),
		Raw:  alert.Raw,
		Fields: []*notifyField{
			{
				Name:  "schema",
				Value: string(alert.Schema),
			},
			{
				Name:  "source",
				Value: alert.Source,
			},
			{
				Name:  "created at",
				Value: alert.CreatedAt.Format("2006-01-02 15:04:05 MST"),
			},
		},
	}

	url, ok := args["secret_url"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "url is required")
	}
	channel, ok := args["channel"].(string)
	if !ok {
		return nil, goerr.Wrap(types.ErrActionInvalidArgument, "channel is required")
	}

	if v, ok := args["text"].(string); ok {
		notify.Text = v
	}
	if v, ok := args["body"].(string); ok {
		notify.Body = v
	}
	if v, ok := args["color"].(string); ok {
		notify.Color = v
	}

	for _, attr := range alert.Attrs {
		notify.Fields = append(notify.Fields, &notifyField{
			Name:  string(attr.Key),
			Value: fmt.Sprintf("%v", attr.Value),
		})
	}

	msg := buildSlackMessage(notify, alert)
	msg.Channel = channel

	if ctxutil.IsDryRun(ctx) {
		return nil, nil
	}

	if err := slack.PostWebhookContext(ctx, url, msg); err != nil {
		raw, _ := json.Marshal(msg)
		return nil, goerr.Wrap(err, "failed to post slack message").With("body", string(raw))
	}

	return nil, nil
}

func buildSlackMessage(notify *notifyContents, _ interface{}) *slack.WebhookMessage {
	color := "#2EB67D"
	if notify.Color != "" {
		color = notify.Color
	}

	var blocks []slack.Block

	if notify.Body != "" {
		blocks = append(blocks,
			slack.NewDividerBlock(),
			slack.NewSectionBlock(&slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: notify.Body,
			}, nil, nil),
		)
	}

	var customFields []*slack.TextBlockObject
	for _, field := range notify.Fields {
		customFields = append(customFields, toBlock(field))
	}
	if len(customFields) > 0 {
		blocks = append(blocks,
			slack.NewDividerBlock(),
			slack.NewSectionBlock(nil, customFields, nil),
		)
	}

	raw := notify.Raw
	if len(raw) > 2990 {
		raw = raw[:2990] + "..."
	}

	blocks = append(blocks,
		slack.NewDividerBlock(),
		slack.NewSectionBlock(&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: fmt.Sprintf("```%s```", raw),
		}, nil, nil),
	)

	msg := &slack.WebhookMessage{
		Text: notify.Text,
		Attachments: []slack.Attachment{
			{
				Color: color,
				Blocks: slack.Blocks{
					BlockSet: blocks,
				},
			},
		},
	}

	return msg
}

func toBlock(field *notifyField) *slack.TextBlockObject {
	text := fmt.Sprintf("*%s*: ", field.Name)
	if field.URL != "" {
		text += fmt.Sprintf("<%s|%s>", field.URL, field.Value)
	} else {
		text += field.Value
	}
	return slack.NewTextBlockObject(slack.MarkdownType, text, false, false)
}
