package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Command holds the schema definition for the Command entity.
type Command struct {
	ent.Schema
}

// Annotations of the User.
func (Command) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "commands"},
	}
}

// Fields of the Command.
func (Command) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("name").Unique().NotEmpty(),
		field.Strings("content"),
		//field.Time("created_at").Default(time.Now),
		//field.Time("updated_at").Default(time.Now).
		//	UpdateDefault(time.Now),
	}
}

// Edges of the Command.
func (Command) Edges() []ent.Edge {
	return nil
}

// Mixin xxxx
func (Command) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
