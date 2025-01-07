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
	}
}

// Edges of the Platform.
func (Platform) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tenants", Tenant.Type), // 一对多关系：一个 Platform 可以有多个 Tenant
	}
}
