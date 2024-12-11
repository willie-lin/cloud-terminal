package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
)

// Tenant holds the schema definition for the Tenant entity.
type Tenant struct {
	ent.Schema
}

func (Tenant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Tenant.
func (Tenant) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").Unique().NotEmpty(),
		field.String("description").Optional(),
	}
}

// Edges of the Tenant.
func (Tenant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
		edge.To("roles", Role.Type),
		edge.To("resources", Resource.Type),
		edge.To("permissions", Permission.Type),
	}
}

// Indexes of the Tenant.
func (Tenant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}

// Policy defines the privacy policy of the Role.
func (Tenant) Policy() ent.Policy {
	return privacy.Policy{
		Query: privacy.QueryPolicy{
			//rule.AllowEmailCheck(),
			//rule.AllowIfAdmin(),            // 允许管理员进行查询
			//rule.AllowIfOwner(),            // 允许用户查询自己的资料
			//rule.AllowIfRole("SuperAdmin"), // 允许超级管理员进行查询
			//rule.AllowIfTenantMember(),     // 允许同一租户成员进行查询
			privacy.AlwaysAllowRule(),
			//privacy.AlwaysDenyRule(),
		},
		Mutation: privacy.MutationPolicy{
			//rule.DenyIfNoViewer(),
			//rule.AllowIfAdmin(),            // 允许管理员进行操作
			//rule.AllowIfOwner(),            // 允许用户修改自己的资料
			//rule.AllowIfRole("SuperAdmin"), // 允许超级管理员进行操作
			//privacy.AlwaysDenyRule(),
			privacy.AlwaysAllowRule(),
		},
	}
}
