package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// Asset holds the schema definition for the Asset entity.
type Asset struct {
	ent.Schema
}

// Annotations of the User.
func (Asset) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "assets"},
	}
}

// Fields of the Assets.
func (Asset) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("name").Unique(),
		field.String("ip"),
		field.String("protocol"),
		field.Int("port"),
		field.String("account_type"),
		field.String("username"),
		field.String("password"),
		field.String("credential_id"),
		field.String("private_key"),
		field.String("passphrase"),
		field.String("description"),
		field.Bool("active"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).
			UpdateDefault(time.Now),
		field.String("tags"),
	}
}

// Edges of the Assets.
func (Asset) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("sessions", Session.Type).
			Ref("assets").
			Unique(),
	}
}
