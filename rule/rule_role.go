package rule

import (
	"context"
	"fmt"

	"entgo.io/ent/privacy"
)

func AllowIfAdminQueryRole() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		fmt.Println("AllowIfAdminQueryRole: allowing")
		return privacy.Skip
	})
}

func AllowIfSuperAdminMutationRole() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		fmt.Println("AllowIfSuperAdminMutationRole: allowing")
		return privacy.Skip
	})
}

func AllowIfAdminMutationRole() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		fmt.Println("AllowIfAdminMutationRole: allowing")
		return privacy.Skip
	})
}
