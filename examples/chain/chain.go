package main

import (
	"time"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/examples/tasks"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type TaskExample struct{}

func (x *TaskExample) Name() string { return "TaskExample" }

func (x *TaskExample) Execute(ctx *types.Context, alert *alertchain.Alert) error {
	// Update serverity
	alert.UpdateSeverity(types.SevUnclassified)

	// Annoate additional info to attributes
	for _, attr := range alert.Attributes {
		attr.Annotate(&alertchain.Annotation{
			Annotation: ent.Annotation{
				Timestamp: time.Now().UTC().Unix(),
				Source:    "example task",
				Name:      "Accessed from",
				Value:     "192.168.0.1",
			},
		})
	}

	// Add references
	alert.AddReference(&ent.Reference{
		Source:  "example task",
		Title:   "github issue",
		URL:     "https://github.com/m-mizutani/alertchain/issues",
		Comment: "test link",
	})

	return nil
}

func Chain() (*alertchain.Chain, error) {
	return &alertchain.Chain{
		Jobs: []*alertchain.Job{
			{
				Tasks: []alertchain.Task{
					&tasks.Evaluator{},
				},
			},
		},
	}, nil
}
