package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/dialect/entsql"
	"github.com/facebook/ent/schema"
	"github.com/facebook/ent/schema/field"
	"time"
)

// Credential holds the schema definition for the Credential entity.
type Credential struct {
	ent.Schema
}

// Annotations of the User.
func (Credential) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "credentials"},
	}
}

// Fields of the Credential.
func (Credential) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("name").Unique(),
		field.String("type"),
		field.String("username"),
		field.String("password"),
		field.String("private_key"),
		field.String("passphrase"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Credential.
func (Credential) Edges() []ent.Edge {
	return nil
}
