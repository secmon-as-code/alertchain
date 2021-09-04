package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Annotation holds the schema definition for the Annotation entity.
type Annotation struct {
	ent.Schema
}

// Fields of the Annotation.
func (Annotation) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("timestamp"),
		field.String("source"),
		field.String("name"),
		field.String("value"),
	}
}

// Edges of the Annotation.
func (Annotation) Edges() []ent.Edge {
	return nil
}
