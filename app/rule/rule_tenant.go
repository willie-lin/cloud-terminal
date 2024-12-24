package rule

import (
	"context"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"github.com/willie-lin/cloud-terminal/app/viewer"
)

//
//// AllowOnlySuperAdminQueryTenant 是一个隐私规则示例，仅允许超级管理员对租户进行查询操作。
//func AllowOnlySuperAdminQueryTenant() privacy.QueryRule {
//	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
//		v := viewer.FromContext(ctx)
//		if v == nil {
//			log.Println("No viewer in context")
//			return privacy.Denyf("viewer is not authenticated")
//		}
//		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)
//		if v.RoleName != "superadmin" {
//			log.Println("Denying query for non-superadmin")
//			return privacy.Denyf("only superadmin can perform this action")
//		}
//		return privacy.Allow
//	})
//}
//
//// AllowOnlySuperAdminMutationTenant 是一个隐私规则示例，仅允许超级管理员对租户进行变更操作。
//func AllowOnlySuperAdminMutationTenant() privacy.MutationRule {
//	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
//		v := viewer.FromContext(ctx)
//		if v == nil {
//			log.Println("No viewer in context")
//			return privacy.Denyf("viewer is not authenticated")
//		}
//		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)
//		if v.RoleName != "superadmin" {
//			log.Println("Denying mutation for non-superadmin")
//			return privacy.Denyf("only superadmin can perform this action")
//		}
//		return privacy.Allow
//	})
//}
//func AllowIfAdminQueryTenant() privacy.QueryRule {
//	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
//		v := viewer.FromContext(ctx)
//		if v == nil {
//			log.Println("No viewer in context")
//			return privacy.Denyf("viewer is not authenticated")
//		}
//		if v.RoleName == "admin" || v.RoleName == "superadmin" {
//			log.Println("Allowing query for admin or superadmin")
//			if tenantQuery, ok := q.(*ent.TenantQuery); ok {
//				tenantQuery.Where(tenant.IDEQ(v.TenantID))
//			}
//			return privacy.Allow
//		}
//		return privacy.Skip
//	})
//}
//
//func AllowIfAdminMutationTenant() privacy.MutationRule {
//	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
//		v := viewer.FromContext(ctx)
//		if v == nil {
//			log.Println("No viewer in context")
//			return privacy.Denyf("viewer is not authenticated")
//		}
//		if v.RoleName == "admin" || v.RoleName == "superadmin" {
//			log.Println("Allowing mutation for admin or superadmin")
//			if tenantMutation, ok := m.(*ent.TenantMutation); ok {
//				tenantMutation.Where(tenant.IDEQ(v.TenantID))
//			}
//			return privacy.Allow
//		}
//		return privacy.Skip
//	})
//}

func AllowIfAdminQueryTenant() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			return privacy.Denyf("viewer is not authenticated")
		}
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			if tenantQuery, ok := q.(*ent.TenantQuery); ok {
				tenantQuery.Where(tenant.IDEQ(v.TenantID)) // 关键：限制查询自己的 Tenant
			}
			return privacy.Allow
		}
		return privacy.Skip
	})
}

func AllowIfAdminMutationTenant() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			return privacy.Denyf("viewer is not authenticated")
		}
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			if tenantMutation, ok := m.(*ent.TenantMutation); ok {
				tenantMutation.Where(tenant.IDEQ(v.TenantID)) // 关键：限制修改自己的 Tenant
			}
			return privacy.Allow
		}
		return privacy.Skip
	})
}
