package rule

import (
	"context"
	"fmt"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"github.com/willie-lin/cloud-terminal/app/viewer"
)

func AllowIfAdminQueryTenant() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		fmt.Println("------------Tenant----------")
		fmt.Println(v.UserID)
		fmt.Println(v.TenantID)
		fmt.Println(v.RoleName)
		fmt.Println("------------Tenant----------")
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
		fmt.Println("------------Tenant----------")
		fmt.Println(v.UserID)
		fmt.Println(v.TenantID)
		fmt.Println(v.RoleName)
		fmt.Println("------------Tenant----------")
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
