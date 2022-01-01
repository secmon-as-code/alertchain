package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

// ExecLog holds the schema definition for the ExecLog entity.
type ExecLog struct {
	ent.Schema
}

// Fields of the ExecLog.
func (ExecLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("timestamp").Immutable(),
		field.String("log").Optional(),
		field.String("errmsg").Optional(),
		field.Strings("err_values").Optional(),
		field.Strings("stack_trace").Optional(),
		field.String("status").GoType(types.ExecStatus("")).Immutable(),
	}
}

// Edges of the ExecLog.
func (ExecLog) Edges() []ent.Edge {
	return nil
}
