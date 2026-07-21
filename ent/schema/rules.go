package schema

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/privacy"
	"github.com/willie-lin/cloud-terminal/viewer"
)

// AllowIfSuperAdmin 超级管理员放行所有操作
func AllowIfSuperAdmin() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		v := viewer.FromContext(ctx)
		if v != nil && v.RoleName == "superadmin" {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AllowIfAdmin 租户管理员及以上放行
func AllowIfAdmin() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		v := viewer.FromContext(ctx)
		if v != nil && (v.RoleName == "admin" || v.RoleName == "superadmin") {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// DenyIfNoViewer 拒绝未经认证的请求
func DenyIfNoViewer() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			return privacy.Denyf("authentication required")
		}
		return privacy.Skip
	})
}

// DenyMutationUnlessAdmin 非管理员禁止写操作
func DenyMutationUnlessAdmin() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		v := viewer.FromContext(ctx)
		if v == nil {
			return privacy.Denyf("authentication required")
		}
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			return privacy.Allow
		}
		return privacy.Denyf("only admins can perform this operation")
	})
}
