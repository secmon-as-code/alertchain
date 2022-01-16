package policy_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/stretchr/testify/require"
)

type httpMock struct {
	do func(req *http.Request) (*http.Response, error)
}

func (x *httpMock) Do(req *http.Request) (*http.Response, error) {
	return x.do(req)
}

func TestRemote(t *testing.T) {
	in := struct{ Number string }{Number: "five"}
	var out1, out2 struct{ Color string }

	var resp struct {
		Result interface{} `json:"result"`
	}

	ctx := context.Background()
	var called int

	remote, err := policy.NewRemoteWithHTTPClient("http://example.com",
		&httpMock{
			do: func(req *http.Request) (*http.Response, error) {
				called++
				out2.Color = "blue"
				resp.Result = out2

				var input map[string]map[string]interface{}
				require.NoError(t, json.NewDecoder(req.Body).Decode(&input))
				assert.Equal(t, "five", input["input"]["Number"])

				raw, err := json.Marshal(resp)
				require.NoError(t, err)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(raw)),
				}, nil
			},
		},
	)
	require.NoError(t, err)

	require.NoError(t, remote.Eval(ctx, in, &out1))
	assert.Equal(t, 1, called)
	assert.Equal(t, "blue", out1.Color)
}
