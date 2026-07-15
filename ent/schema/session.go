package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

func (Session) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}

func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String("session_id").Unique().Immutable().Comment("会话ID"),
		field.String("principal_urn").Comment("主体URN"),
		field.String("resource_urn").Optional().Comment("资源URN"),
		field.String("environment_urn").Optional().Comment("环境URN"),
		field.String("account_urn").Optional().Comment("账号URN"),
		field.Enum("mode").Values("container", "proxy", "unknown").Default("unknown").Comment("连接模式"),
		field.Enum("status").Values("active", "closed", "timeout", "error").Default("active").Comment("会话状态"),
		field.Time("started_at").Comment("开始时间"),
		field.Time("ended_at").Optional().Nillable().Comment("结束时间"),
		field.String("remote_address").Optional().Comment("来源地址"),
	}
}
