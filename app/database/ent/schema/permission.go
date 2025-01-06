package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
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
		edge.From("account", Account.Type).Ref("permissions").Unique().Required(),
		edge.From("access_policies", AccessPolicy.Type).Ref("permissions"),
	}
}

// Indexes of the Permission.
func (Permission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}

// Policy defines the privacy policy of the Role.
//func (Permission) Policy() ent.Policy {
//	return privacy.Policy{
//		Query: privacy.QueryPolicy{
//			rule.AllowIfSuperAdminQueryPermission(),
//			rule.AllowIfAdminQueryPermission(),
//			rule.AllowIfOwnerQueryPermission(),
//			privacy.AlwaysDenyRule(),
//		},
//		Mutation: privacy.MutationPolicy{
//			rule.AllowIfSuperAdminMutationPermission(),
//			rule.AllowIfAdminMutationPermission(),
//			privacy.AlwaysDenyRule(),
//		},
//	}
//}
