package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Platform holds the schema definition for the Platform entity.
type Platform struct {
	ent.Schema
}

func (Platform) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Platform.
func (Platform) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").Unique().NotEmpty(),
		field.String("description").Optional(),
		field.String("region").Optional(),
		field.String("version").Optional(),
		field.Enum("status").Values("active", "maintenance", "disabled").Default("active").Optional(), // 平台状态 (可选)
		field.JSON("config", map[string]interface{}{}).Optional(),                                     // 平台配置 (可选)
	}
}

// Edges of the Platform.
func (Platform) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tenants", Tenant.Type).Comment("平台下的所有租户"), // 一对多关系：一个 Platform 可以有多个 Tenant
	}
}
