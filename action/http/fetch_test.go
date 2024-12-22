package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/m-mizutani/gt"
	httpaction "github.com/secmon-lab/alertchain/action/http"
	"github.com/secmon-lab/alertchain/pkg/domain/model"
)

func TestFetch(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gt.Value(t, r.Method).Equal("GET")
			gt.Value(t, r.URL.Path).Equal("/test")
			w.Header().Add("Content-Type", "application/json")
			gt.R1(w.Write([]byte(`{"message":"Hello"}`))).NoError(t)
		}))
		defer ts.Close()

		ctx := context.Background()

		args := model.ActionArgs{
			"method": "GET",
			"url":    ts.URL + "/test",
		}

		result := gt.R1(httpaction.Fetch(ctx, model.Alert{}, args)).NoError(t)

		res := gt.Cast[map[string]interface{}](t, result)

		gt.Value(t, res["message"].(string)).Equal("Hello")
	})

	t.Run("Invalid argument", func(t *testing.T) {
		ctx := context.Background()

		args := model.ActionArgs{}

		gt.R1(httpaction.Fetch(ctx, model.Alert{}, args)).Error(t)
	})

	t.Run("Invalid JSON response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			gt.R1(w.Write([]byte(`{"message"}`))).NoError(t)
		}))
		defer ts.Close()

		ctx := context.Background()

		args := model.ActionArgs{
			"method": "GET",
			"url":    ts.URL + "/test",
		}

		gt.R1(httpaction.Fetch(ctx, model.Alert{}, args)).Error(t)
	})

	t.Run("HTTP request error", func(t *testing.T) {
		ctx := context.Background()

		args := model.ActionArgs{
			"method": "GET",
			"url":    "http://example.invalid",
		}

		gt.R1(httpaction.Fetch(ctx, model.Alert{}, args)).Error(t)
	})
}
