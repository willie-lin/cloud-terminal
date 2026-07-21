package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/privacy"
)

type Platform struct{ ent.Schema }

func (Platform) Mixin() []ent.Mixin { return []ent.Mixin{IDMixin{}, TimeMixin{}} }

func (Platform) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty(),
		field.String("description").Optional(),
		field.String("region").Optional(),
		field.String("version").Optional(),
		field.Enum("status").Values("active", "maintenance", "disabled").Default("active"),
		field.JSON("config", map[string]interface{}{}).Optional(),
	}
}

func (Platform) Policy() ent.Policy {
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
