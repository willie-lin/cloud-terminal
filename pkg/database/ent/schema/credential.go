package schema

import "github.com/facebook/ent"

// Credential holds the schema definition for the Credential entity.
type Credential struct {
	ent.Schema
}

// Fields of the Credential.
func (Credential) Fields() []ent.Field {
	return nil
}

// Edges of the Credential.
func (Credential) Edges() []ent.Edge {
	return nil
}
