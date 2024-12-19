package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/rule"
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

// Policy of the User.
func (User) Policy() ent.Policy {
	return privacy.Policy{
		Query: privacy.QueryPolicy{

			rule.AllowAdminQueryUser(),
			rule.AllowOwnerQueryUser(),
			privacy.AlwaysDenyRule(),
		},
		Mutation: privacy.MutationPolicy{

			rule.AllowAdminMutationUser(),
			rule.AllowOwnerMutationUser(),
			rule.AllowIfOwnerOrAdminMutationUser(),
			rule.DenyIfNotAdminOrSuperAdmin(),
			privacy.AlwaysDenyRule(),
		},
	}
}
