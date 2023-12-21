package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"time"
)

// TimeMixin implements the ent.Mixin for sharing
// time fields with package schemas.

type TimeMixin struct {
	// We embed the `mixin.Schema` to avoid
	// implementing the rest of the methods.
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		// Immutable 函数告诉我们，生下来后你的出生年月就定了，不能改变。你明明是半老徐娘就不能说自己芳龄十八。
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// editmixin

type EditMixin struct {
	mixin.Schema
}

func (EditMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("creator").Optional().NotEmpty(),
		field.String("editor").Optional(),
		field.Float("deleted").SchemaType(map[string]string{
			dialect.MySQL:    "decimal(1,0)",
			dialect.Postgres: "numeric",
		}),
	}
}
