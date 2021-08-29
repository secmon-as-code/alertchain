package alertchain_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/usecase"
	"github.com/m-mizutani/alertchain/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAlertTest(t *testing.T) (usecase.Interface, infra.Clients, *alertchain.Chain) {
	chain := &alertchain.Chain{}

	clients := infra.Clients{
		DB: db.NewDBMock(t),
	}
	uc := usecase.New(clients, chain)

	return uc, clients, chain
}

type mock struct {
	Exec func(alert *alertchain.Alert) error
}

func (x *mock) Name() string { return "mock" }
func (x *mock) Execute(ctx context.Context, alert *alertchain.Alert) error {
	return x.Exec(alert)
}

func TestRecvAlert(t *testing.T) {
	uc, clients, chain := setupAlertTest(t)

	var done bool
	chain.NewJob().AddTask(&mock{
		Exec: func(alert *alertchain.Alert) error {
			alert.UpdateSeverity(types.SevAffected)
			alert.UpdateStatus(types.StatusClosed)
			done = true
			return nil
		},
	})

	input := alertchain.Alert{
		Alert: ent.Alert{
			Title:    "five",
			Detector: "blue",
		},
	}
	ctx, wg := alertchain.SetWaitGroupToCtx(context.Background())
	alert, err := uc.RecvAlert(ctx, &input)
	require.NoError(t, err)
	require.NotNil(t, alert)

	wg.Wait()
	assert.True(t, done)

	got, err := clients.DB.GetAlert(context.Background(), alert.ID)
	require.NoError(t, err)
	assert.Equal(t, alert.Title, got.Title)
	assert.Equal(t, types.SevAffected, got.Severity)
	assert.Equal(t, types.StatusClosed, got.Status)
}

func TestRecvAlertDoNotUpdate(t *testing.T) {
	t.Run("do not update severity and status by overwriting vars", func(t *testing.T) {
		uc, clients, chain := setupAlertTest(t)

		var done bool
		chain.NewJob().AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error {
				alert.Severity = types.SevAffected
				alert.Status = types.StatusClosed
				if err := alert.Commit(context.Background()); err != nil {
					return err
				}
				done = true
				return nil
			},
		})

		input := alertchain.Alert{
			Alert: ent.Alert{
				Title:    "five",
				Detector: "blue",
			},
		}
		ctx, wg := alertchain.SetWaitGroupToCtx(context.Background())
		alert, err := uc.RecvAlert(ctx, &input)
		require.NoError(t, err)
		require.NotNil(t, alert)

		wg.Wait()
		assert.True(t, done)

		got, err := clients.DB.GetAlert(context.Background(), alert.ID)
		require.NoError(t, err)
		assert.Equal(t, alert.Title, got.Title)
		assert.NotEqual(t, types.SevAffected, got.Severity)
		assert.NotEqual(t, types.StatusClosed, got.Status)
	})
}

func TestRecvAlertMassiveAnnotation(t *testing.T) {
	const multiplex = 32

	uc, clients, chain := setupAlertTest(t)

	job := chain.NewJob()
	job.Timeout = time.Second
	for i := 0; i < multiplex; i++ {
		job.AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error {
				require.Len(t, alert.Attributes, 1)
				alert.Attributes[0].Annotate(&alertchain.Annotation{
					Annotation: ent.Annotation{
						Source:    "x",
						Timestamp: rand.Int63(), /* nosec */
						Name:      "y",
						Value:     "z",
					},
				})
				return nil
			},
		})
	}

	ctx, wg := alertchain.SetWaitGroupToCtx(context.Background())
	input := alertchain.Alert{
		Alert: ent.Alert{
			Title:    "five",
			Detector: "blue",
		},
		Attributes: []*alertchain.Attribute{
			{
				Attribute: ent.Attribute{
					Key:   "color",
					Value: "red",
					Type:  types.AttrUserID,
				},
			},
		},
	}
	created, err := uc.RecvAlert(ctx, &input)
	require.NoError(t, err)
	wg.Wait()

	alert, err := clients.DB.GetAlert(context.Background(), created.Alert.ID)
	require.NoError(t, err)
	require.Len(t, alert.Edges.Attributes[0].Edges.Annotations, multiplex)
	for _, ann := range alert.Edges.Attributes[0].Edges.Annotations {
		assert.Equal(t, "x", ann.Source)
		assert.Equal(t, "y", ann.Name)
		assert.Equal(t, "z", ann.Value)
		assert.Greater(t, ann.Timestamp, int64(0))
	}
}

func TestRecvAlertErrorHandling(t *testing.T) {
	t.Run("exit on error", func(t *testing.T) {
		uc, _, chain := setupAlertTest(t)

		job := chain.NewJob()
		job.ExitOnErr = true
		job.AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error { return nil },
		})
		job.AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error { return errors.New("bomb!") },
		})

		done2ndJob := false
		chain.NewJob().AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error {
				done2ndJob = true
				return nil
			},
		})

		input := alertchain.Alert{
			Alert: ent.Alert{
				Title:    "five",
				Detector: "blue",
			},
		}
		ctx, wg := alertchain.SetWaitGroupToCtx(context.Background())
		_, err := uc.RecvAlert(ctx, &input)
		require.NoError(t, err)
		wg.Wait()
		assert.False(t, done2ndJob)
	})

	t.Run("not exit on error", func(t *testing.T) {
		uc, _, chain := setupAlertTest(t)

		job := chain.NewJob()
		// Default: job.ExitOnErr = false
		job.AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error { return nil },
		})
		job.AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error { return errors.New("bomb!") },
		})

		done2ndJob := false
		chain.NewJob().AddTask(&mock{
			Exec: func(alert *alertchain.Alert) error {
				done2ndJob = true
				return nil
			},
		})

		input := alertchain.Alert{
			Alert: ent.Alert{
				Title:    "five",
				Detector: "blue",
			},
		}
		ctx, wg := alertchain.SetWaitGroupToCtx(context.Background())
		_, err := uc.RecvAlert(ctx, &input)
		require.NoError(t, err)
		wg.Wait()
		assert.True(t, done2ndJob)
	})
}
