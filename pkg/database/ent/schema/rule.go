package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Rule holds the schema definition for the Rule entity.
type Rule struct {
	ent.Schema
}

// Annotations of the User.
func (Asset) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "rules"},
	}
}

// Fields of the Rule.
func (Rule) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("rule"),
		field.String("ip"),
		field.String("source"),
		field.Int64("priority"),
	}
}

// Edges of the Rule.
func (Rule) Edges() []ent.Edge {
	return nil
}
