package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
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
		field.String("password").NotEmpty().MinLen(8).MaxLen(120).Sensitive(),
		field.String("email").NotEmpty().Match(regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)).Unique(),
		field.String("phone").Optional(),
		field.String("totp_secret").Optional(),
		field.Bool("online").Default(true),
		field.Bool("enable_type").Default(true),
		field.Time("last_login_time").Default(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("roles", Role.Type),                            // 用户拥有多个角色
		edge.From("tenant", Tenant.Type).Ref("users").Unique(), // 用户属于一个租户
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username").Unique(),
		index.Fields("email").Unique(),
	}
}

// Policy defines the privacy policy of the User.
func (User) Policy() ent.Policy {
	return privacy.Policy{
		Query: privacy.QueryPolicy{
			//rule.AllowIfAdmin(),            // 允许管理员进行查询
			//rule.AllowIfOwner(),            // 允许用户查询自己的资料
			//rule.AllowIfRole("SuperAdmin"), // 允许超级管理员进行查询
			//rule.AllowIfTenantMember(),     // 允许同一租户成员进行查询
			privacy.AlwaysAllowRule(),
			//privacy.AlwaysDenyRule(),
		},
		Mutation: privacy.MutationPolicy{
			//rule.DenyIfNoViewer(),
			//rule.AllowIfAdmin(),            // 允许管理员进行操作
			//rule.AllowIfOwner(),            // 允许用户修改自己的资料
			//rule.AllowIfRole("SuperAdmin"), // 允许超级管理员进行操作
			//privacy.AlwaysDenyRule(),
			privacy.AlwaysAllowRule(),
		},
	}
}
