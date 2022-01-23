package actions_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/m-mizutani/alertchain/pkg/actions"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyAction struct {
	cfg model.ActionConfig
}

func (x *dummyAction) Run(ctx *types.Context, alert *model.Alert, args ...*model.Attribute) (*model.ChangeRequest, error) {
	return nil, nil
}

func TestRepository(t *testing.T) {
	t.Run("register new action and construct", func(t *testing.T) {
		id := "dummy_" + uuid.NewString()
		var called int
		actions.Register(id, func(config model.ActionConfig) (model.Action, error) {
			called++
			return &dummyAction{cfg: config}, nil
		})

		cfg := model.ActionConfig{
			"color": "blue",
		}
		action, err := actions.New(id, cfg)
		require.NoError(t, err)
		assert.Equal(t, 1, called)
		dummy, ok := action.(*dummyAction)
		require.True(t, ok)
		assert.Equal(t, "blue", dummy.cfg["color"])
	})

	t.Run("fail if trying to register duplicated id", func(t *testing.T) {
		id := "dummy_" + uuid.NewString()
		actions.Register(id, func(config model.ActionConfig) (model.Action, error) {
			return &dummyAction{cfg: config}, nil
		})

		assert.Panics(t, func() {
			// can not register with same id twice
			actions.Register(id, func(config model.ActionConfig) (model.Action, error) {
				return nil, nil
			})
		})
	})
}
