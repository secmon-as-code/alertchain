package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

func main() {
	alerts := genTestAlert()

	const (
		host = "127.0.0.1"
		port = 8080
	)

	client := &http.Client{}
	for _, alert := range alerts {
		url := fmt.Sprintf("http://%s:%d/api/v1/alert", host, port)
		raw, err := json.Marshal(alert)
		if err != nil {
			panic("json.Marshal: " + err.Error())
		}
		req, err := http.NewRequest("POST", url, bytes.NewReader(raw))
		if err != nil {
			panic("http.NewRequest: " + err.Error())
		}

		resp, err := client.Do(req)
		if err != nil {
			panic("http.client.Do: " + err.Error())
		}
		if resp.StatusCode != http.StatusCreated {
			log.Printf("Invalid status code: %d\n%v\n", resp.StatusCode, req)
		}
	}
}

func int64ptr(v int64) *int64 {
	return &v
}
func contexts(c ...types.AttrContext) []string {
	resp := make([]string, len(c))
	for i := range c {
		resp[i] = string(c[i])
	}
	return resp
}

func genTestAlert() []*alertchain.Alert {

	return []*alertchain.Alert{
		{
			Alert: ent.Alert{
				Title:       "Suspicious Login",
				Detector:    "Google Workspace",
				Description: "Account 'blue' accessed to Google Workspace from unusual geolocation",
				DetectedAt:  int64ptr(time.Now().UTC().Unix()),
			},
			Attributes: []*alertchain.Attribute{
				{
					Attribute: ent.Attribute{
						Key:     "client.addr",
						Value:   "203.0.113.1",
						Type:    types.AttrIPAddr,
						Context: contexts(types.CtxRemote, types.CtxClient),
					},
				},
				{
					Attribute: ent.Attribute{
						Key:     "actor.email",
						Value:   "mizutani@example.com",
						Type:    types.AttrEmail,
						Context: contexts(types.CtxRemote, types.CtxClient),
					},
				},
				{
					Attribute: ent.Attribute{
						Key:   "service",
						Value: "GMail",
						Type:  types.AttrNoType,
					},
				},
			},
		},
	}
}
