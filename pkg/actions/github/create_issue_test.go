package github

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
type mockTripper struct{}

func (x *mockTripper) RoundTrip(*http.Request) (*http.Response, error) {
	utils.Logger.Info("roundtrip-----------")
	issue := github.Issue{
		HTMLURL: github.String("https://github.com/m-mizutani/alertchain/issues/0"),
	}
	body, err := json.Marshal(issue)
	if err != nil {
		panic(err.Error())
	}

	return &http.Response{
		StatusCode: http.StatusCreated,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}, nil
}
*/

func TestCreateIssue(t *testing.T) {
	dummyKey := `-----BEGIN RSA PRIVATE KEY-----
MGICAQACEQDC5sfxmKFkC+qY6zUcMhy3AgMBAAECEArjM3YmjZV2D8/fBC9RdJEC
CQD8OxoMSVd5WQIJAMXQXkRIAXSPAghRYVxJJIy1mQIJALf+89/xViEzAghAKoVb
kz+vhw==
-----END RSA PRIVATE KEY-----`

	envKey := "GITHUB_PRIVATE_KEY_" + uuid.NewString()
	os.Setenv(envKey, dummyKey)
	newAction, err := NewCreateIssue(model.ActionConfig{
		"app_id":          123,
		"install_id":      234,
		"private_key_env": envKey,
		"owner":           "blue",
		"repo":            "five",
	})
	require.NoError(t, err)
	require.NotNil(t, newAction)
	action, ok := newAction.(*CreateIssue)
	require.True(t, ok)
	assert.Equal(t, int64(123), action.appID)
	assert.Equal(t, int64(234), action.installID)
	assert.Equal(t, []byte(dummyKey), action.privateKey)
	assert.Equal(t, "blue", action.owner)
	assert.Equal(t, "five", action.repo)

	/*
		// ToBeFixed: Skip for now
			alert := model.NewAlert(&model.Alert{
				Title:       "test",
				Description: "for testing",
				Detector:    "test-detector",
			})
				action.rt = &mockTripper{}

				req, err := action.Run(types.NewContext(utils.Logger), alert)
				require.NoError(t, err)
				require.Len(t, req.NewReferences, 1)
				assert.Equal(t, "https://github.com/m-mizutani/alertchain/issues/0", req.NewReferences[0].URL)
	*/
}

func TestCreateIssueFailure(t *testing.T) {
	envKey := "GITHUB_PRIVATE_KEY_" + uuid.NewString()
	os.Setenv(envKey, "my_secret")

	testCases := []struct {
		desc string
		cfg  model.ActionConfig
	}{
		{
			desc: "fail when app_id is not set",
			cfg: model.ActionConfig{
				"install_id":      234,
				"private_key_env": envKey,
				"owner":           "blue",
				"repo":            "five",
			},
		},
		{
			desc: "fail when app_id is not number",
			cfg: model.ActionConfig{
				"app_id":          true,
				"install_id":      234,
				"private_key_env": envKey,
				"owner":           "blue",
				"repo":            "five",
			},
		},
		{
			desc: "fail when install_id is not set",
			cfg: model.ActionConfig{
				"app_id":          123,
				"private_key_env": envKey,
				"owner":           "blue",
				"repo":            "five",
			},
		},
		{
			desc: "fail when install_id is not number",
			cfg: model.ActionConfig{
				"app_id":          123,
				"install_id":      true,
				"private_key_env": envKey,
				"owner":           "blue",
				"repo":            "five",
			},
		},
		{
			desc: "fail when private_key_env is not set",
			cfg: model.ActionConfig{
				"app_id":     123,
				"install_id": 234,
				"owner":      "blue",
				"repo":       "five",
			},
		},
		{
			desc: "fail when private_key_env var is not found",
			cfg: model.ActionConfig{
				"app_id":          123,
				"install_id":      234,
				"private_key_env": "???",
				"owner":           "blue",
				"repo":            "five",
			},
		},
		{
			desc: "fail when private_key_env is not string",
			cfg: model.ActionConfig{
				"app_id":          123,
				"install_id":      234,
				"private_key_env": 999,
				"owner":           "blue",
				"repo":            "five",
			},
		},

		{
			desc: "fail when owner is not set",
			cfg: model.ActionConfig{
				"app_id":          123,
				"install_id":      234,
				"private_key_env": envKey,
				"repo":            "five",
			},
		},
		{
			desc: "fail when owner is not string",
			cfg: model.ActionConfig{
				"app_id":          123,
				"install_id":      234,
				"private_key_env": envKey,
				"owner":           123,
				"repo":            "five",
			},
		},

		{
			desc: "fail when owner is not set",
			cfg: model.ActionConfig{
				"app_id":          123,
				"install_id":      234,
				"private_key_env": envKey,
				"owner":           "blue",
			},
		},
		{
			desc: "fail when owner is not string",
			cfg: model.ActionConfig{
				"app_id":          123,
				"install_id":      234,
				"private_key_env": envKey,
				"owner":           "blue",
				"repo":            666,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := NewCreateIssue(tC.cfg)
			assert.ErrorIs(t, err, types.ErrInvalidActionConfig)
		})
	}
}
