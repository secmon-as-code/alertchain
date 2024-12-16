package otx_test

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/alertchain/pkg/action/otx"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
)

func TestIndicator(t *testing.T) {
	testCases := []struct {
		title      string
		httpClient *http.Client
		args       model.ActionArgs
		hasErr     bool
		errType    error
	}{
		{
			title: "normal",
			httpClient: mockHTTPClient(func(req *http.Request) *http.Response {
				respJSON := `{"key": "value"}`
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader([]byte(respJSON))),
					Header:     make(http.Header),
				}
			}),
			args: model.ActionArgs{
				"secret_api_key": "dummy",
				"type":           "domain",
				"indicator":      "example.com",
				"section":        "general",
			},
			hasErr: false,
		},
		{
			title: "invalid api key",
			args: model.ActionArgs{
				"type":      "domain",
				"indicator": "example.com",
				"section":   "general",
			},
			hasErr:  true,
			errType: types.ErrActionInvalidArgument,
		},
		{
			title: "invalid type",
			args: model.ActionArgs{
				"secret_api_key": "dummy",
				"type":           "invalid",
				"indicator":      "example.com",
				"section":        "general",
			},
			hasErr:  true,
			errType: types.ErrActionInvalidArgument,
		},
		{
			title: "invalid section",
			args: model.ActionArgs{
				"secret_api_key": "dummy",
				"type":           "domain",
				"indicator":      "example.com",
				"section":        "invalid",
			},
			hasErr:  true,
			errType: types.ErrActionInvalidArgument,
		},
		{
			title: "http request error",
			httpClient: mockHTTPClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 500,
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
					Header:     make(http.Header),
				}
			}),
			args: model.ActionArgs{
				"secret_api_key": "dummy",
				"type":           "domain",
				"indicator":      "example.com",
				"section":        "general",
			},
			hasErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.title, func(t *testing.T) {
			if tt.httpClient != nil {
				otx.ReplaceHTTPClient(tt.httpClient)
				t.Cleanup(func() {
					otx.ReplaceHTTPClient(http.DefaultClient)
				})
			}

			ctx := model.NewContext()
			result, err := otx.Indicator(ctx, model.Alert{}, tt.args)

			if tt.hasErr {
				gt.Error(t, err).Must()
				if tt.errType != nil {
					gt.Error(t, err).Is(tt.errType)
				}
			} else {
				gt.NoError(t, err).Must()
				gt.V(t, result).NotNil()
			}
		})
	}
}

func mockHTTPClient(doFunc func(*http.Request) *http.Response) *http.Client {
	return &http.Client{
		Transport: mockRoundTripper(doFunc),
	}
}

type mockRoundTripper func(*http.Request) *http.Response

func (mrt mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return mrt(req), nil
}

func TestIndicatorRun(t *testing.T) {
	apiKey := os.Getenv("TEST_OTX_API_KEY")
	if apiKey == "" {
		t.Skip("TEST_OTX_API_KEY is not set")
	}

	ctx := model.NewContext()
	args := model.ActionArgs{
		"secret_api_key": apiKey,
		"type":           "IPv4",
		"indicator":      "87.236.176.4",
		"section":        "general",
	}

	result, err := otx.Indicator(ctx, model.Alert{}, args)
	gt.NoError(t, err).Must()
	gt.V(t, result).NotNil()
}
