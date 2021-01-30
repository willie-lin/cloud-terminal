package schema

import "github.com/facebook/ent"

// Assets holds the schema definition for the Assets entity.
type Assets struct {
	ent.Schema
}

// Fields of the Assets.
func (Assets) Fields() []ent.Field {
	return nil
}

// Edges of the Assets.
func (Assets) Edges() []ent.Edge {
	return nil
}
