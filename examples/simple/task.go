package main

import (
	"net"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type Evaluator struct{}

func (x *Evaluator) Name() string { return "Evaluator" }
func (x *Evaluator) Execute(ctx *types.Context, alert *alertchain.Alert) error {
	if alert.Title == "Suspicious Login" {
		evalSuspiciousLogin(alert)
	}

	return nil
}

func evalSuspiciousLogin(alert *alertchain.Alert) {
	attrs := alert.Attributes.FindByKey("srcAddr").FindByType(types.AttrIPAddr)
	if len(attrs) != 1 {
		return // Attribute not found
	}

	addr := net.ParseIP(attrs[0].Value)
	_, internal, _ := net.ParseCIDR("10.1.0.0/16")
	if internal.Contains(addr) {
		alert.UpdateSeverity(types.SevSafe)
	}
}

type CreateTicket struct{}

func (x *CreateTicket) Name() string { return "Create a ticket" }
func (x *CreateTicket) Execute(ctx *types.Context, alert *alertchain.Alert) error {
	alert.AddReference(&ent.Reference{
		Source: "ticket",
		Title:  "Link to ticket",
		URL:    "https://github.com/m-mizutani/alertchain/issues",
	})

	return nil
}