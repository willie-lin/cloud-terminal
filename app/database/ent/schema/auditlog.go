package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// AuditLog holds the schema definition for the AuditLog entity.
type AuditLog struct {
	ent.Schema
}

func (AuditLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the AuditLog.
func (AuditLog) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Time("timestamp").Default(time.Now).Immutable(), // 操作发生的时间，不可修改
		field.Int("actor_id"),
		field.String("actor_username"),
		field.String("action").NotEmpty(),
		field.Int("resource_id").Optional(),
		field.Enum("resource_type").Values("user", "account", "role", "permission", "resource", "tenant", "platform").Optional(), // 可选的资源类型枚举
		field.String("ip_address").Optional(),                      // 用户 IP 地址 (可选)
		field.String("user_agent").Optional(),                      // 用户代理 (可选)
		field.JSON("details", map[string]interface{}{}).Optional(), // 操作详情
	}
}

// Edges of the AuditLog.
func (AuditLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("audit_logs").Comment("执行操作的用户"),
		edge.From("tenant", Tenant.Type).Ref("audit_logs").Comment("操作所属的租户"), // 添加与 Tenant 的关联
	}
}
