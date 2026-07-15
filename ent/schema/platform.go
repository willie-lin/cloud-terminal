package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Platform holds the schema definition for the Platform entity.
type Platform struct {
	ent.Schema
}

func (Platform) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
		IDMixin{},
	}
}

// Fields of the Platform.
func (Platform) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty(),
		field.String("description").Optional(),
		field.String("region").Optional(),
		field.String("version").Optional(),
		field.Enum("status").Values("active", "maintenance", "disabled").Default("active").Optional(),
		field.JSON("config", map[string]interface{}{}).Optional(),
	}
}
