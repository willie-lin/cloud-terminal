package schema

import "github.com/facebook/ent"

// UserGroup holds the schema definition for the UserGroup entity.
type UserGroup struct {
	ent.Schema
}

// Fields of the UserGroup.
func (UserGroup) Fields() []ent.Field {
	return nil
}

// Edges of the UserGroup.
func (UserGroup) Edges() []ent.Edge {
	return nil
}
