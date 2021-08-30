package alertchain

import (
	"github.com/m-mizutani/alertchain/pkg/infra/ent"
	"github.com/m-mizutani/alertchain/pkg/usecase"
	"github.com/m-mizutani/alertchain/types"
)

type Alert struct {
	ent.Alert
	Attributes Attributes       `json:"attributes"`
	TaskLogs   []*ent.TaskLog   `json:"task_logs"`
	References []*ent.Reference `json:"references"`

	id types.AlertID // Immutable AlertID copied from ent.Alert.ID

	usecase.ChangeRequest `json:"-"`

	// To remove "edges" in JSON. DO NOT USE as data field
	EdgesOverride interface{} `json:"edges,omitempty"`
}

func NewAlert(alert *ent.Alert) *Alert {
	if alert == nil {
		return nil
	}

	newAlert := &Alert{
		Alert: *alert,
		id:    alert.ID,

		TaskLogs:   alert.Edges.TaskLogs,
		References: alert.Edges.References,
	}

	for _, attr := range alert.Edges.Attributes {
		newAlert.pushAttribute(attr)
	}

	return newAlert
}
