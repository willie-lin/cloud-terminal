package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/privacy"
)

type AccessPolicy struct{ ent.Schema }

func (AccessPolicy) Mixin() []ent.Mixin { return []ent.Mixin{IDMixin{}, TimeMixin{}} }

type PolicyStatement struct {
	Effect    string      `json:"Effect"`
	Actions   []string    `json:"Action"`
	Resources []string    `json:"Resource"`
	Condition interface{} `json:"Condition,omitempty"`
}

func (AccessPolicy) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Comment("策略名称"),
		field.String("description").Optional(),
		field.String("version").Default("v1").Comment("策略版本"),
		field.JSON("statements", []PolicyStatement{}).Default([]PolicyStatement{}).Comment("策略语句数组"),
		field.Bool("immutable").Default(false),
		field.Int("priority").Default(0).Comment("优先级，数值越小优先级越高"),
		field.Time("effective_date").Optional().Nillable().Comment("生效时间"),
		field.Time("expiry_date").Optional().Nillable().Comment("过期时间"),
	}
}

func (AccessPolicy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("groups", Group.Type).Ref("access_policies"),
		edge.From("users", User.Type).Ref("access_policies"),
		edge.From("resources", Resource.Type).Ref("policies"),
		edge.From("roles", Role.Type).Ref("access_policies"),
	}
}

func (AccessPolicy) Policy() ent.Policy {
	return privacy.Policy{
		Query: privacy.QueryPolicy{
			privacy.AlwaysAllowRule(),
		},
		Mutation: privacy.MutationPolicy{
			DenyMutationUnlessAdmin(),
		},
	}
}
