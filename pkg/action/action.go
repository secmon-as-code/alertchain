package action

import (
	"github.com/secmon-lab/alertchain/pkg/action/bigquery"
	"github.com/secmon-lab/alertchain/pkg/action/chatgpt"
	"github.com/secmon-lab/alertchain/pkg/action/github"
	"github.com/secmon-lab/alertchain/pkg/action/http"
	"github.com/secmon-lab/alertchain/pkg/action/jira"
	"github.com/secmon-lab/alertchain/pkg/action/opsgenie"
	"github.com/secmon-lab/alertchain/pkg/action/otx"
	"github.com/secmon-lab/alertchain/pkg/action/slack"
	"github.com/secmon-lab/alertchain/pkg/domain/interfaces"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

var actionMap = map[types.ActionName]interfaces.RunAction{
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

func Map() map[types.ActionName]interfaces.RunAction {
	var copied = make(map[types.ActionName]interfaces.RunAction, len(actionMap))
	for k, v := range actionMap {
		copied[k] = v
	}

	return copied
}
