package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Permission holds the schema definition for the Permission entity.
type Permission struct {
	ent.Schema
}

// Mixin MiXin Mixin User
func (Permission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Permission.
func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").Unique(),
		field.String("description"),
	}
}

// Edges of the Permission.
func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("roles", Role.Type).Ref("permissions"),
		// Add more edges here
		edge.To("users", User.Type),
	}
}
