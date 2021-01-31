package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
	"time"
)

// Command holds the schema definition for the Command entity.
type Command struct {
	ent.Schema
}

// Fields of the Command.
func (Command) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().NotEmpty(),
		field.String("name"),
		field.Strings("content"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Command.
func (Command) Edges() []ent.Edge {
	return nil
}
