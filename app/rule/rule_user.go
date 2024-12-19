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

// AllowIfSuperAdminQueryUser  允许超级用户查询所有资源。
func AllowIfSuperAdminQueryUser() privacy.QueryRule {
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

// AllowIfSuperAdminMutationUser  允许超级用户变更所有资源。
func AllowIfSuperAdminMutationUser() privacy.MutationRule {
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

// AllowIfAdminQueryUser 允许管理员查询其租户下的用户。
func AllowIfAdminQueryUser() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			log.Println("Allowing query for admin or superadmin")
			// 确保查询限于管理员的租户
			if userQuery, ok := q.(*ent.UserQuery); ok {
				userQuery.Where(user.HasTenantWith(tenant.IDEQ(v.TenantID)))
			}
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AllowIfAdminMutationUser 允许管理员变更其租户下的用户。
func AllowIfAdminMutationUser() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			log.Println("Allowing mutation for admin or superadmin")
			// 确保变更限于管理员的租户
			if userMutation, ok := m.(*ent.UserMutation); ok {
				userMutation.Where(user.HasTenantWith(tenant.IDEQ(v.TenantID)))
			}
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AllowIfOwnerQueryUser  允许用户查询自己的个人信息。
func AllowIfOwnerQueryUser() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)

		if userQuery, ok := q.(*ent.UserQuery); ok {
			userQuery.Where(user.IDEQ(v.UserID))
			log.Println("Allowing query for owner")
			return privacy.Allow
		}

		log.Println("Denying query for non-owner")
		return privacy.Denyf("only owner can perform this action")
	})
}

// AllowIfOwnerMutationUser  允许用户变更自己的个人信息。
func AllowIfOwnerMutationUser() privacy.MutationRule {
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
