package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// Verification holds the schema definition for the Verification entity.
type Verification struct {
	ent.Schema
}

// Annotations of the User.
func (Verification) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "verifications"},
	}
}

// Fields of the Verification.
func (Verification) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("client_ip"),
		field.String("clientUserAgent"),
		field.Time("login_time").Default(time.Now),
		field.Time("logout_time").Default(time.Now).
			UpdateDefault(time.Now),
		field.Bool("remember"),
	}
}

// Edges of the Verification.
func (Verification) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
	}
}
