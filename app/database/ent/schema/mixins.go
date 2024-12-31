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

//
//// Hooks of the AuditMixin.
//func (AuditMixin) Hooks() []ent.Hook {
//	return []ent.Hook{
//		AuditHook,
//	}
//}
//
//// A AuditHook is an example for audit-log hook.
//func AuditHook(next ent.Mutator) ent.Mutator {
//	type AuditLogger interface {
//		SetCreatedBy(uuid.UUID)
//		CreatedBy() (id uuid.UUID, exists bool)
//		SetCreatedAt(time.Time)
//		CreatedAt() (value time.Time, exists bool)
//		SetUpdatedBy(uuid.UUID)
//		UpdatedBy() (id uuid.UUID, exists bool)
//		SetUpdatedAt(time.Time)
//		UpdatedAt() (value time.Time, exists bool)
//	}
//	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
//		ml, ok := m.(AuditLogger)
//		if !ok {
//			return nil, fmt.Errorf("unexpected audit-log call from mutation type %T", m)
//		}
//
//		usr, err := viewer.UserFromContext(ctx) // 替换为你的用户获取逻辑
//		if err != nil {
//			// 根据你的需求决定是否返回错误或记录日志并继续
//			// 这里选择记录日志并继续，不中断操作
//			fmt.Printf("Error getting user from context: %v\n", err)
//			//return nil, fmt.Errorf("getting user from context: %w", err) // 如果需要中断操作则返回错误
//		}
//
//		switch op := m.Op(); {
//		case op.Is(ent.OpCreate):
//			ml.SetCreatedAt(time.Now())
//			if usr != nil {
//				ml.SetCreatedBy(usr.ID)
//			}
//		case op.Is(ent.OpUpdateOne | ent.OpUpdate):
//			ml.SetUpdatedAt(time.Now())
//			if usr != nil {
//				ml.SetUpdatedBy(usr.ID)
//			}
//		}
//		return next.Mutate(ctx, m)
//	})
//}
