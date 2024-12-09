package migrations

import (
	"context"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"log"
)

func AssignDefaultTenant(client *ent.Client) error {
	ctx := context.Background()

	// 创建默认租户
	defaultTenant, err := client.Tenant.Create().
		SetName("Default Tenant").
		SetDescription("Default tenant for existing users").
		Save(ctx)
	if err != nil {
		return err
	}

	// 获取所有没有租户的用户
	users, err := client.User.Query().Where(user.Not(user.HasTenant())).All(ctx)
	if err != nil {
		return err
	}

	// 为每个用户分配默认租户
	for _, u := range users {
		_, err := client.User.UpdateOne(u).SetTenant(defaultTenant).Save(ctx)
		if err != nil {
			log.Printf("Failed to assign default tenant to user %s: %v", u.Username, err)
		}
	}

	return nil
}
