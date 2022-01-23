package chain_test

import (
	"testing"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/infra"
	"github.com/m-mizutani/alertchain/pkg/infra/db"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/alertchain/pkg/service/chain"
	"github.com/stretchr/testify/require"
)

func newClients(t *testing.T, eval func(in interface{}) (interface{}, error)) *infra.Clients {
	return infra.New(db.NewDBMock(t), policy.NewMock(t, eval))
}

func TestChain(t *testing.T) {
	clients := newClients(t, func(in interface{}) (interface{}, error) {
		return struct{}{}, nil
	})
	c, err := chain.New(clients, []model.ActionDefinition{}, []model.JobDefinition{})
	require.NoError(t, err)
	require.NotNil(t, c)
}
