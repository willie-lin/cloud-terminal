package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Job holds the schema definition for the Job entity.
type Job struct {
	ent.Schema
}

// Annotations of the User.
func (Job) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "jobs"},
	}
}

// Fields of the Job.
func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.Int("cronjobid"),
		field.String("name"),
		field.String("func"),
		field.String("cron"),
		field.String("mode"),
		field.String("resourceIds"),
		field.String("status"),
		field.String("metadata"),
		//field.Time("created_at").Default(time.Now),
		//field.Time("updated_at").Default(time.Now).
		//	UpdateDefault(time.Now),
	}
}

// Edges of the Job.
func (Job) Edges() []ent.Edge {
	return nil
}
