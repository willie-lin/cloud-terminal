package schema

import (
	"entgo.io/ent/privacy"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Task struct{ ent.Schema }

func (Task) Mixin() []ent.Mixin { return []ent.Mixin{IDMixin{}, TimeMixin{}} }

func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.String("reason").NotEmpty().Comment("申请理由"),
		field.Int("duration_minutes").Positive().Comment("申请时长（分钟）"),
		field.Enum("status").Values("pending", "approved", "rejected", "expired").Default("pending"),
		field.Time("reviewed_at").Optional().Nillable().Comment("审批时间"),
		field.String("reviewer_comment").Optional().Comment("审批意见"),
		field.String("issued_token").Optional().Sensitive().Comment("STS token（审批通过后签发）"),
		field.Time("expires_at").Optional().Nillable().Comment("token 过期时间"),
	}
}

func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("requester", User.Type).Ref("tasks").Unique().Required().Comment("申请人"),
		edge.From("resource", Resource.Type).Ref("tasks").Unique().Required().Comment("目标资源"),
		edge.From("reviewer", User.Type).Ref("reviewed_tasks").Unique().Comment("审批人"),
	}
}

func (Task) Policy() ent.Policy {
	return privacy.Policy{
		Query: privacy.QueryPolicy{
			privacy.AlwaysAllowRule(),
		},
		Mutation: privacy.MutationPolicy{
			DenyMutationUnlessAdmin(),
		},
	}
}

