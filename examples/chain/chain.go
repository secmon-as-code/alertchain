package main

import (
	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/types"
)

type myEvaluator struct{}

func (x *myEvaluator) Name() string                              { return "myEvaluator" }
func (x *myEvaluator) Description() string                       { return "Eval alert" }
func (x *myEvaluator) IsExecutable(alert *alertchain.Alert) bool { return false }

func (x *myEvaluator) Execute(alert *alertchain.Alert) error {
	if alert.Title == "Something wrong" {
		alert.Severity = types.SevAffected
	}
	if err := alert.Commit(); err != nil {
		return err
	}
	return nil
}

func Chain() *alertchain.Chain {
	return &alertchain.Chain{
		Stages: []alertchain.Tasks{
			{&myEvaluator{}},
		},
	}
}
