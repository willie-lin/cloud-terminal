package api

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/permission"
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
	for _, p := range defaultPermissions {
		//log.Printf("Checking permission: %s - %s", p.Name, p.ResourceType)
		_, err := client.Permission.Query().
			Where(permission.NameEQ(p.Name)).
			Only(context.Background())
		if ent.IsNotFound(err) {
			//log.Printf("Permission not found, creating: %s", p.Name)
			_, err = client.Permission.Create().
				SetName(p.Name).
				SetDescription(p.Description).
				SetActions(p.Actions). // 直接设置动作列表
				SetResourceType(p.ResourceType).
				SetIsDisabled(false).
				Save(context.Background())
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
	// 获取所有权限
	allPermissions, err := client.Permission.Query().All(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch all permissions: %w", err)
	}
	log.Printf("Fetched %d permissions", len(allPermissions))

	for _, template := range defaultRoleTemplates {
		// 检查角色是否已经存在
		role, err := client.Role.Query().
			Where(role.NameEQ(template.Name)).
			Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("Creating role: %s", template.Name)
			newRole, err := client.Role.Create().
				SetName(template.Name).
				SetDescription(template.Description).
				//AddTenantIDs(tenantID).
				Save(context.Background())
			if err != nil {
				return fmt.Errorf("failed to create role %s: %w", template.Name, err)
			}
			log.Printf("Created role: %s", newRole.Name)
		} else if err != nil {
			return fmt.Errorf("failed to query role %s: %w", template.Name, err)
		}

		if template.Name == "Admin" {
			// 管理员角色拥有所有权限
			for _, perm := range allPermissions {
				log.Printf("Associating permission %s with role %s", perm.Name, role.Name)
				if err := client.Role.UpdateOne(role).AddPermissions(perm).Exec(context.Background()); err != nil {
					log.Printf("Error associating permission %s with role %s: %v", perm.Name, role.Name, err)
					return fmt.Errorf("failed to associate permission %s with role %s: %w", perm.Name, role.Name, err)
				}
				log.Printf("Associated permission %s with role %s", perm.Name, role.Name)
			}
		} else {
			// 其他角色拥有特定权限
			for _, permName := range template.Permissions {
				log.Printf("Querying permission: %s", permName)
				p, err := client.Permission.Query().
					Where(permission.NameEQ(permName)).
					Only(context.Background())
				if err != nil {
					log.Printf("Error querying permission %s: %v", permName, err)
					return fmt.Errorf("failed to query permission %s: %w", permName, err)
				}
				log.Printf("Associating permission %s with role %s", p.Name, role.Name)
				if err := client.Role.UpdateOne(role).AddPermissions(p).Exec(context.Background()); err != nil {
					return fmt.Errorf("failed to associate permission %s with role %s: %w", permName, role.Name, err)
				}
				log.Printf("Associated permission %s with role %s", p.Name, role.Name)
			}
		}
	}
	return nil
}
