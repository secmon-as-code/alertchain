package schema

import "entgo.io/ent"

// ExecLog holds the schema definition for the ExecLog entity.
type ExecLog struct {
	ent.Schema
}

// Fields of the ExecLog.
func (ExecLog) Fields() []ent.Field {
	return nil
}

// Edges of the ExecLog.
func (ExecLog) Edges() []ent.Edge {
	return nil
}
