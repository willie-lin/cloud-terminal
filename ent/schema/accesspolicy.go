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
		BaseMixin{},
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
		field.JSON("statements", []PolicyStatement{}).Default([]PolicyStatement{}),
		field.Bool("immutable").Default(false),
		field.Int("priority").Default(0).Comment("策略优先级，数值越小优先级越高"),
	}
}

// Edges of the AccessPolicy.
func (AccessPolicy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("access_policies").Comment("应用此策略的账户"),
		edge.From("roles", Role.Type).Ref("access_policies").Comment("分配此策略的角色"),
		edge.From("tenant", Tenant.Type).Ref("access_policies").Unique().Comment("所属租户"),
		edge.To("environment", Environment.Type).Unique().Comment("关联的环境模板"),
	}
}
