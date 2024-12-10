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

// Resource holds the schema definition for the Resource entity.
type Resource struct {
	ent.Schema
}

func (Resource) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Resource.
func (Resource) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("type").NotEmpty(),
		field.String("identifier").NotEmpty(),
		field.String("description").Optional(),
	}
}

// Edges of the Resource.
func (Resource) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).Ref("resources").Unique(),
		edge.To("permissions", Permission.Type),
	}
}

// Indexes of the Resource.
func (Resource) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
		index.Fields("identifier").Unique(),
	}
}

// Policy defines the privacy policy of the Role.
func (Resource) Policy() ent.Policy {
	return privacy.Policy{
		Mutation: privacy.MutationPolicy{
			rule.DenyIfNoViewer(),
			rule.AllowIfAdmin(),                 // 允许管理员进行操作
			rule.AllowIfOwner(),                 // 允许资源拥有者进行修改
			rule.AllowIfRole("ResourceManager"), // 允许特定角色进行修改
			rule.AllowIfTenantMember(),          // 允许同一租户成员进行操作
			privacy.AlwaysDenyRule(),
		},
		Query: privacy.QueryPolicy{
			rule.AllowIfAdmin(),        // 允许管理员进行查询
			rule.AllowIfOwner(),        // 允许资源拥有者进行查询
			rule.AllowIfTenantMember(), // 允许同一租户成员进行查询
			privacy.AlwaysDenyRule(),
		},
	}
}
