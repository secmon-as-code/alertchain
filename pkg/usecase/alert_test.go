package usecase_test

import (
	"sync"
	"testing"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type blockTask struct {
	wg sync.WaitGroup
}

func (x *blockTask) Name() string                              { return "blocker" }
func (x *blockTask) Description() string                       { return "blocking" }
func (x *blockTask) IsExecutable(alert *alertchain.Alert) bool { return false }
func (x *blockTask) Execute(alert *alertchain.Alert) error {
	x.wg.Done()
	return nil
}

func setupAlertTest(t *testing.T) (usecase.Interface, infra.Clients, *alertchain.Chain) {
	chain := &alertchain.Chain{}

	clients := infra.Clients{
		DB: db.NewDBMock(t),
	}
	uc := usecase.New(clients, chain)

	return uc, clients, chain
}

func TestRecvAlert(t *testing.T) {
	uc, clients, chain := setupAlertTest(t)
	stage := chain.NewStage()
	blocker := &blockTask{}
	blocker.wg.Add(1)
	stage.AddTask(blocker)

	input := alertchain.Alert{
		Alert: ent.Alert{
			Title:    "five",
			Detector: "blue",
		},
	}
	alert, err := uc.RecvAlert(&input)
	require.NoError(t, err)
	require.NotNil(t, alert)

	blocker.wg.Wait()

	got, err := clients.DB.GetAlert(alert.ID)
	require.NoError(t, err)
	assert.Equal(t, alert.Title, got.Title)
}
