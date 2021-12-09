package alertchain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/types"
	"github.com/m-mizutani/zlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newContext() *types.Context {
	var logger = zlog.New()
	return types.NewContextWith(context.Background(), logger)
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

type TaskWhite struct {
	alert *alertchain.Alert
}

func (x *TaskWhite) Name() string { return "white" }
func (x *TaskWhite) Execute(ctx *types.Context, alert *alertchain.Alert) error {
	x.alert = alert
	alert.AddAttributes([]*alertchain.Attribute{
		{
			Key:   "gamma",
			Value: "C",
		},
	})
	return nil
}

func TestChainBasic(t *testing.T) {
	mock := db.NewDBMock(t)
	chain := alertchain.New(alertchain.WithDB(mock), alertchain.WithJobs(
		&alertchain.Job{
			Tasks: []alertchain.Task{
				&TaskBlue{v: "less"},
				&TaskOrange{},
			},
		},
		&alertchain.Job{
			Tasks: []alertchain.Task{
				&TaskBlue{v: "full"},
				&TaskOrange{},
			},
		},
	))

	alert, err := chain.Execute(newContext(), &alertchain.Alert{
		Title:       "test-alert",
		Description: "x",
		Detector:    "y",
	})
	require.NoError(t, err)

	t.Run("called tasks properly", func(t *testing.T) {
		assert.Equal(t, 1, chain.Jobs()[0].Tasks[0].(*TaskBlue).called)
		assert.Equal(t, 1, chain.Jobs()[0].Tasks[1].(*TaskOrange).called)
		assert.Equal(t, 1, chain.Jobs()[1].Tasks[0].(*TaskBlue).called)
		assert.Equal(t, 1, chain.Jobs()[1].Tasks[1].(*TaskOrange).called)
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
	mock := db.NewDBMock(t)
	chain := alertchain.New(alertchain.WithDB(mock), alertchain.WithJobs(
		&alertchain.Job{
			ExitOnErr: false,
			Tasks: []alertchain.Task{
				&TaskBlue{},
				&TaskRed{err: errors.New("not failed")},
			},
		},
		&alertchain.Job{
			ExitOnErr: true,
			Tasks: []alertchain.Task{
				&TaskBlue{},
				&TaskRed{err: errors.New("failed")},
			},
		},
		&alertchain.Job{
			ExitOnErr: true,
			Tasks: []alertchain.Task{
				&TaskBlue{},
			},
		},
	))

	_, err := chain.Execute(newContext(), &alertchain.Alert{
		Title:       "test-alert",
		Description: "x",
		Detector:    "y",
	})

	require.NotErrorIs(t, err, chain.Jobs()[0].Tasks[1].(*TaskRed).err)
	require.ErrorIs(t, err, chain.Jobs()[1].Tasks[1].(*TaskRed).err)

	assert.Equal(t, 1, chain.Jobs()[0].Tasks[0].(*TaskBlue).called)
	assert.Equal(t, 1, chain.Jobs()[0].Tasks[1].(*TaskRed).called)
	assert.Equal(t, 1, chain.Jobs()[1].Tasks[0].(*TaskBlue).called)
	assert.Equal(t, 1, chain.Jobs()[1].Tasks[1].(*TaskRed).called)
	assert.Equal(t, 0, chain.Jobs()[2].Tasks[0].(*TaskBlue).called) // Stopped by TaskRed's error
}

func TestChainTimeout(t *testing.T) {
	mock := db.NewDBMock(t)
	chain := alertchain.New(alertchain.WithDB(mock), alertchain.WithJobs(
		&alertchain.Job{
			ExitOnErr: true,
			Timeout:   time.Millisecond * 500,
			Tasks: []alertchain.Task{
				&TaskGray{},
			},
		},
		&alertchain.Job{
			ExitOnErr: true,
			Timeout:   time.Millisecond * 100,
			Tasks: []alertchain.Task{
				&TaskGray{},
			},
		},
		&alertchain.Job{
			Tasks: []alertchain.Task{
				&TaskBlue{},
			},
		},
	))

	_, err := chain.Execute(newContext(), &alertchain.Alert{
		Title:    "test-alert",
		Detector: "y",
	})
	require.NoError(t, err) // Not error by TaskGray even if timeout. But task should return error when timed out
	assert.Equal(t, 1, chain.Jobs()[0].Tasks[0].(*TaskGray).called)
	assert.False(t, chain.Jobs()[0].Tasks[0].(*TaskGray).timeout)
	assert.Equal(t, 1, chain.Jobs()[1].Tasks[0].(*TaskGray).called)
	assert.True(t, chain.Jobs()[1].Tasks[0].(*TaskGray).timeout)
	assert.Equal(t, 1, chain.Jobs()[2].Tasks[0].(*TaskBlue).called)
}

func TestChainAlert(t *testing.T) {
	task := &TaskWhite{}
	mock := db.NewDBMock(t)
	chain := alertchain.New(alertchain.WithDB(mock), alertchain.WithJobs(
		&alertchain.Job{
			Tasks: []alertchain.Task{task},
		}))

	sent := &alertchain.Alert{
		Title:       "words",
		Detector:    "blue",
		Description: "five",
		DetectedAt:  time.Now().UTC(),
		Attributes: alertchain.Attributes{
			{
				Key:   "alpha",
				Value: "A",
			},
			{
				Key:   "beta",
				Value: "B",
			},
		},
		References: alertchain.References{
			{
				Source:  "orange",
				Title:   "scarred red",
				URL:     "https://example.com/x",
				Comment: "?",
			},
		},
	}
	created, err := chain.Execute(newContext(), sent)
	require.NoError(t, err)
	require.NotNil(t, task.alert)
	assert.Equal(t, sent.Title, task.alert.Title)
	assert.Equal(t, sent.Detector, task.alert.Detector)
	assert.Equal(t, sent.Description, task.alert.Description)
	assert.NotEqual(t, sent.DetectedAt, task.alert.DetectedAt)            // alert in task is generated only unixtime second
	assert.Equal(t, sent.DetectedAt.Unix(), task.alert.DetectedAt.Unix()) // matched with only unixtime second

	assert.Len(t, sent.Attributes.FindByKey("alpha").FindByValue("A"), 1)
	assert.Len(t, sent.Attributes.FindByKey("beta").FindByValue("B"), 1)
	assert.Len(t, sent.Attributes.FindByKey("gamma").FindByValue("C"), 0)
	assert.Len(t, task.alert.Attributes.FindByKey("alpha").FindByValue("A"), 1)
	assert.Len(t, task.alert.Attributes.FindByKey("beta").FindByValue("B"), 1)
	assert.Len(t, task.alert.Attributes.FindByKey("gamma").FindByValue("C"), 0)

	// created alert has added attribute in last job
	assert.Len(t, created.Attributes.FindByKey("gamma").FindByValue("C"), 1)

	got, err := mock.GetAlert(newContext(), created.ID())
	require.NoError(t, err)
	require.NotNil(t, got)
	// retrieved alert also has added attribute
	assert.Len(t, alertchain.NewAlert(got).Attributes.FindByKey("gamma").FindByValue("C"), 1)

	require.Len(t, task.alert.References, 1)
	assert.Equal(t, sent.References[0], task.alert.References[0])
}
