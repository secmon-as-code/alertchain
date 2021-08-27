package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/alertchain/types"
)

// TaskLog holds the schema definition for the TaskLog entity.
type TaskLog struct {
	ent.Schema
}

// Fields of the TaskLog.
func (TaskLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("task_name").Immutable(),
		field.Bool("optional").Default(false).Immutable(),
		field.Int64("stage").Immutable(),
		field.Int64("started_at").Immutable(),
		field.Int64("exited_at").Optional(),
		field.String("log").Optional(),
		field.String("errmsg").Optional(),
		field.Strings("err_values").Optional(),
		field.Strings("stack_trace").Optional(),
		field.String("status").GoType(types.TaskStatus("")).Default(string(types.TaskRunning)),
	}
}

// Edges of the TaskLog.
func (TaskLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("annotated", Annotation.Type),
	}
}
