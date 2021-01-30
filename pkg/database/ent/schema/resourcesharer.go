package schema

import "github.com/facebook/ent"

// ResourceSharer holds the schema definition for the ResourceSharer entity.
type ResourceSharer struct {
	ent.Schema
}

// Fields of the ResourceSharer.
func (ResourceSharer) Fields() []ent.Field {
	return nil
}

// Edges of the ResourceSharer.
func (ResourceSharer) Edges() []ent.Edge {
	return nil
}
