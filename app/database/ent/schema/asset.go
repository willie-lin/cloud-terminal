package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

func (Asset) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "assets"},
	}
}

// Mixin MiXin Mixin User
func (Asset) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Asset holds the schema definition for the Asset entity.
type Asset struct {
	ent.Schema
}

// Fields of the Asset.
func (Asset) Fields() []ent.Field {
	return []ent.Field{
		//field.UUID("id", uuid.UUID{}).Unique(),
		//field.String("id", uuid.ClockSequence()),
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("asset_name"),
		field.String("asset_type"),
		field.String("asset_details"),
		field.Int("group_id"),
		// 其他资产信息字段
	}
}

// Edges of the Asset.
//func (Asset) Edges() []ent.Edge {
//	return []ent.Edge{
//		edge.From("group", AssetGroup.Type).Ref("assets").Unique(),
//		edge.To("users", User.Type),
//	}
//}
