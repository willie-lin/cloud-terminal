package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Environment holds the schema definition for the Environment entity.
type Environment struct {
	ent.Schema
}

func (Environment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}

func (Environment) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty().Comment("环境名称"),
		field.String("description").Optional(),
		field.String("image").NotEmpty().Comment("容器镜像"),
		field.Int("port").Default(22).Comment("SSH端口"),
		field.JSON("resource_limit", map[string]interface{}{}).Optional().Comment("资源限制 {cpu, memory, disk}"),
		field.JSON("env_vars", map[string]interface{}{}).Optional().Comment("环境变量"),
		field.JSON("volumes", []map[string]interface{}{}).Optional().Comment("挂载卷"),
		field.Enum("status").Values("active", "inactive").Default("active"),
	}
}

func (Environment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("environments").Unique().Comment("所属租户"),
		edge.From("access_policies", AccessPolicy.Type).Ref("environment").Unique().Comment("关联的访问策略"),
	}
}
