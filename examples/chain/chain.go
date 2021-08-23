package main

import (
	"context"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/types"
)

type myEvaluator struct{}

func (x *myEvaluator) Name() string        { return "myEvaluator" }
func (x *myEvaluator) Description() string { return "Eval alert" }
func (x *myEvaluator) Optionable(alert *alertchain.Alert) bool {
	return false
}

func (x *myEvaluator) Execute(ctx context.Context, alert *alertchain.Alert) error {
	if alert.Title == "Something wrong" {
		alert.Severity = types.SevAffected
	}
	if err := alert.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func Chain() *alertchain.Chain {
	chain := &alertchain.Chain{}
	chain.NewStage().AddTask(&myEvaluator{})
	return chain
}
