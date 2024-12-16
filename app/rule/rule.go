package rule

import (
	"context"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/viewer"
)

// DenyIfNoViewer is a rule that returns Deny decision if the viewer is missing in the context.
func DenyIfNoViewer() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		view := viewer.FromContext(ctx)
		if view == nil {
			return privacy.Denyf("viewer-context is missing")
		}
		return privacy.Skip
	})
}

// AllowIfAdmin is a rule that returns Allow decision if the viewer is admin.
//func AllowIfAdmin() privacy.QueryMutationRule {
//	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
//		view := viewer.FromContext(ctx)
//		if view.Admin {
//			return privacy.Allow
//		}
//		return privacy.Skip
//	})
//}

//// AllowIfRole is a rule that returns Allow decision if the viewer has a specific role.
//func AllowIfRole(roleName string) privacy.QueryMutationRule {
//	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
//		view := viewer.FromContext(ctx)
//		if view.HasRole(roleName) {
//			return privacy.Allow
//		}
//		return privacy.Skip
//	})
//}

// AllowIfOwner is a rule that returns Allow decision if the viewer is the owner of the entity.
func AllowIfOwner() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		view := viewer.FromContext(ctx)
		entityOwnerID := ctx.Value("ENTITY_OWNER_ID") // 假设我们在上下文中传递了实体的所有者ID
		if entityOwnerID == view.UserID {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AllowIfTenantMember is a rule that allows access if the viewer belongs to the same tenant as the entity.
func AllowIfTenantMember() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		view := viewer.FromContext(ctx)
		tenantID := ctx.Value("ENTITY_TENANT_ID") // 假设我们在上下文中传递了实体的租户ID
		if tenantID == view.TenantID {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

//// DenyIfNotTenant DenyIfNotTenant 检查当前用户是否属于查询的租户
//var DenyIfNotTenant privacy.QueryRule = privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
//	viewer := viewer.FromContext(ctx)
//	if viewer == nil {
//		return privacy.Denyf("Viewer not found in context")
//	}
//	tenantID := viewer.TenantID
//	q.Where(entql.FieldEQ("tenant_id", tenantID))
//	return privacy.Skip
//})
//
//// TenantPolicy 为租户查询设置隐私策略
//var TenantPolicy = privacy.QueryPolicy{
//	DenyIfNotTenant,
//}
