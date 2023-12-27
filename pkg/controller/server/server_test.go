package server_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"testing"

	"github.com/m-mizutani/alertchain/pkg/chain"
	"github.com/m-mizutani/alertchain/pkg/chain/core"
	"github.com/m-mizutani/alertchain/pkg/controller/graphql"
	"github.com/m-mizutani/alertchain/pkg/controller/server"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/alertchain/pkg/infra/memory"
	"github.com/m-mizutani/alertchain/pkg/infra/policy"
	"github.com/m-mizutani/alertchain/pkg/service"
	"github.com/m-mizutani/gt"
)

//go:embed testdata/scc.json
var sccData []byte

func TestSCC(t *testing.T) {
	var called int
	srv := server.New(func(ctx *model.Context, schema types.Schema, data any) ([]*model.Alert, error) {
		called++
		gt.V(t, schema).Equal("scc")
		alert := gt.Cast[map[string]any](t, data)
		name := gt.Cast[string](t, alert["notificationConfigName"])
		gt.V(t, name).Equal("organizations/000000123456/notificationConfigs/pubsub_notification")
		return nil, nil
	})

	req := httptest.NewRequest("POST", "/alert/raw/scc", bytes.NewReader(sccData))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	gt.N(t, w.Result().StatusCode).Equal(http.StatusOK)
	gt.N(t, called).Equal(1)
}

func TestPubSub(t *testing.T) {
	var called int
	srv := server.New(func(ctx *model.Context, schema types.Schema, data any) ([]*model.Alert, error) {
		called++
		gt.V(t, schema).Equal("scc")
		alert := gt.Cast[map[string]any](t, data)
		name := gt.Cast[string](t, alert["color"])
		gt.V(t, name).Equal("blue")
		return nil, nil
	})

	req := model.PubSubRequest{
		Message: model.PubSubMessage{
			Data: []byte(`{"color":"blue"}`),
		},
	}

	body := gt.R1(json.Marshal(req)).NoError(t)

	httpReq := httptest.NewRequest("POST", "/alert/pubsub/scc", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, httpReq)
	gt.N(t, w.Result().StatusCode).Equal(http.StatusOK)
	gt.N(t, called).Equal(1)
}

//go:embed testdata/alert.rego
var alertRego string

//go:embed testdata/action.rego
var actionRego string

func TestGraphQL(t *testing.T) {
	dbClient := memory.New()
	chain := gt.R1(chain.New(
		core.WithPolicyAlert(gt.R1(policy.New(
			policy.WithPackage("alert"),
			policy.WithPolicyData("alert.rego", alertRego),
		)).NoError(t)),
		core.WithPolicyAction(gt.R1(policy.New(
			policy.WithPackage("action"),
			policy.WithPolicyData("action.rego", actionRego),
		)).NoError(t)),
		core.WithDatabase(dbClient),
	)).NoError(t)

	resolver := graphql.NewResolver(service.New(dbClient))
	srv := server.New(chain.HandleAlert, server.WithResolver(resolver))

	var alertID string
	t.Run("receive alert", func(t *testing.T) {
		var output struct {
			Alerts []*model.Alert `json:"alerts"`
		}
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, httptest.NewRequest("POST", "/alert/raw/test_service", strings.NewReader(`{"foo":"bar"}`)))
		gt.N(t, w.Result().StatusCode).Equal(http.StatusOK)
		gt.NoError(t, json.Unmarshal(w.Body.Bytes(), &output))
		gt.N(t, len(output.Alerts)).Equal(1)
		alertID = string(output.Alerts[0].ID)
	})
	print(alertID)

	t.Run("query workflow via GraphQL", func(t *testing.T) {
		q := `query my_query {
			workflows(limit: 1) {
			  id
			  alert {
				id
			  }
			}
		  }
		`
		body := gt.R1(json.Marshal(map[string]string{
			"query": q,
		})).NoError(t)
		req := httptest.NewRequest("POST", "/graphql", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		gt.N(t, w.Result().StatusCode).Equal(http.StatusOK)

		var output struct {
			Data struct {
				Workflows []*model.WorkflowRecord `json:"workflows"`
			} `json:"data"`
		}
		gt.NoError(t, json.Unmarshal(w.Body.Bytes(), &output))
		gt.N(t, len(output.Data.Workflows)).Equal(1)
		gt.S(t, string(output.Data.Workflows[0].Alert.ID)).Equal(alertID)
	})
}
