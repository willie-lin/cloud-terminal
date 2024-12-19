package api

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/permission"
	"github.com/willie-lin/cloud-terminal/app/database/ent/privacy"
	"github.com/willie-lin/cloud-terminal/app/database/ent/role"
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
			err = client.Role.UpdateOne(role).AddPermissions(permissionsToAdd...).Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to associate permissions with role %s: %w", role.Name, err)
			}
		}
	}
	return nil
}
