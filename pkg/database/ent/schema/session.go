package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
		field.String("connection_id"),
		field.String("asset_id"),
		field.String("username"),
		field.String("password"),
		field.String("creator"),
		field.String("client_ip"),
		field.Int("width"),
		field.Int("height"),
		field.String("status"),
		field.String("recording"),
		field.String("private_key"),
		field.String("passphrase"),
		field.Int("code"),
		field.String("message"),
		field.Time("connected").Default(time.Now),
		field.Time("disconnected").Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("assets", Asset.Type),
	}
}
