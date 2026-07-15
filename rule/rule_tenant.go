package rule

import (
	"context"
	"fmt"

	"entgo.io/ent/privacy"
	"github.com/willie-lin/cloud-terminal/viewer"
)

func AllowIfAdminQueryTenant() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		v := viewer.FromContext(ctx)
		fmt.Println("------------Tenant----------")
		if v == nil {
			return privacy.Denyf("viewer is not authenticated")
		}
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

func AllowIfAdminMutationTenant() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		v := viewer.FromContext(ctx)
		fmt.Println("------------Tenant----------")
		if v == nil {
			return privacy.Denyf("viewer is not authenticated")
		}
		if v.RoleName == "admin" || v.RoleName == "superadmin" {
			return privacy.Allow
		}
		return privacy.Skip
	})
}
