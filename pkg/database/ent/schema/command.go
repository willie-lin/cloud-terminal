package schema

import "github.com/facebook/ent"

// Command holds the schema definition for the Command entity.
type Command struct {
	ent.Schema
}

// Fields of the Command.
func (Command) Fields() []ent.Field {
	return nil
}

// Edges of the Command.
func (Command) Edges() []ent.Edge {
	return nil
}
