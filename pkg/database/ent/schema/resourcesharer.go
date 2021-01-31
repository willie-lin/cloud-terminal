package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/dialect/entsql"
	"github.com/facebook/ent/schema"
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
	return nil
}

// Edges of the ResourceSharer.
func (ResourceSharer) Edges() []ent.Edge {
	return nil
}
