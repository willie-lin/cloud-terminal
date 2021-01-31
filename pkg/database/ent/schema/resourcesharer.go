package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/dialect/entsql"
	"github.com/facebook/ent/schema"
	"github.com/facebook/ent/schema/field"
)

// ResourceSharer holds the schema definition for the ResourceSharer entity.
type ResourceSharer struct {
	ent.Schema
}

// Annotations of the User.
func (ResourceSharer) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "resourceSharers"},
	}
}

// Fields of the ResourceSharer.
func (ResourceSharer) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("resource_id"),
		field.String("resource_type"),
		field.String("user_id"),
		field.String("userGroup_id"),
	}
}

// Edges of the ResourceSharer.
func (ResourceSharer) Edges() []ent.Edge {
	return nil
}
