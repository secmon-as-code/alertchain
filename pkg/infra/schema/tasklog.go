package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TaskLog holds the schema definition for the TaskLog entity.
type TaskLog struct {
	ent.Schema
}

// Fields of the TaskLog.
func (TaskLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Immutable(),
		field.Int64("stage").Immutable(),
	}
}

// Edges of the TaskLog.
func (TaskLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("annotated", Annotation.Type),
		edge.To("exec_logs", ExecLog.Type),
	}
}
