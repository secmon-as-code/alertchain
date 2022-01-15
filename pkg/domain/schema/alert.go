package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

// Alert holds the schema definition for the Alert entity.
type Alert struct {
	ent.Schema
}

// Fields of the Alert.
func (Alert) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").GoType(types.AlertID("")).Immutable(),
		field.String("title").Optional(),
		field.String("description").Optional(),
		field.String("detector").Optional(),
		field.String("status").GoType(types.AlertStatus("")).Default(string(types.StatusNew)),
		field.String("severity").GoType(types.Severity("")).Optional(),
		field.Int64("detected_at").Optional(),

		field.Int64("created_at").Immutable(),
		field.Int64("closed_at").Optional(),
	}
}

// Edges of the Alert.
func (Alert) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To(edgeAlertToAttrs, Attribute.Type),
		edge.To(edgeAlertToRef, Reference.Type),
		edge.To(edgeAlertToJob, Job.Type),
	}
}
