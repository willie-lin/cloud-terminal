package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Mixin MiXin Mixin User
func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").Unique().NotEmpty(),
		field.String("description").Optional(),
		field.Bool("is_disabled").Default(false), // 标记角色是否被禁用
		field.Bool("is_default").Default(false),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("roles").Unique(),
		edge.From("users", User.Type).Ref("role"),
		//edge.To("users", User.Type),
		edge.To("access_policies", AccessPolicy.Type),
		edge.To("child_roles", Role.Type).From("parent_role"), // 自引用：一个 Role 可以有多个子 Role
	}
}

// Indexes of the Role.
func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}

// Policy defines the privacy policy of the Role.
//func (Role) Policy() ent.Policy {
//	return privacy.Policy{
//		Query: privacy.QueryPolicy{
//			rule.AllowIfSuperAdminQueryRole(), // 允许 superuser 查询所有角色
//			rule.AllowIfAdminQueryRole(),      // 允许 admin 查询其租户下的角色
//			rule.AllowIfOwnerQueryRole(),      // 允许 user 查询自己的角色
//			privacy.AlwaysDenyRule(),          // 最后的拒绝策略
//			//privacy.AlwaysAllowRule(),
//		},
//		Mutation: privacy.MutationPolicy{
//			rule.AllowIfSuperAdminMutationRole(), // 允许 superuser 变更所有角色
//			rule.AllowIfAdminMutationRole(),      // 允许 admin 变更其租户下的角色
//			privacy.AlwaysDenyRule(),             // 最后的拒绝策略
//			//privacy.AlwaysAllowRule(),
//		},
//	}
//}
