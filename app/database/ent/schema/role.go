package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Mixin MiXin Mixin User
func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").Unique(),
		field.String("description"),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		//edge.From("users", User.Type).Ref("roles").Unique(),
		//edge.To("permissions", Permission.Type),
		edge.From("tenant", Tenant.Type).Ref("roles").Unique(),
		edge.From("users", User.Type).Ref("roles"),
		edge.To("permissions", Permission.Type),
	}
}
