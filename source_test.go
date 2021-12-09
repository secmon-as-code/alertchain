package alertchain_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SourceBlue struct {
	title string
	err   error
}

func (x *SourceBlue) Name() string { return "blue" }
func (x *SourceBlue) Run(handler alertchain.Handler) error {
	x.err = handler(context.Background(), &alertchain.Alert{
		Title:    x.title,
		Detector: "blue",
	})
	return nil
}

func TestSource(t *testing.T) {
	mock := db.NewDBMock(t)
	chain, err := alertchain.New(alertchain.WithDB(mock), alertchain.WithSources(
		&SourceBlue{title: "one"},
		&SourceBlue{title: "five"},
	))
	require.NoError(t, err)

	chain.Start()

	alerts, err := mock.GetAlerts(newContext(), 0, 10)
	require.NoError(t, err)
	require.Len(t, alerts, 2)
	assert.Contains(t, []string{alerts[0].Title, alerts[1].Title}, "one")
	assert.Contains(t, []string{alerts[0].Title, alerts[1].Title}, "five")
}
