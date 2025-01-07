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
		field.Time("timestamp").Default(time.Now),
		field.Int("actor_id"),
		field.String("actor_username"),
		field.String("action"),
		field.Int("resource_id").Optional(),
		field.String("resource_type").Optional(),
		field.JSON("details", map[string]interface{}{}).Optional(), // 使用 JSON 存储操作详情
	}
}

// Edges of the AuditLog.
func (AuditLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("audit_logs").Required(),
	}
}
