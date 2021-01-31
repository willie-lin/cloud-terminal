package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/dialect/entsql"
	"github.com/facebook/ent/schema"
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
	return nil
}

// Edges of the Verification.
func (Verification) Edges() []ent.Edge {
	return nil
}
