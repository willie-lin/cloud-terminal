package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// UserGroup holds the schema definition for the UserGroup entity.
type UserGroup struct {
	ent.Schema
}

// Fields of the UserGroup.
func (UserGroup) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
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
