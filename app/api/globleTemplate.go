package api

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/account"
	"github.com/willie-lin/cloud-terminal/app/database/ent/permission"
	"github.com/willie-lin/cloud-terminal/app/database/ent/platform"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/role"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"os"
)

type PermissionTemplate struct {
	Name         string
	Description  string
	Actions      []string
	ResourceType string
}

type RoleTemplate struct {
	Name        string
	Description string
	Permissions []string
}

var defaultPermissions = []PermissionTemplate{
	{"UserManagement", "管理所有用户", []string{"create", "delete", "update", "read"}, "user"},
	{"RoleManagement", "管理所有角色", []string{"create", "delete", "update", "read"}, "role"},
	{"LogAccess", "访问系统日志", []string{"read"}, "log"},
	{"ApplicationDevelopment", "开发和维护应用", []string{"read", "write", "deploy"}, "application"},
	{"AuditLogAccess", "审计系统日志", []string{"read"}, "audit"},
	{"PersonalDataAccess", "访问和编辑个人数据", []string{"read", "update"}, "personal_data"},
}

var defaultRoleTemplates = []RoleTemplate{
	{"Admin", "管理员角色", []string{}}, // 管理员拥有所有权限
	{"Developer", "开发者角色", []string{"ApplicationDevelopment", "LogAccess"}},
	{"Auditor", "审计员角色", []string{"AuditLogAccess"}},
	{"User", "普通用户角色", []string{"PersonalDataAccess"}},
}

func InitializeGlobalPermissions(client *ent.Client) error {
	// 使用 privacy.DecisionContext 跳过隐私检查
	ctx := privacy.DecisionContext(context.Background(), privacy.Allow)
	for _, p := range defaultPermissions {
		//log.Printf("Checking permission: %s - %s", p.Name, p.ResourceType)
		_, err := client.Permission.Query().
			Where(permission.NameEQ(p.Name)).
			Only(ctx)
		if ent.IsNotFound(err) {
			//log.Printf("Permission not found, creating: %s", p.Name)
			_, err = client.Permission.Create().
				SetName(p.Name).
				SetDescription(p.Description).
				SetActions(p.Actions). // 直接设置动作列表
				SetResourceType(p.ResourceType).
				SetIsDisabled(false).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to create permission %s: %w", p.Name, err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to query permission %s: %w", p.Name, err)
		}
	}
	log.Printf("Global permissions initialized successfully")
	return nil
}

func InitializeTenantRolesAndPermissions(client *ent.Client) error {
	// 使用 privacy.DecisionContext 跳过隐私检查
	ctx := privacy.DecisionContext(context.Background(), privacy.Allow)
	allPermissions, err := client.Permission.Query().All(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch all permissions: %w", err)
	}
	permissionMap := make(map[string]*ent.Permission, len(allPermissions))
	for _, perm := range allPermissions {
		permissionMap[perm.Name] = perm
	}

	for _, template := range defaultRoleTemplates {
		role, err := client.Role.Query().Where(role.NameEQ(template.Name)).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return fmt.Errorf("failed to query role %s: %w", template.Name, err)
		}

		if ent.IsNotFound(err) {
			role, err = client.Role.Create().
				SetName(template.Name).
				SetDescription(template.Description).
				SetIsDisabled(false).
				SetIsDefault(true).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to create role %s: %w", template.Name, err)
			}
			log.Printf("Created role: %s", role.Name)
		} else {
			log.Printf("Role %s already exists", template.Name)
		}

		var permissionNames []string

		if template.Name == "Admin" {
			permissionNames = make([]string, 0, len(permissionMap))
			for k := range permissionMap {
				permissionNames = append(permissionNames, k)
			}
		} else {
			permissionNames = template.Permissions
		}
		var permissionsToAdd []*ent.Permission
		for _, permName := range permissionNames {
			perm, ok := permissionMap[permName]
			if !ok {
				log.Printf("Permission %s not found for role %s", permName, template.Name)
				continue
			}
			permissionsToAdd = append(permissionsToAdd, perm)
		}
		if len(permissionsToAdd) > 0 {
			log.Printf("Adding %d permissions to role %s", len(permissionsToAdd), role.Name)
			err = client.Role.UpdateOne(role).Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to associate permissions with role %s: %w", role.Name, err)
			}
		}
	}
	return nil
}

func InitSuperAdminAndSuperRoles(client *ent.Client) error {
	// 使用 privacy.DecisionContext 跳过隐私检查
	ctx := privacy.DecisionContext(context.Background(), privacy.Allow)

	// 0. 检查或创建 Platform
	defaultPlatform, err := client.Platform.Query().Where(platform.NameEQ("CloudSecPlatform")).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return err
	}
	if ent.IsNotFound(err) {
		defaultPlatform, err = client.Platform.Create().SetName("CloudSecPlatform").Save(ctx)
		if err != nil {
			return err
		}
		log.Print("Created CloudSecPlatform platform")
	} else {
		log.Print("CloudSecPlatform platform already exists.")
	}

	// 1. 检查或创建 "management" 租户（用于管理）
	managementTenant, err := client.Tenant.Query().Where(tenant.NameEQ("management")).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return err
	}
	if ent.IsNotFound(err) {
		managementTenant, err = client.Tenant.Create().
			SetName("management").
			SetPlatform(defaultPlatform).
			Save(ctx)
		if err != nil {
			return err
		}
		log.Print("Created management tenant")
	} else {
		log.Print("Management tenant already exists.")
	}

	// 2. 检查或创建 "default" 租户（用于普通业务）
	_, err = client.Tenant.Query().Where(tenant.NameEQ("default")).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return err
	}
	if ent.IsNotFound(err) {
		_, err = client.Tenant.Create().
			SetName("default").
			SetPlatform(defaultPlatform).
			Save(ctx)
		if err != nil {
			return err
		}
		log.Print("Created default tenant")
	} else {
		log.Print("Default tenant already exists.")
	}

	// 3. 检查或创建 "system" Account，并关联到 "management" Tenant
	systemAccount, err := client.Account.Query().Where(account.NameEQ("system")).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return err
	}
	if ent.IsNotFound(err) {
		systemAccount, err = client.Account.Create().SetName("system").SetTenant(managementTenant).Save(ctx)
		if err != nil {
			return err
		}
		log.Print("Created system account")
	} else {
		log.Print("System account already exists.")
	}

	roles := []string{"super_admin", "platform_admin", "tenant_admin"}

	for _, roleName := range roles {
		_, err := client.Role.Query().Where(role.NameEQ(roleName)).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return err // 其他错误，直接返回
		}
		if ent.IsNotFound(err) {
			// 角色不存在，创建它
			_, err = client.Role.Create().
				SetName(roleName).
				SetAccount(systemAccount).
				Save(ctx)
			if err != nil {
				return err
			}
			log.Printf("Created role: %s", roleName)
		} else {
			log.Printf("Role already exist: %s", roleName)
		}

	}

	// 创建初始 Super Admin 用户
	superAdminUser, err := client.User.Query().Where(user.UsernameEQ("SuperAdmin")).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return err
	}

	if ent.IsNotFound(err) {
		initialPassword := os.Getenv("INITIAL_SUPERADMIN_PASSWORD")
		if initialPassword == "" {
			initialPassword = "67727a41b5b1d4dfca981e4045b1bb2f1e7fef0e3e8825c028949d186cad4c00" //设置一个默认密码，但是一定要提示用户修改
			log.Print("WARNING: INITIAL_SUPERADMIN_PASSWORD env not set, use default password 'changeme', please change it ASAP!")

		}
		//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(initialPassword), bcrypt.DefaultCost)
		fmt.Println(initialPassword)
		hashedPassword, err := utils.GenerateFromPassword([]byte(initialPassword), utils.DefaultCost)
		fmt.Println(hashedPassword)

		if err != nil {
			return err
		}
		superAdminUser, err = client.User.Create().
			SetUsername("SuperAdmin").
			SetPassword(string(hashedPassword)).
			SetEmail("superadmin@example.com"). // 设置一个默认邮箱
			SetPhoneNumber("19288888888").
			SetAccount(systemAccount).
			Save(ctx)
		if err != nil {
			return err
		}
		log.Print("Created initial Super Admin user: superadmin")

		superAdminRole, err := client.Role.Query().Where(role.NameEQ("super_admin")).Only(ctx)
		if err != nil {
			return err
		}
		err = superAdminUser.Update().AddRoles(superAdminRole).Exec(ctx)
		if err != nil {
			return err
		}
	} else {
		log.Print("Initial Super Admin user already exists.")
	}

	return nil
}
