package api

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/accesspolicy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/account"
	"github.com/willie-lin/cloud-terminal/app/database/ent/platform"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/role"
	"github.com/willie-lin/cloud-terminal/app/database/ent/schema"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"os"
)

// {"tenant:create", "允许创建租户"},
//		{"tenant:read", "允许读取租户信息"},
//		{"tenant:update", "允许更新租户信息"},
//		{"tenant:delete", "允许删除租户"},
//		{"tenant:list", "允许列出所有租户"},
//
//		// 账户管理权限 (例如，如果需要区分账户级别的权限)
//		{"account:create", "允许创建账户信息"},
//		{"account:read", "允许读取账户信息"},
//		{"account:update", "允许更新账户信息"},
//		{"account:delete", "允许删除账户信息"},
//		{"account:list", "允许列出所有账户信息"},
//
//		// 用户管理权限
//		{"user:create", "允许创建用户"},
//		{"user:read", "允许读取用户信息"},
//		{"user:update", "允许更新用户信息"},
//		{"user:delete", "允许删除用户"},
//		{"user:list", "允许列出所有用户"},
//
//		// 角色管理权限
//		{"role:create", "允许创建角色"},
//		{"role:read", "允许读取角色信息"},
//		{"role:update", "允许更新角色信息"},
//		{"role:delete", "允许删除角色"},
//		{"role:list", "允许列出所有角色"},
//		{"role:assign", "允许为用户分配角色"},
//		{"role:unassign", "允许取消用户分配角色"},
//
//
//		// 策略管理权限
//		{"policy:create", "允许创建访问策略"},
//		{"policy:read", "允许读取访问策略信息"},
//		{"policy:update", "允许更新访问策略信息"},
//		{"policy:delete", "允许删除访问策略"},
//		{"policy:list", "允许列出所有策略"},
//		{"policy:attach", "允许将策略附加到主体"}, // 新增
//		{"policy:detach", "允许将策略从主体分离"}, // 新增
//
//		// 资源管理权限 (例如，如果需要区分账户级别的权限)
//		{"resource:create", "允许创建资源信息"},
//		{"resource:read", "允许读取资源信息"},
//		{"resource:update", "允许更新资源信息"},
//		{"resource:delete", "允许删除资源信息"},
//		{"resource:list", "允许列出所有资源信息"},
//		{"resource:grant", "授予对资源的访问权限"},  // 新增，更通用
//		{"resource:revoke", "撤销对资源的访问权限"}, // 新增，更通用
//
//		// 自管理权限
//		{"user:self-update", "允许用户更新自己的信息"},
//
//		// 审计日志权限
//		{"audit:read", "允许读取审计日志"},
//	}

//
//func InitPermissions(client *ent.Client, account *ent.Account) error {
//	ctx := context.Background()
//
//	permissionsToCreate := []struct {
//		Name        string
//		Description string
//	}{
//		// 租户管理权限
//		{"tenant:create", "允许创建租户"},
//		{"tenant:read", "允许读取租户信息"},
//		{"tenant:update", "允许更新租户信息"},
//		{"tenant:delete", "允许删除租户"},
//		{"tenant:list", "允许列出所有租户"},
//
//		// 账户管理权限 (例如，如果需要区分账户级别的权限)
//		{"account:create", "允许创建账户信息"},
//		{"account:read", "允许读取账户信息"},
//		{"account:update", "允许更新账户信息"},
//		{"account:delete", "允许删除账户信息"},
//		{"account:list", "允许列出所有账户信息"},
//
//		// 用户管理权限
//		{"user:create", "允许创建用户"},
//		{"user:read", "允许读取用户信息"},
//		{"user:update", "允许更新用户信息"},
//		{"user:delete", "允许删除用户"},
//		{"user:list", "允许列出所有用户"},
//
//		// 角色管理权限
//		{"role:create", "允许创建角色"},
//		{"role:read", "允许读取角色信息"},
//		{"role:update", "允许更新角色信息"},
//		{"role:delete", "允许删除角色"},
//		{"role:list", "允许列出所有角色"},
//		{"role:assign", "允许为用户分配角色"},
//		{"role:unassign", "允许取消用户分配角色"},
//
//		// 权限管理权限
//		{"permission:create", "允许创建权限"},
//		{"permission:read", "允许读取权限信息"},
//		{"permission:update", "允许更新权限信息"},
//		{"permission:delete", "允许删除权限"},
//		{"permission:list", "允许列出所有权限"},
//
//		// 策略管理权限
//		{"policy:create", "允许创建访问策略"},
//		{"policy:read", "允许读取访问策略信息"},
//		{"policy:update", "允许更新访问策略信息"},
//		{"policy:delete", "允许删除访问策略"},
//		{"policy:list", "允许列出所有策略"},
//		{"policy:attach", "允许将策略附加到主体"}, // 新增
//		{"policy:detach", "允许将策略从主体分离"}, // 新增
//
//		// 资源管理权限 (例如，如果需要区分账户级别的权限)
//		{"resource:create", "允许创建资源信息"},
//		{"resource:read", "允许读取资源信息"},
//		{"resource:update", "允许更新资源信息"},
//		{"resource:delete", "允许删除资源信息"},
//		{"resource:list", "允许列出所有资源信息"},
//		{"resource:grant", "授予对资源的访问权限"},  // 新增，更通用
//		{"resource:revoke", "撤销对资源的访问权限"}, // 新增，更通用
//
//		// 自管理权限
//		{"user:self-update", "允许用户更新自己的信息"},
//
//		// 审计日志权限
//		{"audit:read", "允许读取审计日志"},
//	}
//
//	for _, p := range permissionsToCreate {
//		//existingPermission, err := client.Permission.Query().Where(permission.NameEQ(p.Name)).Where(permission.HasAccountWith(account.IDEQ(account.ID))).Only(ctx)
//		existingPermission, err := client.Permission.Query().
//			Where(
//				permission.NameEQ(p.Name),
//				//permission.HasAccountWith(account.IDEQ(system.ID)),
//			).
//			Only(ctx)
//
//		if err != nil && !ent.IsNotFound(err) {
//			return fmt.Errorf("query permission %s failed: %w", p.Name, err)
//		}
//
//		if ent.IsNotFound(err) {
//			_, err := client.Permission.Create().SetName(p.Name).SetDescription(p.Description).SetAccount(account).Save(ctx)
//			if err != nil {
//				return fmt.Errorf("create permission %s failed: %w", p.Name, err)
//			}
//			log.Printf("Created permission: %s for account: %v", p.Name, account.ID)
//		} else {
//			log.Printf("Permission: %s already exists for account: %v, ID: %v", p.Name, account.ID, existingPermission.ID)
//		}
//	}
//
//	return nil
//}

func InitSuperAdminAndSuperRoles(client *ent.Client) error {
	// 使用 privacy.DecisionContext 跳过隐私检查
	ctx := privacy.DecisionContext(context.Background(), privacy.Allow)
	platformName := utils.PlatformName // 或者其他你想要的默认平台名称

	// 0. 检查或创建 Platform
	defaultPlatform, err := client.Platform.Query().Where(platform.NameEQ(platformName)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return err
	}
	if ent.IsNotFound(err) {
		defaultPlatform, err = client.Platform.Create().
			SetName(platformName).
			Save(ctx)
		if err != nil {
			return err
		}
		log.Print("Created CloudSecPlatform platform")
	} else {
		log.Print("CloudSecPlatform platform already exists.")
	}

	// 1. 检查或创建 "management" 租户（用于管理）
	tenantName := utils.ManagementTenant
	managementTenant, err := client.Tenant.Query().Where(tenant.NameEQ(tenantName)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return err
	}
	if ent.IsNotFound(err) {
		managementTenant, err = client.Tenant.Create().
			SetName(tenantName).
			SetPlatform(defaultPlatform).
			Save(ctx)
		if err != nil {
			return err
		}
		log.Printf("Created management tenant for %s platform", platformName)
	} else {
		log.Printf("Management tenant already exists for %s platform.", platformName)
	}

	// 2. 检查或创建 "default" 租户（用于普通业务）
	defaultTenantName := "default"
	_, err = client.Tenant.Query().Where(tenant.NameEQ(defaultTenantName)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return err
	}
	if ent.IsNotFound(err) {
		_, err = client.Tenant.Create().
			SetName(defaultTenantName).
			SetPlatform(defaultPlatform).
			Save(ctx)
		if err != nil {
			return err
		}
		log.Print("Created default tenant for %s platform", platformName)
	} else {
		log.Print("Default tenant already exists for %s platform.", platformName)
	}

	// 3. 检查或创建 "system" Account，并关联到 "management" Tenant
	accountName := utils.Account
	sa, err := client.Account.
		Query().
		Where(account.NameEQ(accountName)).
		Where(account.HasTenantWith(tenant.IDEQ(managementTenant.ID))).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return err
	}
	if ent.IsNotFound(err) {
		sa, err = client.Account.Create().
			SetName(accountName).
			SetTenant(managementTenant).
			Save(ctx)
		if err != nil {
			return err
		}
		log.Printf("Created system account for %s platform", platformName)
	} else {
		log.Printf("System account already exists for %s platform.", platformName)
	}

	// 4. 创建平台角色，租户角色，和超级管理员账户
	rolesToCreate := []string{"super_admin", "platform_admin", "tenant_admin"}

	for _, roleName := range rolesToCreate {
		//_, err := client.Role.Query().Where(role.NameEQ(roleName)).Where(role.HasTenantWith(tenant.NameEQ(managementTenant.Name))).Only(ctx)
		r, err := client.Role.Query().Where(role.NameEQ(roleName)).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return fmt.Errorf("query role %s failed: %w", roleName, err)
		}
		if ent.IsNotFound(err) {
			r, err = client.Role.Create().SetName(roleName).Save(ctx)
			if err != nil {
				return fmt.Errorf("create role %s failed: %w", roleName, err)
			}
			_, err = sa.Update().AddRoles(r).Save(ctx)
			if err != nil {
				return fmt.Errorf("关联 system 账户和 role %s 角色失败: %w", err)
			}
			log.Printf("Created role: %s (ID: %v)", r.Name, r.ID)
		} else {
			log.Printf("Role: %s (ID: %v) already exists.", r.Name, r.ID)
		}
	}
	// 5. 创建超级管理员策略
	superAdminRole, err := client.Role.Query().Where(role.NameEQ("super_admin")).Only(ctx)
	if err != nil {
		return fmt.Errorf("query super admin role failed: %w", err)
	}

	superAdminPolicy, err := client.AccessPolicy.Query().Where(accesspolicy.NameEQ("super_admin_policy")).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("query super admin policy failed: %w", err)
	}

	if ent.IsNotFound(err) {
		statements := []schema.PolicyStatement{
			{
				Effect: "Allow",
				Actions: []string{
					utils.ActionConfigRead,
					utils.ActionConfigUpdateDatabase,
					utils.ActionConfigUpdateNetwork,
					// ActionTenantCreate Tenant Actions
					utils.ActionTenantCreate,
					utils.ActionTenantRead,
					utils.ActionTenantUpdate,
					utils.ActionTenantDelete,
					// ActionUserCreate User Actions
					utils.ActionUserCreate,
					utils.ActionUserRead,
					utils.ActionUserUpdate,
					utils.ActionUserDelete,

					// ActionRoleCreate ActionUserCreate User Actions
					utils.ActionRoleCreate,
					utils.ActionRoleRead,
					utils.ActionRoleUpdate,
					utils.ActionRoleDelete,

					// ActionProjectCreate Project Actions
					utils.ActionProjectCreate,
					utils.ActionProjectRead,
					utils.ActionProjectUpdate,
					utils.ActionProjectDelete,

					// ActionAuditLogRead Audit Log Actions
					utils.ActionAuditLogRead,
					utils.ActionAuditLogExport,
				}, // 超级管理员拥有所有操作权限
				Resources: []string{
					utils.ResourceUserAll, // 匹配所有租户和账户的用户
					utils.ResourceConfigAll,
					utils.ResourceTenantAll,
					utils.ResourceAccountAll,
					utils.ResourceRoleAll,
					utils.ResourcePolicyAll,
					utils.ResourceProjectAll,
					utils.ResourceAuditLogAll,
				}, // 超级管理员拥有所有资源权限
			},
		}

		superAdminPolicy, err = client.AccessPolicy.Create().
			SetName("super_admin_policy").
			SetStatements(statements).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("创建 super admin policy 失败: %w", err)
		}
		log.Printf("Created super admin policy: %s (ID: %v)", superAdminPolicy.Name, superAdminPolicy.ID)
	} else if err == nil {
		log.Printf("Super admin policy: %s (ID: %v) already exists.", superAdminPolicy.Name, superAdminPolicy.ID)
	} else {
		return fmt.Errorf("unexpected error when querying super admin policy: %w", err)
	}

	// 检查策略是否已经关联到角色
	exists, err := client.Role.Query().
		Where(role.IDEQ(superAdminRole.ID)).
		Where(role.HasAccessPoliciesWith(accesspolicy.IDEQ(superAdminPolicy.ID))).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("checking role policy existence: %w", err)
	}

	if !exists {
		_, err = superAdminRole.Update().AddAccessPolicies(superAdminPolicy).Save(ctx)
		if err != nil {
			return fmt.Errorf("关联 super_admin 角色和 super_admin 策略失败: %w", err)
		}
		log.Printf("关联 super_admin 角色和 super_admin 策略")
	} else {
		log.Printf("super_admin 角色和 super_admin 策略已经关联")
	}

	// 6. 创建超级管理员用户并关联到账户
	superAdminUsername := "SuperAdmin"
	_, err = client.User.Query().Where(user.UsernameEQ(superAdminUsername)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("query super admin user failed: %w", err)
	}
	if ent.IsNotFound(err) {
		initialPassword := os.Getenv("INITIAL_SUPERADMIN_PASSWORD")
		if initialPassword == "" {
			initialPassword = "67727a41b5b1d4dfca981e4045b1bb2f1e7fef0e3e8825c028949d186cad4c00"
			log.Printf("WARNING: INITIAL_SUPERADMIN_PASSWORD env not set, using default password, please change it ASAP for platform: %s!", platformName)
		}
		hashedPassword, err := utils.GenerateFromPassword([]byte(initialPassword), utils.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password failed: %w", err)
		}

		_, err = client.User.Create().
			SetUsername(superAdminUsername).
			SetPassword(string(hashedPassword)).
			SetEmail("superadmin@example.com").
			SetPhoneNumber("19288888888").
			SetAccount(sa).
			SetRole(superAdminRole).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("create super admin failed: %w", err)
		}
		log.Printf("Created initial Super Admin user: superadmin for platform: %s", platformName)

	} else {
		log.Printf("Initial Super Admin user already exists for platform: %s", platformName)
	}
	return nil
}
