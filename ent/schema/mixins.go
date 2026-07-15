package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
	"time"
)

// TimeMixin provides created_at/updated_at timestamp fields.
type TimeMixin struct {
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Immutable().
			Default(time.Now).Comment("创建时间"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).Comment("更新时间"),
	}
}

// IDMixin provides a string UUID primary key.
type IDMixin struct {
	mixin.Schema
}

func (IDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			DefaultFunc(func() string {
				return uuid.NewString()
			}).
			Immutable().
			Comment("UUID primary key"),
	}
}
