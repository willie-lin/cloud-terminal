package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Test holds the schema definition for the Test entity.
type Test struct {
	ent.Schema
}

// Fields of the Test.
func (Test) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("name"),
		field.String("password"),
		field.String("email"),
	}
}

// Edges of the Test.
func (Test) Edges() []ent.Edge {
	return nil
}

// Mixin xxxx
func (Test) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
