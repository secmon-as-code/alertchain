package action

import (
	"github.com/secmon-lab/alertchain/action/bigquery"
	"github.com/secmon-lab/alertchain/action/chatgpt"
	"github.com/secmon-lab/alertchain/action/github"
	"github.com/secmon-lab/alertchain/action/http"
	"github.com/secmon-lab/alertchain/action/jira"
	"github.com/secmon-lab/alertchain/action/opsgenie"
	"github.com/secmon-lab/alertchain/action/otx"
	"github.com/secmon-lab/alertchain/action/slack"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

var actionMap = map[types.ActionName]model.RunAction{
	"github.create_issue":   github.CreateIssue,
	"github.create_comment": github.CreateComment,
	"jira.create_issue":     jira.CreateIssue,
	"jira.add_comment":      jira.AddComment,
	"jira.add_attachment":   jira.AddAttachment,
	`opsgenie.create_alert`: opsgenie.CreateAlert,
	"chatgpt.query":         chatgpt.Query,
	"slack.post":            slack.Post,
	"http.fetch":            http.Fetch,
	"otx.indicator":         otx.Indicator,
	"bigquery.insert_alert": bigquery.InsertAlert,
	"bigquery.insert_data":  bigquery.InsertData,
}

func Map() map[types.ActionName]model.RunAction {
	var copied = make(map[types.ActionName]model.RunAction, len(actionMap))
	for k, v := range actionMap {
		copied[k] = v
	}

	return copied
}
