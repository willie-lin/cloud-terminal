package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "users"},
	}
}

// Mixin MiXin Mixin User
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("username").Unique(),
		field.String("password"),
		field.String("email").Unique(),
		field.String("nickname"),
		field.String("totp_secret"),
		field.Bool("online").Default(false),
		field.Enum("enable_type").Values("Enabled", "Disabled").Default("Enabled"),
		field.Enum("user_type").Values("Admin", "Auditor", "SuperUser", "User").Default("User"),
		field.Time("last_login_time").Default(time.Now),
	}
}

// Edges of the User.
//func (User) Edges() []ent.Edge {
//	return []ent.Edge{
//		edge.To("groups", UserGroups.Type),
//		edge.To("assets", Assets.Type),
//	}
//}
