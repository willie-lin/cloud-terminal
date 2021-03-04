package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Property holds the schema definition for the Property entity.
type Property struct {
	ent.Schema
}

// Annotations of the User.
func (Property) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "properties"},
	}
}

// Fields of the Property.
func (Property) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("value"),
	}
}

// Edges of the Property.
func (Property) Edges() []ent.Edge {
	return nil
}
