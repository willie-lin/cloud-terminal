package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/dialect/entsql"
	"github.com/facebook/ent/schema"
	"github.com/facebook/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Users"},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		//field.String("ID").NotEmpty().Unique(),
		//field.String("ID").MaxLen(30).NotEmpty().Unique().Immutable(),
		field.UUID("id", uuid.UUID{}),
		field.String("Username"),
		field.String("Password"),
		field.String("Nickname"),
		field.String("TOTPSecret"),
		field.Bool("Online"),
		field.Bool("Enable"),
		field.Time("created_at").Default(time.Now),
		field.String("Type"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
