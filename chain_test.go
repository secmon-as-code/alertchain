package alertchain_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/alertchain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TaskBlue struct {
	v int
}

func (x *TaskBlue) Name() string { return "blue" }
func (x *TaskBlue) Execute(ctx context.Context, alert *alertchain.Alert) error {
	return nil
}

type TaskOrange struct{}

func (x *TaskOrange) Name() string { return "orange" }
func (x *TaskOrange) Execute(ctx context.Context, alert *alertchain.Alert) error {
	return nil
}

func TestChain(t *testing.T) {
	chain := &alertchain.Chain{
		Jobs: []*alertchain.Job{
			{
				Tasks: []alertchain.Task{
					&TaskBlue{v: 1},
					&TaskOrange{},
				},
			},
			{
				Tasks: []alertchain.Task{
					&TaskBlue{v: 2},
				},
			},
		},
	}

	t.Run("polite way", func(t *testing.T) {
		task := chain.LookupTask(&TaskBlue{})
		require.NotNil(t, task)
		assert.Equal(t, "blue", task.Name())
		blue, ok := task.(*TaskBlue)
		require.True(t, ok)

		t.Run("should lookup first task of the type", func(t *testing.T) {
			assert.Equal(t, 1, blue.v)
		})
	})

	t.Run("like sugar syntax", func(t *testing.T) {
		task := chain.LookupTask(&TaskBlue{}).(*TaskBlue)
		require.NotNil(t, task)
	})
}
