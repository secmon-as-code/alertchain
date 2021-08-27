package main

import (
	"context"
	"fmt"
	"time"

	"github.com/m-mizutani/alertchain"
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/types"
)

type taskExample struct{}

func (x *taskExample) Name() string        { return "taskExample" }
func (x *taskExample) Description() string { return "Example of task" }
func (x *taskExample) Optionable(alert *alertchain.Alert) bool {
	return false
}

func (x *taskExample) Execute(ctx context.Context, alert *alertchain.Alert) error {
	w := alertchain.LogOutput(ctx)

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

	fmt.Fprintf(w, "done")
	return nil
}

func Chain() *alertchain.Chain {
	return &alertchain.Chain{
		Stages: []*alertchain.Stage{
			{
				Tasks: []alertchain.Task{
					&taskExample{},
				},
			},
		},
	}
}
