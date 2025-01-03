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

// Permission holds the schema definition for the Permission entity.
type Permission struct {
	ent.Schema
}

// Mixin MiXin Mixin User
func (Permission) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Permission.
func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").Unique().NotEmpty(),
		field.Strings("actions"),
		field.String("resource_type"),
		field.String("description").Optional(),
		field.Bool("is_disabled").Default(false), // 标记角色是否被禁用
	}
}

// Edges of the Permission.
func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("roles", Role.Type).Ref("permissions"),
		edge.From("tenant", Tenant.Type).Ref("permissions").Unique(), // 权限到租户的关系

		//edge.From("tenant", Tenant.Type).Ref("permissions").Unique(),

		//edge.From("resources", Resource.Type).Ref("permissions"),
		//edge.To("resource", Resource.Type), // 新增的资源关系
		// Add more edges here
		//edge.To("users", User.Type),
		//edge.To("roles", Role.Type), // 权限分配给多个角色
		//edge.To("actions", Action.Type), // 权限包含多个操作
		edge.To("resources", Resource.Type), // 权限针对多个资源类型
	}
}

// Indexes of the Permission.
func (Permission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}

// Policy defines the privacy policy of the Role.
func (Permission) Policy() ent.Policy {
	return privacy.Policy{
		Query: privacy.QueryPolicy{
			rule.AllowIfSuperAdminQueryPermission(),
			rule.AllowIfAdminQueryPermission(),
			rule.AllowIfOwnerQueryPermission(),
			privacy.AlwaysDenyRule(),
		},
		Mutation: privacy.MutationPolicy{
			rule.AllowIfSuperAdminMutationPermission(),
			rule.AllowIfAdminMutationPermission(),
			privacy.AlwaysDenyRule(),
		},
	}
}
