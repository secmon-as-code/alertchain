package core

import (
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

func (x *Core) GetAction(name types.ActionName) (model.RunAction, bool) {
	action, ok := x.actionMap[name]
	return action, ok
}
