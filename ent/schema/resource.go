package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Resource struct{ ent.Schema }

func (Resource) Mixin() []ent.Mixin { return []ent.Mixin{IDMixin{}, TimeMixin{}} }

func (Resource) Fields() []ent.Field {
	return []ent.Field{
		field.String("urn").Unique().NotEmpty().Comment("逻辑主键 urn:ct:<env>:<region>:<type>:<name>"),
		field.String("name").NotEmpty().Comment("资源名称"),
		field.Enum("type").Values("mysql", "redis", "k8s-service", "ssh", "rdp", "vnc", "telnet", "http", "custom").Default("ssh"),
		field.String("ip").NotEmpty(),
		field.Int("port").Default(22),
		field.Enum("env").Values("prod", "staging", "dev", "test", "dr").Default("dev"),
		field.String("region").Default("default"),
		field.String("description").Optional(),
		field.Enum("status").Values("active", "inactive").Default("active"),
		field.JSON("details", map[string]interface{}{}).Optional(),
		field.JSON("auth_data", map[string]interface{}{}).Optional().Sensitive(),
	}
}

func (Resource) Indexes() []ent.Index {
	return []ent.Index{index.Fields("urn").Unique(), index.Fields("type", "env", "region")}
}

func (Resource) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("resources").Unique(),
		edge.To("audit_logs", AuditLog.Type),
		edge.To("policies", AccessPolicy.Type).Comment("附加到资源的策略，实现 ResourcePolicy"),
	}
}
