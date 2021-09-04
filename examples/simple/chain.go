package main

import (
	"github.com/m-mizutani/alertchain"
)

func Chain() (*alertchain.Chain, error) {
	return &alertchain.Chain{
		Jobs: []*alertchain.Job{
			{
				Tasks: []alertchain.Task{&Evaluator{}},
			},
			{
				Tasks: []alertchain.Task{&CreateTicket{}},
			},
		},
	}, nil
}
