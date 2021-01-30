package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/dialect/entsql"
	"github.com/facebook/ent/schema"
	"github.com/facebook/ent/schema/field"
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
		field.String("AssetId"),
		field.String("Username"),
		field.String("Password"),
		field.String("Creator"),
		field.String("ClientIP"),
		field.Int("Width"),
		field.Int("Height"),
		field.String("Status"),
		field.String("Recording"),
		field.String("PrivateKey"),
		field.String("Passphrase"),
		field.Int("Code"),
		field.String("Message"),
		field.Time("ConnectedTime"),
		field.Time("DisconnectedTime"),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return nil
}
