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
		//field.String("email").Unique().NotEmpty(),
	}
}

func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("accounts").Unique().Required(), // 多对一关系：一个 Account 属于一个 Tenant
		edge.To("users", User.Type),                   // 一对多关系：一个 Account 可以有多个 User
		edge.To("roles", Role.Type),                   // 一对多关系：一个 Account 可以有多个 Role
		edge.To("access_policies", AccessPolicy.Type), // 一对多关系：一个 Account 可以有多个 AccessPolicy
		edge.To("permissions", Permission.Type),       // 一对多关系：一个 Account 可以有多个 Permission
	}
}
