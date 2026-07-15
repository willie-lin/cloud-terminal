package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Account holds the schema definition for the Account entity.
// UNR core: User → Account ← Resource (the binding layer)
type Account struct {
	ent.Schema
}

func (Account) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").NotEmpty().Comment("账号名称"),
		field.String("description").Optional().Comment("账号描述"),
		field.Enum("status").Values("active", "inactive", "suspended").Default("active"),
	}
}

func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type).Comment("账号下的用户"),
		edge.To("roles", Role.Type).Comment("账号下的角色"),
		edge.To("access_policies", AccessPolicy.Type).Comment("账号下的访问策略"),
		edge.To("resource", Resource.Type).Unique().Comment("绑定的目标服务器"),
	}
}
