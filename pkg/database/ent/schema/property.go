package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/dialect/entsql"
	"github.com/facebook/ent/schema"
	"github.com/facebook/ent/schema/field"
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
