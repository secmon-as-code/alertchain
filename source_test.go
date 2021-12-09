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
	if x.wait > 0 {
		time.Sleep(x.wait)
	}
	x.err = handler(context.Background(), &alertchain.Alert{
		Title:    x.title,
		Detector: "blue",
	})
	return nil
}

func TestSource(t *testing.T) {
	mock := db.NewDBMock(t)
	chain, err := alertchain.New(alertchain.WithDB(mock), alertchain.WithSources(
		&SourceBlue{title: "one", wait: 0},
		&SourceBlue{title: "five", wait: time.Millisecond * 900},
	))
	require.NoError(t, err)

	chain.Start()

	alerts, err := mock.GetAlerts(newContext(), 0, 10)
	require.NoError(t, err)
	require.Len(t, alerts, 2)
	assert.Equal(t, "five", alerts[0].Title)
	assert.Equal(t, "one", alerts[1].Title)
}
