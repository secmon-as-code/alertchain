package policy

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type Mock struct {
	t        *testing.T
	EvalMock func(in interface{}) (interface{}, error)
}

func NewMock(t *testing.T, eval func(in interface{}) (interface{}, error)) *Mock {
	return &Mock{
		t:        t,
		EvalMock: eval,
	}
}

func (x *Mock) Eval(ctx context.Context, in interface{}, out interface{}) error {
	result, err := x.EvalMock(in)
	if err != nil {
		return err
	}

	raw, err := json.Marshal(result)
	require.NoError(x.t, err)
	require.NoError(x.t, json.Unmarshal(raw, out))
	return nil
}
