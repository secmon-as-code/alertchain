package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

// Attribute holds the schema definition for the Attribute entity.
type Attribute struct {
	ent.Schema
}

// Fields of the Attribute.
func (Attribute) Fields() []ent.Field {
	return []ent.Field{
		field.String("key"),
		field.String("value"),
		field.String("type").GoType(types.AttrType("")),
		field.Strings("context"),
	}
}

// Edges of the Attribute.
func (Attribute) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To(edgeAttrsToAnn, Annotation.Type),
		edge.From(edgeAttrToAlert, Alert.Type).Ref(edgeAlertToAttrs).Unique(),
	}
}
