package schema

import (
	"entgo.io/ent/privacy"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Role struct{ ent.Schema }

func (Role) Mixin() []ent.Mixin { return []ent.Mixin{IDMixin{}, TimeMixin{}} }

func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty(),
		field.String("description").Optional(),
		field.Bool("is_disabled").Default(false),
		field.Bool("is_default").Default(false),
		field.JSON("trust_policy", map[string]interface{}{}).Optional().Comment("信任策略：控制谁能 Assume 此角色"),
		field.Time("effective_date").Optional().Nillable().Comment("角色生效时间"),
		field.Time("expiry_date").Optional().Nillable().Comment("角色过期时间"),
	}
}

func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).Ref("roles").Comment("可以 Assume 该角色的用户"),
		edge.To("access_policies", AccessPolicy.Type).Comment("附加到此角色的权限策略"),
		edge.To("child_roles", Role.Type).From("parent_role"),
		edge.To("permissions_boundary", AccessPolicy.Type).Unique().Comment("权限边界"),
	}
}


func (Role) Policy() ent.Policy {
	return privacy.Policy{
		Query: privacy.QueryPolicy{
			privacy.AlwaysAllowRule(),
		},
		Mutation: privacy.MutationPolicy{
			DenyMutationUnlessAdmin(),
		},
	}
}

func (Role) Indexes() []ent.Index { return []ent.Index{index.Fields("name").Unique()} }
