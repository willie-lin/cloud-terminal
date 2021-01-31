package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/dialect/entsql"
	"github.com/facebook/ent/schema"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
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
