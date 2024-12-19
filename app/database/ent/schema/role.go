package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/rule"
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
		edge.From("users", User.Type).Ref("roles"),    // 角色被多个用户拥有
		edge.To("permissions", Permission.Type),       // 角色拥有多个权限
		edge.From("tenant", Tenant.Type).Ref("roles"), // 正确用法
	}
}

// Indexes of the Role.
func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}

// Policy defines the privacy policy of the Role.
func (Role) Policy() ent.Policy {
	return privacy.Policy{
		Query: privacy.QueryPolicy{

			rule.AllowOnlySuperAdminQueryRole(),      // 仅允许 superadmin 查询
			rule.AllowIfAdminOrSuperAdminQueryRole(), // 允许 admin 和 superadmin 查询
			privacy.AlwaysDenyRule(),
		},
		Mutation: privacy.MutationPolicy{
			rule.AllowOnlySuperAdminMutationRole(),      // 仅允许 superadmin 变更
			rule.AllowIfAdminOrSuperAdminMutationRole(), // 允许 admin 和 superadmin 变更
			privacy.AlwaysDenyRule(),                    // 最后的拒绝策略
		},
	}
}
