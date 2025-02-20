package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"regexp"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Mixin MiXin Mixin User
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("avatar").Optional(), // 新增的头像字段
		field.String("nickname").MinLen(2).MaxLen(30).Unique().Optional(),
		field.String("bio").MaxLen(200).Optional(),
		field.String("username").NotEmpty().MinLen(6).MaxLen(30).Unique(),
		field.String("password").NotEmpty().MinLen(20).MaxLen(120).Sensitive(),
		field.String("email").NotEmpty().Match(regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)).Unique(),
		field.Bool("email_verified").Default(true), // 邮箱是否已验证
		field.String("phone_number").Optional().Default(""),
		field.Bool("phone_number_verified").Default(false), // 邮箱是否已验证
		field.String("totp_secret").Optional(),
		field.Bool("online").Default(true),
		field.Bool("status").Default(true),
		field.Int("login_attempts").Default(0), // 登录尝试次数
		field.Time("lockout_time").Optional(),  // 账户锁定时间
		field.Time("last_login_time").Default(time.Now),
		field.JSON("social_logins", map[string]string{}).Optional(), // 社交登录信息
		field.Bool("is_default").Default(false),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("users").Unique().Required().Comment("用户所属的账户"), // 多对一关系：一个 User 属于一个 Account
		edge.To("role", Role.Type).Unique().Required().Comment("用户拥有的角色"),                      // 正确：edge.To
		edge.To("audit_logs", AuditLog.Type).Comment("用户的审计日志"),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username").Unique(),
		index.Fields("email").Unique(),
	}
}

// Policy of the User.
//func (User) Policy() ent.Policy {
//	return privacy.Policy{
//		Query: privacy.QueryPolicy{
//			rule.AllowIfSuperAdminQueryUser(),
//			rule.AllowIfAdminQueryUser(),
//			rule.AllowIfOwnerQueryUser(),
//			privacy.AlwaysDenyRule(),
//		},
//		Mutation: privacy.MutationPolicy{
//			rule.AllowIfSuperAdminMutationUser(),
//			rule.AllowIfAdminMutationUser(),
//			rule.AllowIfOwnerMutationUser(),
//			privacy.AlwaysDenyRule(),
//		},
//	}
//}
