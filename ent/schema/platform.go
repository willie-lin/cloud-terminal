package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Platform holds the schema definition for the Platform entity.
type Platform struct {
	ent.Schema
}

func (Platform) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
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
		field.Enum("status").Values("active", "maintenance", "disabled").Default("active").Optional(),
		field.JSON("config", map[string]interface{}{}).Optional(),
	}
}
