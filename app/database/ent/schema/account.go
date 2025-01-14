package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Account holds the schema definition for the Account entity.
type Account struct {
	ent.Schema
}

func (Account) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").Unique().NotEmpty(),
		field.String("description").Optional(), // 可选的描述信息
		field.Enum("status").Values("active", "suspended", "deleted").Default("active"),
	}
}

func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("accounts").Unique().Required().Comment("账户所属的租户"), // 多对一关系：一个 Account 属于一个 Tenant
		edge.To("users", User.Type).Comment("账户下的所有用户"),                                         // 一对多关系：一个 Account 可以有多个 User
		edge.To("roles", Role.Type).Comment("账户下的所有角色"),
		edge.To("resources", Resource.Type).Comment("账户下的所有资源"),
		edge.To("access_policies", AccessPolicy.Type).Comment("直接关联到账户的策略"), // 账户与策略的多对多关系
	}
}
