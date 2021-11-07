package alertchain_test

import (
	"context"
	"testing"
	"time"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SourceBlue struct {
	title string
	wait  time.Duration
	err   error
}

func (x *SourceBlue) Name() string { return "blue" }
func (x *SourceBlue) Run(handler alertchain.Handler) error {
	time.Sleep(x.wait)
	x.err = handler(context.Background(), &alertchain.Alert{
		Title:    x.title,
		Detector: "blue",
	})
	return nil
}

func TestSource(t *testing.T) {
	mock := db.NewDBMock(t)
	chain := alertchain.New(alertchain.OptDB(mock), alertchain.OptSources(
		&SourceBlue{title: "one", wait: time.Millisecond * 10},
		&SourceBlue{title: "five", wait: time.Millisecond * 1500},
	))

	chain.StartSources()

	alerts, err := mock.GetAlerts(newContext(), 0, 10)
	require.NoError(t, err)
	require.Len(t, alerts, 2)
	assert.Equal(t, "five", alerts[0].Title)
	assert.Equal(t, "one", alerts[1].Title)
}
