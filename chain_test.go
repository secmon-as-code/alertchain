package alertchain_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/utils"
	"github.com/m-mizutani/alertchain/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	if logLevel, ok := os.LookupEnv("LOG_LEVEL"); ok {
		utils.Logger.SetLogLevel(logLevel)
	}
}

type TaskBlue struct {
	called int
	v      string
}

func (x *TaskBlue) Name() string { return "blue" }
func (x *TaskBlue) Execute(ctx *types.Context, alert *alertchain.Alert) error {
	x.called++
	alert.AddAttributes([]*alertchain.Attribute{
		{
			Key:   "blue",
			Value: "time-" + x.v,
		},
	})
	return nil
}

type TaskOrange struct {
	called int
}

func (x *TaskOrange) Name() string { return "orange" }
func (x *TaskOrange) Execute(ctx *types.Context, alert *alertchain.Alert) error {
	x.called++
	for _, attr := range alert.Attributes.FindByKey("blue") {
		attr.Annotate(&alertchain.Annotation{
			Timestamp: time.Now(),
			Source:    "orange",
			Name:      "circit",
			Value:     "good",
		})
	}

	return nil
}

type TaskRed struct {
	called int
	err    error
}

func (x *TaskRed) Name() string { return "red" }
func (x *TaskRed) Execute(ctx *types.Context, alert *alertchain.Alert) error {
	x.called++
	return x.err
}

type TaskGray struct {
	called  int
	timeout bool
}

func (x *TaskGray) Name() string { return "gray" }
func (x *TaskGray) Execute(ctx *types.Context, alert *alertchain.Alert) error {
	x.called++
	select {
	case <-ctx.Done():
		x.timeout = true
	case <-time.After(time.Millisecond * 200):
		x.timeout = false
	}
	return nil
}

func TestChainBasic(t *testing.T) {
	chain := alertchain.New(db.NewDBMock(t))
	chain.AddJobs(&alertchain.Job{
		Tasks: []alertchain.Task{
			&TaskBlue{v: "less"},
			&TaskOrange{},
		},
	}, &alertchain.Job{
		Tasks: []alertchain.Task{
			&TaskBlue{v: "full"},
			&TaskOrange{},
		},
	})

	alert, err := chain.Execute(types.NewContext(), &alertchain.Alert{
		Title:       "test-alert",
		Description: "x",
		Detector:    "y",
	})
	require.NoError(t, err)

	t.Run("called tasks properly", func(t *testing.T) {
		assert.Equal(t, 1, chain.Jobs[0].Tasks[0].(*TaskBlue).called)
		assert.Equal(t, 1, chain.Jobs[0].Tasks[1].(*TaskOrange).called)
		assert.Equal(t, 1, chain.Jobs[1].Tasks[0].(*TaskBlue).called)
		assert.Equal(t, 1, chain.Jobs[1].Tasks[1].(*TaskOrange).called)
	})

	t.Run("added attrs", func(t *testing.T) {
		assert.Len(t, alert.Attributes.FindByKey("blue"), 2)
		assert.Len(t, alert.Attributes.FindByKey("blue").FindByValue("time-full"), 1)
		assert.Len(t, alert.Attributes.FindByKey("blue").FindByValue("time-less"), 1)
	})

	t.Run("added annotation", func(t *testing.T) {
		attrs := alert.Attributes.FindByKey("blue").FindByValue("time-less")
		require.Len(t, attrs, 1)
		require.Len(t, attrs[0].Annotations, 1)
		assert.Equal(t, "orange", attrs[0].Annotations[0].Source)
		assert.Equal(t, "circit", attrs[0].Annotations[0].Name)
		assert.Equal(t, "good", attrs[0].Annotations[0].Value)
	})
}

func TestChainError(t *testing.T) {
	chain := alertchain.New(db.NewDBMock(t))
	chain.AddJobs(alertchain.Jobs{
		{
			ExitOnErr: false,
			Tasks: []alertchain.Task{
				&TaskBlue{},
				&TaskRed{err: errors.New("not failed")},
			},
		},
		{
			ExitOnErr: true,
			Tasks: []alertchain.Task{
				&TaskBlue{},
				&TaskRed{err: errors.New("failed")},
			},
		},
		{
			ExitOnErr: true,
			Tasks: []alertchain.Task{
				&TaskBlue{},
			},
		},
	}...)

	_, err := chain.Execute(types.NewContext(), &alertchain.Alert{
		Title:       "test-alert",
		Description: "x",
		Detector:    "y",
	})

	require.NotErrorIs(t, err, chain.Jobs[0].Tasks[1].(*TaskRed).err)
	require.ErrorIs(t, err, chain.Jobs[1].Tasks[1].(*TaskRed).err)

	assert.Equal(t, 1, chain.Jobs[0].Tasks[0].(*TaskBlue).called)
	assert.Equal(t, 1, chain.Jobs[0].Tasks[1].(*TaskRed).called)
	assert.Equal(t, 1, chain.Jobs[1].Tasks[0].(*TaskBlue).called)
	assert.Equal(t, 1, chain.Jobs[1].Tasks[1].(*TaskRed).called)
	assert.Equal(t, 0, chain.Jobs[2].Tasks[0].(*TaskBlue).called) // Stopped by TaskRed's error
}

func TestChainTimeout(t *testing.T) {
	chain := alertchain.New(db.NewDBMock(t))
	chain.AddJobs(alertchain.Jobs{
		{
			ExitOnErr: true,
			Timeout:   time.Millisecond * 500,
			Tasks: []alertchain.Task{
				&TaskGray{},
			},
		},
		{
			ExitOnErr: true,
			Timeout:   time.Millisecond * 100,
			Tasks: []alertchain.Task{
				&TaskGray{},
			},
		},
		{
			Tasks: []alertchain.Task{
				&TaskBlue{},
			},
		},
	}...)

	_, err := chain.Execute(types.NewContext(), &alertchain.Alert{
		Title:    "test-alert",
		Detector: "y",
	})
	require.NoError(t, err) // Not error by TaskGray even if timeout. But task should return error when timed out
	assert.Equal(t, 1, chain.Jobs[0].Tasks[0].(*TaskGray).called)
	assert.False(t, chain.Jobs[0].Tasks[0].(*TaskGray).timeout)
	assert.Equal(t, 1, chain.Jobs[1].Tasks[0].(*TaskGray).called)
	assert.True(t, chain.Jobs[1].Tasks[0].(*TaskGray).timeout)
	assert.Equal(t, 1, chain.Jobs[2].Tasks[0].(*TaskBlue).called)
}
