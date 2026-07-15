package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// AuditLog holds the schema definition for the AuditLog entity.
type AuditLog struct {
	ent.Schema
}

func (AuditLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (AuditLog) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("session_id").Unique().Comment("会话ID"),
		field.String("username").Comment("用户名"),
		field.String("action").Comment("操作动作"),
		field.String("result").Comment("操作结果"),
		field.Time("started_at").Comment("操作开始时间"),
		field.Time("ended_at").Optional().Nillable().Comment("操作结束时间"),
		field.JSON("detail", map[string]interface{}{}).Optional().Comment("操作详情"),
		field.String("s3_path").Optional().Comment("S3存储路径"),
	}
}

func (AuditLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("audit_logs").Unique(),
	}
}

func (AuditLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("session_id").Unique(),
		index.Fields("started_at"),
	}
}
