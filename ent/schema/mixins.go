package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
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

// AuditMixin implements the ent.Mixin for sharing
// audit-log capabilities with package schemas.
type AuditMixin struct {
	mixin.Schema
}

// Fields of the AuditMixin.
func (AuditMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("created_by", uuid.UUID{}).Optional(),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.UUID("updated_by", uuid.UUID{}).Optional(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// BaseMixin provides the standard created_at/updated_at timestamp fields.
// Used by all new-architecture schema entities.
type BaseMixin struct {
	mixin.Schema
}

func (BaseMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Immutable().
			Default(time.Now).Comment("创建时间"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).Comment("更新时间"),
	}
}
