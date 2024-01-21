package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/m-mizutani/alertchain/pkg/controller/cli"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/gt"
)

func TestServe(t *testing.T) {
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

func TestPlay(t *testing.T) {
	ctx := context.Background()
	args := []string{
		"alertchain",
		"-l", "debug",
		"play",
		"-d", "examples/test/policy",
		"-s", "examples/test/scenarios",
		"-o", "examples/test/output",
	}
	gt.NoError(t, cli.New().Run(ctx, args))

	gt.F(t, "examples/test/output/scenario1/data.json").Reader(func(t testing.TB, r io.Reader) {
		var data model.ScenarioLog
		gt.NoError(t, json.NewDecoder(r).Decode(&data))
		gt.Equal(t, data.ID, "scenario1")
		gt.Equal(t, data.Title, "Test 1")
		gt.A(t, data.Results).Length(1).
			At(0, func(t testing.TB, v *model.PlayLog) {
				gt.Equal(t, v.Alert.Title, "Trojan:EC2/DropPoint!DNS")

				gt.A(t, v.Actions).Length(2).
					At(0, func(t testing.TB, v *model.ActionLog) {
						gt.Equal(t, v.Seq, 0)
						gt.A(t, v.Run).Length(1).
							At(0, func(t testing.TB, v model.Action) {
								gt.Equal(t, v.Uses, "chatgpt.query")
							})
					}).
					At(1, func(t testing.TB, v *model.ActionLog) {
						gt.Equal(t, v.Seq, 1)
						gt.A(t, v.Run).Length(1).
							At(0, func(t testing.TB, v model.Action) {
								gt.Equal(t, v.Uses, "slack.post")
							})
					})
			})
	})
}
