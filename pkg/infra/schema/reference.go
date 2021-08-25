package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Reference holds the schema definition for the Reference entity.
type Reference struct {
	ent.Schema
}

// Fields of the Reference.
func (Reference) Fields() []ent.Field {
	return []ent.Field{
		field.String("source"),
		field.String("title"),
		field.String("url"),
		field.String("comment").Optional(),
	}
}

// Edges of the Reference.
func (Reference) Edges() []ent.Edge {
	return nil
}
