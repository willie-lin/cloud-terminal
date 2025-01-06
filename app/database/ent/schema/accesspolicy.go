package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// AccessPolicy holds the schema definition for the AccessPolicy entity.
type AccessPolicy struct {
	ent.Schema
}

func (AccessPolicy) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the AccessPolicy.
func (AccessPolicy) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").Unique().NotEmpty(),
		field.String("description").Optional(),
		field.Enum("effect").Values("allow", "deny").Default("allow"), // 策略效果
		//field.JSON("condition", &model.Condition{}).Optional(),
		field.String("resource_type").NotEmpty(),
		field.String("action").NotEmpty(),
	}
}

// Edges of the AccessPolicy.
func (AccessPolicy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("access_policies").Unique().Required(),
		edge.From("roles", Role.Type).Ref("access_policies"),
		edge.To("resources", Resource.Type).StorageKey(edge.Table("access_policy_resources")),
		edge.To("permissions", Permission.Type).StorageKey(edge.Table("access_policy_permissions")),
		edge.From("users", User.Type).Ref("access_policies"),
	}
}
