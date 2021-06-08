package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "users"},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		//field.String("ID").NotEmpty().Unique(),
		//field.String("ID").MaxLen(30).NotEmpty().Unique().Immutable(),
		field.String("id").Unique(),
		field.String("username").Unique(),
		field.String("password"),
		field.String("email"),
		field.String("nickname"),
		field.String("totpSecret"),
		field.Bool("online"),
		field.Bool("enable"),
		//field.Time("created_at").Default(time.Now),
		//field.Time("updated_at").Default(time.Now).
		//	UpdateDefault(time.Now),
		field.String("type"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("groups", Group.Type).Ref("users"),
		edge.To("assets", Asset.Type),
	}
}

// Mixin User
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Index of the User

//func (User) Indexes() []ent.Index {
//	return []ent.Index{
//		index.Fields("username"),
//	}
//}
