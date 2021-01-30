package schema

import "github.com/facebook/ent"

// Property holds the schema definition for the Property entity.
type Property struct {
	ent.Schema
}

// Fields of the Property.
func (Property) Fields() []ent.Field {
	return nil
}

// Edges of the Property.
func (Property) Edges() []ent.Edge {
	return nil
}
