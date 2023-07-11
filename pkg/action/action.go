package action

import (
	"github.com/m-mizutani/alertchain/pkg/action/chatgpt"
	"github.com/m-mizutani/alertchain/pkg/action/github"
	"github.com/m-mizutani/alertchain/pkg/action/http"
	"github.com/m-mizutani/alertchain/pkg/action/otx"
	"github.com/m-mizutani/alertchain/pkg/action/slack"
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

var actionMap = map[types.ActionName]interfaces.RunAction{
	"github.create_issue":   github.CreateIssue,
	"github.create_comment": github.CreateComment,
	"chatgpt.query":         chatgpt.Query,
	"slack.post":            slack.Post,
	"http.fetch":            http.Fetch,
	"otx.indicator":         otx.Indicator,
}

func Map() map[types.ActionName]interfaces.RunAction {
	var copied = make(map[types.ActionName]interfaces.RunAction, len(actionMap))
	for k, v := range actionMap {
		copied[k] = v
	}

	return copied
}
