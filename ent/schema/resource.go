package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Resource holds the schema definition for the Resource entity.
type Resource struct {
	ent.Schema
}

func (Resource) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}

func (Resource) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Comment("资源名称"),
		field.String("host").NotEmpty().Comment("目标地址"),
		field.Int("port").Default(22).Comment("目标端口"),
		field.Enum("type").Values("ssh", "rdp", "vnc", "telnet").Comment("协议类型"),
		field.String("description").Optional(),
		field.Enum("status").Values("active", "inactive").Default("active"),
		field.JSON("metadata", map[string]interface{}{}).Optional().Comment("扩展属性"),
	}
}

func (Resource) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("resources").Unique().Comment("所属租户"),
		edge.From("accounts", Account.Type).Ref("resource").Unique().Comment("关联的凭据"),
	}
}
