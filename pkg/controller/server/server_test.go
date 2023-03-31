package server_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/m-mizutani/alertchain/pkg/controller/server"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
	"github.com/m-mizutani/gt"
)

//go:embed testdata/scc.json
var sccData []byte

func TestSCC(t *testing.T) {
	var called int
	srv := server.New(func(ctx *types.Context, label string, data any) error {
		called++
		gt.V(t, label).Equal("scc")
		alert := gt.Cast[map[string]any](t, data)
		name := gt.Cast[string](t, alert["notificationConfigName"])
		gt.V(t, name).Equal("organizations/000000123456/notificationConfigs/pubsub_notification")
		return nil
	})

	req := httptest.NewRequest("POST", "/alert/raw/scc", bytes.NewReader(sccData))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	gt.N(t, w.Result().StatusCode).Equal(http.StatusOK)
	gt.N(t, called).Equal(1)
}

func TestPubSub(t *testing.T) {
	var called int
	srv := server.New(func(ctx *types.Context, label string, data any) error {
		called++
		gt.V(t, label).Equal("scc")
		alert := gt.Cast[map[string]any](t, data)
		name := gt.Cast[string](t, alert["color"])
		gt.V(t, name).Equal("blue")
		return nil
	})

	msg := pubsub.Message{
		ID:   "xxx",
		Data: []byte(`{"color":"blue"}`),
	}
	body := gt.R1(json.Marshal(msg)).NoError(t)

	req := httptest.NewRequest("POST", "/alert/pubsub/scc", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	gt.N(t, w.Result().StatusCode).Equal(http.StatusOK)
	gt.N(t, called).Equal(1)
}