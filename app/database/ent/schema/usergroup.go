package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// UserGroup holds the schema definition for the UserGroup entity.
type UserGroup struct {
	ent.Schema
}

// Fields of the UserGroup.
func (UserGroup) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("group_name"),
		// 其他用户组信息字段
	}
}

// Edges of the UserGroup.
//func (UserGroup) Edges() []ent.Edge {
//	return []ent.Edge{
//		edge.To("users", User.Type),
//	}
//}
