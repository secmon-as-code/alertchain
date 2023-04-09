package github_test

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/action/github"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/gt"
)

func TestIssueTemplate(t *testing.T) {
	var buf bytes.Buffer
	gt.NoError(t, github.ExecuteTemplate(&buf, model.Alert{
		AlertMetaData: model.AlertMetaData{
			Title:       "blue",
			Description: "orange",
			Params: []types.Parameter{
				{
					Key:   "magic",
					Value: "five",
				},
			},
		},
		Schema: "fire",
		Raw:    `{"foo": "bar"}`,
	}))

	s := buf.String()
	gt.B(t, strings.Contains(s, "orange")).True()
	gt.B(t, strings.Contains(s, "| magic | `five` |")).True()
	gt.B(t, strings.Contains(s, `{"foo": "bar"}`)).True()

	println(s)
}

func TestIssuer(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_GITHUB_ISSUER"); !ok {
		t.Skip("Skipping test because TEST_GITHUB_ISSUER is not set")
	}

	cfg := model.ActionConfigValues{
		"app_id":      gt.R1(strconv.Atoi(os.Getenv("TEST_GITHUB_APP_ID"))).NoError(t),
		"install_id":  gt.R1(strconv.Atoi(os.Getenv("TEST_GITHUB_INSTALL_ID"))).NoError(t),
		"private_key": os.Getenv("TEST_GITHUB_PRIVATE_KEY"),
		"owner":       os.Getenv("TEST_GITHUB_OWNER"),
		"repo":        os.Getenv("TEST_GITHUB_REPO"),
	}

	requiredVars := []string{"app_id", "install_id", "private_key", "owner", "repo"}
	for _, key := range requiredVars {
		gt.V(t, cfg[key]).NotEqual("")
	}

	factory := &github.IssuerFactory{}
	issuer := gt.R1(factory.New("test", cfg)).NoError(t)

	ctx := types.NewContext()
	alert := model.Alert{
		AlertMetaData: model.AlertMetaData{
			Title:       "blue",
			Description: "orange",
			Params: []types.Parameter{
				{
					Key:   "magic",
					Value: "five",
				},
			},
		},
		CreatedAt: time.Now(),
		Raw:       `{"foo": "bar"}`,
	}

	params := model.ActionParams{}
	gt.NoError(t, issuer.Run(ctx, alert, params))
}
