package schema

import "entgo.io/ent"

// AssetGroup holds the schema definition for the AssetGroup entity.
type AssetGroup struct {
	ent.Schema
}

// Fields of the AssetGroup.
func (AssetGroup) Fields() []ent.Field {
	return nil
}

// Edges of the AssetGroup.
func (AssetGroup) Edges() []ent.Edge {
	return nil
}
