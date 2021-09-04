package main

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/types"
	"github.com/pkg/errors"
)

type Evaluator struct{}

func (x *Evaluator) Name() string { return "Evaluation" }

func (x *Evaluator) Execute(ctx *types.Context, alert *alertchain.Alert) error {
	if alert.Title == "Suspicious Login" {
		attrs := alert.Attributes.FindByKey("srcAddr").FindByType(types.AttrIPAddr)
		if len(attrs) != 1 {
			return nil // Attribute not found
		}

		addr := net.ParseIP(attrs[0].Value)
		_, internal, _ := net.ParseCIDR("10.1.0.0/16")
		if internal.Contains(addr) {
			alert.UpdateSeverity(types.SevSafe)
		}
	}

	return nil
}

type CreateTicket struct{}

func (x *CreateTicket) Name() string { return "This is my first task" }

func (x *CreateTicket) Execute(ctx *types.Context, alert *alertchain.Alert) error {
	if alert.Severity == types.SevSafe {
		return nil // nothing to do if alert is "safe"
	}

	msg := struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}{
		Title: alert.Title,
		Body:  alert.Description,
	}

	url := "https://your-ticket-system.example.com/ticket"
	body, _ := json.Marshal(&msg)
	if _, err := http.Post(url, "application/json", bytes.NewReader(body)); err != nil {
		return errors.Wrap(err, "failed to create ticket")
	}

	return nil
}
