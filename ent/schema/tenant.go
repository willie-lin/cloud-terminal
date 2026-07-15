package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
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
		field.Enum("status").Values("active", "inactive", "suspended").Default("active"),
	}
}

// Edges of the Tenant.
func (Tenant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("environments", Environment.Type).Comment("租户下的环境模板"),
		edge.To("resources", Resource.Type).Comment("租户下的资源"),
		edge.To("access_policies", AccessPolicy.Type).Comment("租户下的访问策略"),
	}
}
