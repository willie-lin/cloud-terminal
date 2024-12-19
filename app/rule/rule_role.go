package rule

import (
	"context"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/viewer"
	"log"
)

// AllowOnlySuperAdminMutationRole 是一个隐私规则示例，仅允许超级管理员对角色进行变更操作。
func AllowOnlySuperAdminMutationRole() privacy.MutationRule {
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

// AllowOnlySuperAdminQueryRole 是一个隐私规则示例，仅允许超级管理员对角色进行查询操作。
func AllowOnlySuperAdminQueryRole() privacy.QueryRule {
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

// AllowIfAdminOrSuperAdminMutationRole 是一个隐私规则示例，允许管理员或超级管理员对角色进行变更操作。
func AllowIfAdminOrSuperAdminMutationRole() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)
		if v.RoleName != "admin" && v.RoleName != "superadmin" {
			log.Println("Denying mutation for non-admin or non-superadmin")
			return privacy.Denyf("only admin or superadmin can perform this action")
		}
		return privacy.Allow
	})
}

// AllowIfAdminOrSuperAdminQueryRole 是一个隐私规则示例，允许管理员或超级管理员对角色进行查询操作。
func AllowIfAdminOrSuperAdminQueryRole() privacy.QueryRule {
	return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			log.Println("No viewer in context")
			return privacy.Denyf("viewer is not authenticated")
		}
		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, v.TenantID, v.RoleName)
		if v.RoleName != "admin" && v.RoleName != "superadmin" {
			log.Println("Denying query for non-admin or non-superadmin")
			return privacy.Denyf("only admin or superadmin can perform this action")
		}
		return privacy.Allow
	})
}
