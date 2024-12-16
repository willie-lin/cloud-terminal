package rule

import (
	"context"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/app/viewer"
)

// UserPrivacyPolicy 定义用户相关的隐私策略
var UserPrivacyPolicy = privacy.QueryPolicy{
	DenyIfNotTenant,
	//DenyIfNotAdminRule,
}

// DenyIfNotTenant 是一个示例规则，拒绝不属于当前租户的用户访问
var DenyIfNotTenant privacy.QueryRule = privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
	viewer := viewer.FromContext(ctx)
	if viewer == nil {
		return privacy.Denyf("Viewer not found in context")
	}
	tenantID := viewer.TenantID
	query, ok := q.(*ent.UserQuery)
	if !ok {
		return privacy.Denyf("Unexpected query type")
	}
	query.Where(user.HasTenantWith(tenant.IDEQ(tenantID)))
	return privacy.Skip

})

// DenyIfNotAdminRule 是一个示例规则，拒绝非管理员用户的访问
//var DenyIfNotAdminRule privacy.QueryRule = privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
//	viewer := viewer.FromContext(ctx)
//	if viewer == nil {
//		return privacy.Denyf("Viewer not found in context")
//	}
//	if viewer.RoleID != "admin" {
//		return privacy.Denyf("Access denied for non-admin user")
//	}
//	return privacy.Skip
//})
