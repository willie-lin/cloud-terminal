package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Resource holds the schema definition for the Resource entity.
type Resource struct {
	ent.Schema
}

func (Resource) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Resource.
func (Resource) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("name").Unique().NotEmpty(),                      // 资源的名称，例如"Database1"
		field.String("type").NotEmpty(),                               // 资源类型，例如"Database", "File", "IP"
		field.String("rrn").Unique().NotEmpty(),                       // 资源的 ARN
		field.JSON("properties", map[string]interface{}{}).Optional(), // 存储资源的其他属性
		field.JSON("tags", map[string]string{}).Optional(),            // 资源的标签
		field.String("description").Optional(),                        // 资源的描述
	}
}

// Edges of the Resource.
func (Resource) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("resources").Required().Comment("资源所属的账户"), // 资源属于哪个 Account
		edge.To("parent", Resource.Type).From("children").Comment("父资源"),                  //资源层级结构
	}
}

// Indexes of the Resource.
func (Resource) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
	}
}

// Policy defines the privacy policy of the Role.
//func (Resource) Policy() ent.Policy {
//	return privacy.Policy{
//		Query: privacy.QueryPolicy{
//			//rule.AllowEmailCheck(),
//			//rule.AllowIfAdmin(),            // 允许管理员进行查询
//			//rule.AllowIfOwner(),            // 允许用户查询自己的资料
//			//rule.AllowIfRole("SuperAdmin"), // 允许超级管理员进行查询
//			//rule.AllowIfTenantMember(),     // 允许同一租户成员进行查询
//			privacy.AlwaysAllowRule(),
//			//privacy.AlwaysDenyRule(),
//		},
//		Mutation: privacy.MutationPolicy{
//			//rule.DenyIfNoViewer(),
//			//rule.AllowIfAdmin(),            // 允许管理员进行操作
//			//rule.AllowIfOwner(),            // 允许用户修改自己的资料
//			//rule.AllowIfRole("SuperAdmin"), // 允许超级管理员进行操作
//			//privacy.AlwaysDenyRule(),
//			privacy.AlwaysAllowRule(),
//		},
//	}
//}
