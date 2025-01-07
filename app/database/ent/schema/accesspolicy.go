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
		field.JSON("statements", []map[string]interface{}{}), // 使用 JSON 存储策略语句
		field.String("resource_type").NotEmpty(),
		field.String("action").NotEmpty(),
		field.Bool("immutable").Default(false), // 添加 immutable 字段，默认为 false
	}
}

// Edges of the AccessPolicy.
func (AccessPolicy) Edges() []ent.Edge {
	return []ent.Edge{

		edge.From("tenant", Tenant.Type).Ref("access_policies").Required(),
		//edge.From("account", Account.Type).Ref("access_policies").Unique().Required(),
		//edge.From("users", User.Type).Ref("access_policies"),
		////edge.From("roles", Role.Type).Ref("access_policies"),
		//edge.To("roles", Role.Type),
		//edge.To("permissions", Permission.Type),
		//edge.To("resources", Resource.Type),
		////edge.To("resources", Resource.Type).StorageKey(edge.Table("access_policy_resources")),
		////edge.To("permissions", Permission.Type).StorageKey(edge.Table("access_policy_permissions")),

	}
}
