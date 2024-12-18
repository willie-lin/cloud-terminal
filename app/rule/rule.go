package rule

import (
	"context"
	"fmt"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/app/viewer"
)

// AllowIfOwner 是一个隐私规则示例，允许用户访问自己的数据。
func AllowIfOwner() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			fmt.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		fmt.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)

		// 检查查询是否是 User 查询，并添加查询条件
		userQuery, ok := q.(*ent.UserQuery)
		if !ok {
			return privacy.Denyf("not a UserQuery")
		}
		userQuery.Where(user.IDEQ(v.UserID))

		fmt.Println("Allowing access for owner")
		return privacy.Allow
	})
}

// AllowIfOwnerMutation 是一个隐私规则示例，允许用户修改自己的数据。
func AllowIfOwnerMutation() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			fmt.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		fmt.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)

		// 检查变更是否是 User 变更，并添加查询条件
		if userMutation, ok := m.(*ent.UserMutation); ok {
			ids, err := userMutation.IDs(ctx)
			if err == nil && len(ids) == 1 && ids[0] == v.UserID {
				fmt.Println("Allowing mutation for owner")
				return privacy.Allow
			}
		}

		fmt.Println("Denying mutation for non-owner")
		return privacy.Denyf("only owner can perform this action")
	})
}

// AllowIfAdmin 是一个查询规则，允许管理员和超级管理员进行查询
func AllowIfAdmin() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			fmt.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		fmt.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			fmt.Println("Allowing access for admin or superadmin")
			return privacy.Allow
		}
		fmt.Println("Denying access for non-admin or non-superadmin")
		return privacy.Skip
	})
}

// AllowIfAdminMutation 是一个变更规则，允许管理员和超级管理员进行变更操作
// AllowIfAdminMutation 是一个变更规则，允许管理员和超级管理员进行变更操作
func AllowIfAdminMutation() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			fmt.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		fmt.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			fmt.Println("Allowing mutation for admin or superadmin")
			return privacy.Allow
		}
		fmt.Println("Denying mutation for non-admin or non-superadmin")
		return privacy.Skip
	})
}

// DenyIfNotOwnerOrAdmin 是一个隐私规则示例，拒绝非所有者或非管理员进行特定操作。
func DenyIfNotOwnerOrAdmin() privacy.MutationRule {
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
