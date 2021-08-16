package usecase_test

import (
	"testing"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAlertTest(t *testing.T) (usecase.Usecase, infra.DBClient) {
	chain := &alertchain.Chain{
		Stages: []alertchain.Tasks{
			{},
		},
	}

	infrastructure := infra.Infra{
		DB: db.NewDBMock(t),
	}
	uc := usecase.New(infrastructure, chain)

	return uc, infrastructure.DB
}

func TestRecvAlert(t *testing.T) {
	uc, inf := setupAlertTest(t)

	input := alertchain.Alert{
		Alert: ent.Alert{
			Title:    "five",
			Detector: "blue",
		},
	}
	alert, err := uc.RecvAlert(&input)
	require.NoError(t, err)
	require.NotNil(t, alert)

	got, err := inf.GetAlert(alert.ID)
	require.NoError(t, err)
	assert.Equal(t, alert.Title, got.Title)
}
