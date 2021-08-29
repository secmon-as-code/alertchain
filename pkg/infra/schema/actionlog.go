package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ActionLog holds the schema definition for the ActionLog entity.
type ActionLog struct {
	ent.Schema
}

// Fields of the ActionLog.
func (ActionLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Immutable(),
	}
}

// Edges of the ActionLog.
func (ActionLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("argument", Attribute.Type),
		edge.To("exec_logs", ExecLog.Type),
	}
}
