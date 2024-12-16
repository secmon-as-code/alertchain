package core

import (
	"github.com/secmon-lab/alertchain/pkg/domain/interfaces"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

func (x *Core) GetAction(name types.ActionName) (interfaces.RunAction, bool) {
	action, ok := x.actionMap[name]
	return action, ok
}
