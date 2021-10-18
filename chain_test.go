package alertchain_test

import (
	"testing"
	"time"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TaskBlue struct {
	v string
}

func (x *TaskBlue) Name() string { return "blue" }
func (x *TaskBlue) Execute(ctx *types.Context, alert *alertchain.Alert) error {
	alert.AddAttributes([]*alertchain.Attribute{
		{
			Key:   "blue",
			Value: "time-" + x.v,
		},
	})
	return nil
}

type TaskOrange struct{}

func (x *TaskOrange) Name() string { return "orange" }
func (x *TaskOrange) Execute(ctx *types.Context, alert *alertchain.Alert) error {
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

func TestChain(t *testing.T) {
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
