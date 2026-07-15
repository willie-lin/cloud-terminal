package rule

import (
	"context"
	"fmt"

	"entgo.io/ent/privacy"
)

func AllowIfSuperAdminQueryUser() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		fmt.Println("AllowIfSuperAdminQueryUser: allowing")
		return privacy.Skip
	})
}

func AllowIfOwnerQueryUser() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		fmt.Println("AllowIfOwnerQueryUser: allowing")
		return privacy.Skip
	})
}

func AllowIfAdminQueryUser() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		fmt.Println("AllowIfAdminQueryUser: allowing")
		return privacy.Skip
	})
}

func AllowIfSuperAdminMutationUser() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		fmt.Println("AllowIfSuperAdminMutationUser: allowing")
		return privacy.Skip
	})
}

func AllowIfAdminMutationUser() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		fmt.Println("AllowIfAdminMutationUser: allowing")
		return privacy.Skip
	})
}

func AllowIfOwnerMutationUser() privacy.QueryMutationRule {
	return privacy.ContextQueryMutationRule(func(ctx context.Context) error {
		fmt.Println("AllowIfOwnerMutationUser: allowing")
		return privacy.Skip
	})
}
