package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"regexp"
	"time"
)

type User struct{ ent.Schema }

func (User) Mixin() []ent.Mixin { return []ent.Mixin{IDMixin{}, TimeMixin{}} }

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("avatar").Optional(),
		field.String("nickname").MinLen(2).MaxLen(30).Unique().Optional(),
		field.String("bio").MaxLen(200).Optional(),
		field.String("username").NotEmpty().MinLen(6).MaxLen(30).Unique(),
		field.String("password").NotEmpty().MinLen(20).MaxLen(120).Sensitive(),
		field.String("email").NotEmpty().Match(regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)).Unique(),
		field.Bool("email_verified").Default(true),
		field.String("phone_number").Optional().Default(""),
		field.Bool("phone_number_verified").Default(false),
		field.String("totp_secret").Optional(),
		field.Bool("online").Default(true),
		field.Bool("status").Default(true),
		field.Int("login_attempts").Default(0),
		field.Time("lockout_time").Optional(),
		field.Time("last_login_time").Default(time.Now),
		field.JSON("social_logins", map[string]string{}).Optional(),
		field.Bool("is_default").Default(false),
		field.String("ssh_public_key").Optional().Comment("SSH公钥"),
		field.JSON("attributes", map[string]interface{}{}).Optional().Comment("预留扩展属性"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("group", Group.Type).Ref("users").Unique().Comment("用户所属的组"),
		edge.To("roles", Role.Type).Comment("用户可 Assume 的角色"),
		edge.To("audit_logs", AuditLog.Type).Comment("用户的审计日志"),
		edge.To("access_policies", AccessPolicy.Type).Comment("直接分配给用户的策略"),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username").Unique(),
		index.Fields("email").Unique(),
	}
}
