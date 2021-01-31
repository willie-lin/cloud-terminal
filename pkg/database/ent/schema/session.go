package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/dialect/entsql"
	"github.com/facebook/ent/schema"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
	"time"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Annotations of the User.
func (Session) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sessions"},
	}
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("protocol"),
		field.String("ip"),
		field.Int("port"),
		field.String("connectionId"),
		field.String("assetId"),
		field.String("username"),
		field.String("password"),
		field.String("creator"),
		field.String("clientIP"),
		field.Int("width"),
		field.Int("height"),
		field.String("status"),
		field.String("recording"),
		field.String("privateKey"),
		field.String("passphrase"),
		field.Int("code"),
		field.String("message"),
		field.Time("connectedTime").Default(time.Now),
		field.Time("disconnectedTime").Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("assets", Asset.Type),
	}
}
