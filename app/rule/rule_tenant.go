package rule

import (
	"context"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/app/viewer"
	"log"
)

// AllowOnlySuperAdminQueryTenant 是一个隐私规则示例，仅允许超级管理员对租户进行查询操作。
func AllowOnlySuperAdminQueryTenant() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)
		if v.RoleName != "superadmin" {
			log.Println("Denying query for non-superadmin")
			return privacy.Denyf("only superadmin can perform this action")
		}
		return privacy.Allow
	})
}

// AllowOnlySuperAdminMutationTenant 是一个隐私规则示例，仅允许超级管理员对租户进行变更操作。
func AllowOnlySuperAdminMutationTenant() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)
		if v.RoleName != "superadmin" {
			log.Println("Denying mutation for non-superadmin")
			return privacy.Denyf("only superadmin can perform this action")
		}
		return privacy.Allow
	})
}

// AllowIfAdminQueryTenant 允许管理员查询其租户下的资源。
func AllowIfAdminQueryTenant() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			log.Println("Allowing query for admin or superadmin")
			// 确保查询限于管理员的租户
			if tenantQuery, ok := q.(*ent.TenantQuery); ok {
				tenantQuery.Where(tenant.HasUsersWith(user.IDEQ(v.UserID)))
			} else if userQuery, ok := q.(*ent.UserQuery); ok {
				userQuery.Where(user.HasTenantWith(tenant.IDEQ(v.TenantID)))
			}
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AllowIfAdminMutationTenant 允许管理员变更其租户下的资源。
func AllowIfAdminMutationTenant() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			log.Println("Allowing mutation for admin or superadmin")
			// 确保变更限于管理员的租户
			if tenantMutation, ok := m.(*ent.TenantMutation); ok {
				tenantMutation.Where(tenant.HasUsersWith(user.IDEQ(v.UserID)))
			} else if userMutation, ok := m.(*ent.UserMutation); ok {
				userMutation.Where(user.HasTenantWith(tenant.IDEQ(v.TenantID)))
			}
			return privacy.Allow
		}
		return privacy.Skip
	})
}
