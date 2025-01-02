package server_test

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/alertchain/pkg/chain"
	"github.com/secmon-lab/alertchain/pkg/controller/graphql"
	"github.com/secmon-lab/alertchain/pkg/controller/server"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
	"github.com/secmon-lab/alertchain/pkg/domain/types"
	"github.com/secmon-lab/alertchain/pkg/infra/memory"
	"github.com/secmon-lab/alertchain/pkg/infra/policy"
	"github.com/secmon-lab/alertchain/pkg/service"
)

//go:embed testdata/scc.json
var sccData []byte

func TestSCC(t *testing.T) {
	var called int
	srv := server.New(func(ctx context.Context, schema types.Schema, data any) ([]*model.Alert, error) {
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
	srv := server.New(func(ctx context.Context, schema types.Schema, data any) ([]*model.Alert, error) {
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

func sendGraphQLRequest(t *testing.T, srv *server.Server, query string, out any) {
	body := gt.R1(json.Marshal(map[string]string{
		"query": query,
	})).NoError(t)
	req := httptest.NewRequest("POST", "/graphql", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	gt.N(t, w.Result().StatusCode).Equal(http.StatusOK)
	t.Log("body", w.Body.String())
	gt.NoError(t, json.Unmarshal(w.Body.Bytes(), out))
}

func TestGraphQL(t *testing.T) {
	dbClient := memory.New()
	chain := gt.R1(chain.New(
		chain.WithPolicyAlert(gt.R1(policy.New(
			policy.WithPackage("alert"),
			policy.WithPolicyData("alert.rego", alertRego),
		)).NoError(t)),
		chain.WithPolicyAction(gt.R1(policy.New(
			policy.WithPackage("action"),
			policy.WithPolicyData("action.rego", actionRego),
		)).NoError(t)),
		chain.WithDatabase(dbClient),
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
		}`

		var output struct {
			Data struct {
				Workflows []*model.WorkflowRecord `json:"workflows"`
			} `json:"data"`
		}
		sendGraphQLRequest(t, srv, q, &output)
		gt.N(t, len(output.Data.Workflows)).Equal(1)
		gt.S(t, string(output.Data.Workflows[0].Alert.ID)).Equal(alertID)
	})

	t.Run("query workflow via GraphQL for attrs", func(t *testing.T) {
		q := `query my_query {
			workflows(limit: 1) {
				id
				alert {
					id
					initAttrs {
						key
						value
					}
					lastAttrs {
						key
						value
					}
				}
			}
		}`

		var output struct {
			Data struct {
				Workflows []*model.WorkflowRecord `json:"workflows"`
			} `json:"data"`
		}
		sendGraphQLRequest(t, srv, q, &output)
		gt.N(t, len(output.Data.Workflows)).Equal(1)
		gt.S(t, string(output.Data.Workflows[0].Alert.ID)).Equal(alertID)

		gt.A(t, output.Data.Workflows[0].Alert.InitAttrs).Length(1).At(0, func(t testing.TB, v *model.AttributeRecord) {
			gt.V(t, v.Key).Equal("test_attr")
			gt.V(t, v.Value).Equal("test_value")
		})
		gt.A(t, output.Data.Workflows[0].Alert.LastAttrs).Length(2).
			At(0, func(t testing.TB, v *model.AttributeRecord) {
				gt.V(t, v.Key).Equal("test_attr")
				gt.V(t, v.Value).Equal("test_value")
			}).
			At(1, func(t testing.TB, v *model.AttributeRecord) {
				gt.V(t, v.Key).Equal("added_attr")
				gt.V(t, v.Value).Equal("swirls")
			})
	})
}
