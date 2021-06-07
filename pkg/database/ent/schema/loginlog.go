package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

// LoginLog holds the schema definition for the LoginLog entity.
type LoginLog struct {
	ent.Schema
}

// Annotations of the User.
func (LoginLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "loginlogs"},
	}
}

// Fields of the LoginLog.
func (LoginLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("user_id"),
		field.String("client_ip"),
		field.String("clent_uset_agent"),
		field.Time("login_time").Default(time.Now),
		field.Time("logout_time").Default(time.Now).UpdateDefault(time.Now),
		field.Bool("remember"),
	}
}

// Edges of the LoginLog.
func (LoginLog) Edges() []ent.Edge {
	return nil
}
