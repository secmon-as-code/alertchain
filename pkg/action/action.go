package action

import (
	"github.com/m-mizutani/alertchain/pkg/action/chatgpt"
	"github.com/m-mizutani/alertchain/pkg/action/github"
	"github.com/m-mizutani/alertchain/pkg/domain/interfaces"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

var actionMap = map[types.ActionName]interfaces.RunAction{
	"github.create_issue":   github.CreateIssue,
	"chatgpt.comment_alert": chatgpt.CommentAlert,
}

func Map() map[types.ActionName]interfaces.RunAction {
	var copied = make(map[types.ActionName]interfaces.RunAction, len(actionMap))
	for k, v := range actionMap {
		copied[k] = v
	}

	return copied
}
