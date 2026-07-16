package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type AuditLog struct{ ent.Schema }

func (AuditLog) Mixin() []ent.Mixin { return []ent.Mixin{IDMixin{}, TimeMixin{}} }

func (AuditLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("session_id").Unique().Comment("会话ID"),
		field.String("username").Comment("用户名"),
		field.String("action").Comment("操作动作"),
		field.String("result").Comment("操作结果"),
		field.Time("started_at").Comment("开始时间"),
		field.Time("ended_at").Optional().Nillable().Comment("结束时间"),
		field.String("resource_urn_snapshot").Optional().Comment("操作时刻的资源 URN 快照"),
		field.JSON("detail", map[string]interface{}{}).Optional().Comment("操作详情"),
		field.String("s3_path").Optional().Comment("S3存储路径"),
	}
}

func (AuditLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("audit_logs").Unique(),
		edge.From("resource", Resource.Type).Ref("audit_logs").Unique().Comment("被操作的资源"),
	}
}

func (AuditLog) Indexes() []ent.Index {
	return []ent.Index{index.Fields("session_id").Unique(), index.Fields("started_at")}
}
