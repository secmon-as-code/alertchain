package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/alertchain/types"
)

// ActionLog holds the schema definition for the ActionLog entity.
type ActionLog struct {
	ent.Schema
}

// Fields of the ActionLog.
func (ActionLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Immutable(),
		field.Int64("started_at").Immutable(),
		field.Int64("exited_at").Optional(),
		field.String("log").Optional(),
		field.String("errmsg").Optional(),
		field.Strings("err_values").Optional(),
		field.Strings("stack_trace").Optional(),
		field.String("status").GoType(types.TaskStatus("")).Default(string(types.TaskRunning)),
	}
}

// Edges of the ActionLog.
func (ActionLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("argument", Attribute.Type),
	}
}
