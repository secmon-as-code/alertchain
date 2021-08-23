package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Finding holds the schema definition for the Finding entity.
type Finding struct {
	ent.Schema
}

// Fields of the Finding.
func (Finding) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("timestamp"),
		field.String("source"),
		field.String("name"),
		field.String("value"),
	}
}

// Edges of the Finding.
func (Finding) Edges() []ent.Edge {
	return nil
}
