package main_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/controller/cli"
	"github.com/m-mizutani/gt"
)

func TestE2E(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		args := []string{
			"alertchain",
			"serve",
			"-p",
			"--addr", "127.0.0.1:8080",
			"-d", "examples/e2e",
		}
		gt.NoError(t, cli.New().Run(ctx, args))
	}()

	var called int
	callbackHandler := func(w http.ResponseWriter, r *http.Request) {
		t.Log("called!")
		called++
		w.WriteHeader(http.StatusOK)
		gt.R1(w.Write([]byte("OK"))).NoError(t)
	}
	go func() {
		gt.NoError(t, http.ListenAndServe("127.0.0.1:9876", http.HandlerFunc(callbackHandler)))
	}()

	send := func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"color":"blue"}`))
		req := gt.R1(http.NewRequest("POST", "http://127.0.0.1:8080/alert/raw/my_alert", body)).NoError(t)
		resp := gt.R1(http.DefaultClient.Do(req)).NoError(t)
		gt.N(t, resp.StatusCode).Equal(200)
	}
	sendIgnoredAlert := func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"color":"red"}`))
		req := gt.R1(http.NewRequest("POST", "http://127.0.0.1:8080/alert/raw/my_alert", body)).NoError(t)
		resp := gt.R1(http.DefaultClient.Do(req)).NoError(t)
		gt.N(t, resp.StatusCode).Equal(200)
	}

	time.Sleep(time.Second)

	send(t) // 1
	gt.N(t, called).Equal(0)
	send(t) // 2
	gt.N(t, called).Equal(0)
	send(t) // 3
	gt.N(t, called).Equal(1)
	sendIgnoredAlert(t) // ignored
	gt.N(t, called).Equal(1)
	send(t) // 4
	gt.N(t, called).Equal(2)
}
