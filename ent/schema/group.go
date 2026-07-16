package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Group struct{ ent.Schema }

func (Group) Mixin() []ent.Mixin { return []ent.Mixin{IDMixin{}, TimeMixin{}} }

func (Group) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty().Comment("组名称"),
		field.String("description").Optional().Comment("组描述"),
		field.Enum("status").Values("active", "inactive").Default("active"),
		field.JSON("attributes", map[string]interface{}{}).Optional().Comment("预留扩展属性"),
	}
}

func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type).Comment("组内的用户"),
		edge.To("access_policies", AccessPolicy.Type).Comment("授予组的策略"),
	}
}
