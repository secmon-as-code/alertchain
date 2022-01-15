package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/alertchain/pkg/domain/model"
)

// Job holds the schema definition for the Job entity.
type Job struct {
	ent.Schema
}

// Fields of the Job.
func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Immutable(),
		field.Int64("step").Immutable(),
		field.JSON("input", &model.Alert{}),
	}
}

// Edges of the Job.
func (Job) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From(edgeJobToAlert, Alert.Type).Ref(edgeAlertToJob).Unique(),
		edge.To(edgeJobToActionLog, ActionLog.Type),
	}
}
