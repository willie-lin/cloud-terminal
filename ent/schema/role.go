package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").Unique().NotEmpty(),
		field.String("description").Optional(),
		field.Bool("is_disabled").Default(false),
		field.Bool("is_default").Default(false),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("roles").Unique().Comment("角色所属的账户"),
		edge.From("users", User.Type).Ref("role").Comment("拥有该角色的用户"),
		edge.To("access_policies", AccessPolicy.Type).Comment("分配给角色的策略"),
		edge.To("child_roles", Role.Type).From("parent_role"),
	}
}

// Indexes of the Role.
func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}
