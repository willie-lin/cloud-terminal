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
		field.String("ssh_public_key").Optional().Comment("SSH公钥，用于登录ContainerSSH"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("users").Unique().Required().Comment("用户所属的账户"),
		edge.To("role", Role.Type).Unique().Required().Comment("用户拥有的角色"),
		edge.To("audit_logs", AuditLog.Type).Comment("用户的审计日志"),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username").Unique(),
		index.Fields("email").Unique(),
	}
}
