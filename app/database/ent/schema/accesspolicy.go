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

type PolicyStatement struct {
	Effect    string      `json:"Effect"`
	Actions   []string    `json:"Action"`
	Resources []string    `json:"Resource"`
	Condition interface{} `json:"Condition,omitempty"`
}

// Fields of the AccessPolicy.
func (AccessPolicy) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").NotEmpty(),
		field.String("description").Optional(),
		//field.Enum("effect").Values("allow", "deny").Default("allow"),          // 策略效果
		field.JSON("statements", []PolicyStatement{}).Default([]PolicyStatement{}), // 使用 JSON 存储策略语句
		field.Bool("immutable").Default(false),                                     // 添加 immutable 字段，默认为 false
	}
}

// Edges of the AccessPolicy.
func (AccessPolicy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("access_policies").Unique(),
		//edge.From("user", User.Type).Ref("access_policies"),
		edge.To("roles", Role.Type),
	}
}
