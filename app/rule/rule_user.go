package rule

import (
	"context"
	"fmt"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/app/viewer"
	"log"
)

// AllowAdminQueryUser 是一个隐私规则示例，允许管理员对用户进行查询操作。
func AllowAdminQueryUser() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			log.Println("Allowing query for admin or superadmin")
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AllowOwnerQueryUser 是一个隐私规则示例，允许所有者对用户进行查询操作。
func AllowOwnerQueryUser() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)

		userQuery, ok := q.(*ent.UserQuery)
		if !ok {
			return privacy.Denyf("not a UserQuery")
		}
		userQuery.Where(user.IDEQ(v.UserID))

		log.Println("Allowing query for owner")
		return privacy.Allow
	})
}

// AllowOwnerMutationUser 是一个隐私规则示例，允许用户修改自己的数据。
func AllowOwnerMutationUser() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)

		if userMutation, ok := m.(*ent.UserMutation); ok {
			ids, err := userMutation.IDs(ctx)
			if err == nil && len(ids) == 1 && ids[0] == v.UserID {
				log.Println("Allowing mutation for owner")
				return privacy.Allow
			}
		}

		log.Println("Denying mutation for non-owner")
		return privacy.Denyf("only owner can perform this action")
	})
}

// AllowAdminMutationUser 是一个变更规则，允许管理员和超级管理员进行变更操作
func AllowAdminMutationUser() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			log.Println("Allowing mutation for admin or superadmin")
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AllowIfOwnerOrAdminMutationUser 是一个隐私规则示例，允许所有者或管理员进行用户实体的变更操作。
func AllowIfOwnerOrAdminMutationUser() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			fmt.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		fmt.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)

		// 如果是管理员或超级管理员，允许操作
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			fmt.Println("Allowing mutation for admin or superadmin")
			return privacy.Allow
		}

		// 检查是否是所有者
		if userMutation, ok := m.(*ent.UserMutation); ok {
			ids, err := userMutation.IDs(ctx)
			if err == nil && len(ids) == 1 && ids[0] == v.UserID {
				fmt.Println("Allowing mutation for owner")
				return privacy.Allow
			}
		}

		fmt.Println("Denying mutation for non-owner and non-admin")
		return privacy.Denyf("only owner or admin can perform this action")
	})
}
