package rule

import (
	"context"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/role"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/app/viewer"
	"log"
)

// AllowIfSuperAdminQueryRole 允许超级用户查询所有资源。
func AllowIfSuperAdminQueryRole() privacy.QueryRule {
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

// AllowIfSuperAdminMutationRole  允许超级用户变更所有资源。
func AllowIfSuperAdminMutationRole() privacy.MutationRule {
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

// AllowIfAdminQueryRole 允许管理员查询其租户下的角色。
func AllowIfAdminQueryRole() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			log.Println("Allowing query for admin or superadmin")
			// 确保查询限于管理员的租户
			if roleQuery, ok := q.(*ent.RoleQuery); ok {
				roleQuery.Where(role.HasTenantWith(tenant.IDEQ(v.TenantID)))
				//roleQuery.Where(role.HasUsersWith(user.HasTenantWith(tenant.IDEQ(v.TenantID))))
			}
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AllowIfAdminMutationRole 允许管理员变更其租户下的角色。
func AllowIfAdminMutationRole() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			log.Println("Allowing mutation for admin or superadmin")
			// 确保变更限于管理员的租户
			if roleMutation, ok := m.(*ent.RoleMutation); ok {
				roleMutation.Where(role.HasTenantWith(tenant.IDEQ(v.TenantID)))
				//roleMutation.Where(role.HasUsersWith(user.HasTenantWith(tenant.IDEQ(v.TenantID))))
			}
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AllowIfOwnerQueryRole 允许用户查询自己的角色。
func AllowIfOwnerQueryRole() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)
		if roleQuery, ok := q.(*ent.RoleQuery); ok {
			roleQuery.Where(role.HasUsersWith(user.IDEQ(v.UserID)))
			log.Println("Allowing query for user roles")
			return privacy.Allow
		}
		log.Println("Denying query for non-owner")
		return privacy.Denyf("only owner can perform this action")
	})
}
