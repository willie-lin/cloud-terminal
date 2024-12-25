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
			//rule.AllowOnlySuperAdminQueryTenant(), // 仅允许 superadmin 查询
			rule.AllowIfAdminQueryTenant(), // 允许 admin 查询其租户下的资源
			privacy.AlwaysDenyRule(),       // 最后的拒绝策略
		},
		Mutation: privacy.MutationPolicy{
			//rule.AllowOnlySuperAdminMutationTenant(), // 允许 superuser 变更所有租户
			rule.AllowIfAdminMutationTenant(), // 允许 admin 变更其租户下的资源
			privacy.AlwaysDenyRule(),
		},
	}
}
