package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/privacy"
)

type Tenant struct{ ent.Schema }

func (Tenant) Mixin() []ent.Mixin { return []ent.Mixin{IDMixin{}, TimeMixin{}} }

func (Tenant) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty(),
		field.String("description").Optional(),
		field.Enum("status").Values("active", "inactive", "suspended").Default("active"),
	}
}

func (Tenant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("environments", Environment.Type),
		edge.To("resources", Resource.Type),
		edge.To("access_policies", AccessPolicy.Type),
	}
}

func (Tenant) Policy() ent.Policy {
	return privacy.Policy{
		Query: privacy.QueryPolicy{
			AllowIfSuperAdmin(),
			privacy.AlwaysAllowRule(),
		},
		Mutation: privacy.MutationPolicy{
			AllowIfSuperAdmin(),
			privacy.AlwaysDenyRule(),
		},
	}
}
