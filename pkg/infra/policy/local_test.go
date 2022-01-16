package policy_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalClient(t *testing.T) {
	t.Run("import recursive if specifing directory", func(t *testing.T) {
		client, err := policy.NewLocal("./testdata")
		require.NoError(t, err)
		in := map[string]string{
			"color":  "blue",
			"number": "five",
		}
		out := map[string]map[string]interface{}{}
		require.NoError(t, client.Eval(context.Background(), in, &out))
		assert.Equal(t, true, out["color"]["allow"])
		assert.Equal(t, true, out["number"]["allow"])
	})

	t.Run("import a file if specifing file path", func(t *testing.T) {
		client, err := policy.NewLocal("./testdata/policy.rego")
		require.NoError(t, err)
		in := map[string]string{
			"color":  "blue",
			"number": "five",
		}
		out := map[string]map[string]interface{}{}
		require.NoError(t, client.Eval(context.Background(), in, &out))
		assert.Equal(t, true, out["color"]["allow"])
		assert.Equal(t, nil, out["number"]["allow"])
	})

	t.Run("fail by specifying invalid path", func(t *testing.T) {
		_, err := policy.NewLocal("./testdata/not_found.rego")
		require.Error(t, err)
	})
}
