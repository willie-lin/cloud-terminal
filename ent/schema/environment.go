package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/privacy"
)

type Environment struct{ ent.Schema }

func (Environment) Mixin() []ent.Mixin { return []ent.Mixin{IDMixin{}, TimeMixin{}} }

func (Environment) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty(),
		field.String("description").Optional(),
		field.String("image").NotEmpty(),
		field.Int("port").Default(22),
		field.JSON("resource_limit", map[string]interface{}{}).Optional(),
		field.JSON("env_vars", map[string]interface{}{}).Optional(),
		field.JSON("volumes", []map[string]interface{}{}).Optional(),
		field.Enum("status").Values("active", "inactive").Default("active"),
	}
}

func (Environment) Edges() []ent.Edge {
	return []ent.Edge{edge.From("tenant", Tenant.Type).Ref("environments").Unique()}
}

func (Environment) Policy() ent.Policy {
	return privacy.Policy{
		Query: privacy.QueryPolicy{
			AllowIfAdmin(),
			privacy.AlwaysAllowRule(),
		},
		Mutation: privacy.MutationPolicy{
			DenyMutationUnlessAdmin(),
		},
	}
}
