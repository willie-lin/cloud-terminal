package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/privacy"
)

type Session struct{ ent.Schema }

func (Session) Mixin() []ent.Mixin { return []ent.Mixin{IDMixin{}, TimeMixin{}} }

func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String("session_id").Unique().Immutable(),
		field.String("principal_urn"),
		field.String("resource_urn").Optional(),
		field.String("environment_urn").Optional(),
		field.Enum("mode").Values("container", "proxy", "unknown").Default("unknown"),
		field.Enum("status").Values("active", "closed", "timeout", "error").Default("active"),
		field.Time("started_at"),
		field.Time("ended_at").Optional().Nillable(),
		field.String("remote_address").Optional(),
	}
}

func (Session) Policy() ent.Policy {
	return privacy.Policy{
		Query: privacy.QueryPolicy{
			privacy.AlwaysAllowRule(),
		},
		Mutation: privacy.MutationPolicy{
			DenyMutationUnlessAdmin(),
		},
	}
}
