package rule

import (
	"context"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/permission"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/role"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/app/viewer"
	"log"
)

// AllowIfSuperAdminQueryPermission  允许超级用户查询所有资源。
func AllowIfSuperAdminQueryPermission() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		if v.RoleName == "superadmin" {
			log.Println("Allowing query for superadmin")
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AllowIfSuperAdminMutationPermission   允许超级用户变更所有资源。
func AllowIfSuperAdminMutationPermission() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		if v.RoleName == "superadmin" {
			log.Println("Allowing mutation for superadmin")
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AllowIfAdminQueryPermission 允许管理员查询其租户下的权限。
func AllowIfAdminQueryPermission() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			log.Println("Allowing query for admin or superadmin")
			// 确保查询限于管理员的租户
			if permissionQuery, ok := q.(*ent.PermissionQuery); ok {
				permissionQuery.Where(permission.HasRolesWith(role.HasUsersWith(user.HasTenantWith(tenant.IDEQ(v.TenantID)))))
			}
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AllowIfAdminMutationPermission 允许管理员变更其租户下的权限。
func AllowIfAdminMutationPermission() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			log.Println("Allowing mutation for admin or superadmin")
			// 确保变更限于管理员的租户
			if permissionMutation, ok := m.(*ent.PermissionMutation); ok {
				permissionMutation.Where(permission.HasRolesWith(role.HasUsersWith(user.HasTenantWith(tenant.IDEQ(v.TenantID)))))
			}
			return privacy.Allow
		}
		return privacy.Skip
	})
}

//// AllowIfOwnerQueryPermission 允许用户查询自己的权限。
//func AllowIfOwnerQueryPermission() privacy.QueryRule {
//	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
//		v := viewer.FromContext(ctx)
//		if v == nil {
//			log.Println("No viewer in context")
//			return privacy.Denyf("viewer is not authenticated")
//		}
//		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)
//		if permissionQuery, ok := q.(*ent.PermissionQuery); ok {
//			permissionQuery.Where(permission.HasRolesWith(role.HasUsersWith(user.IDEQ(v.UserID))))
//			log.Println("Allowing query for user permissions")
//			return privacy.Allow
//		}
//		//log.Println("Denying query for non-owner")
//		//return privacy.Denyf("only owner can perform this action")
//		log.Println("Denying query for non-owner")
//		return privacy.Skip // 重要：这里要返回 privacy.Skip，而不是 privacy.Denyf
//	})
//
//}

func AllowIfOwnerQueryPermission() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)
		if permissionQuery, ok := q.(*ent.PermissionQuery); ok {
			permissionQuery.Where(permission.HasRolesWith(role.HasUsersWith(user.IDEQ(v.UserID))))
			log.Println("Allowing query for user permissions")
			return privacy.Allow
		}
		log.Println("Skipping query for non-permission entity")
		return privacy.Skip
	})
}
