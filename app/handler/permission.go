package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/permission"
	"github.com/willie-lin/cloud-terminal/app/database/ent/role"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/app/viewer"
	"net/http"
)

type PermissionDTO struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	ResourceIDs []uuid.UUID `json:"resource_ids"` // 可选的资源ID
}

func CheckPermissionName(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		dto := new(PermissionDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding permission: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		// 检查permission是否已经存在
		exists, err := client.Permission.Query().Where(permission.NameEQ(dto.Name)).Exist(context.Background())
		if err != nil {
			log.Printf("Error checking name: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error checking name from database"})
		}
		return c.JSON(http.StatusOK, map[string]bool{"exists": exists})
	}
}

// GetAllPermissions  获取所有permission
func GetAllPermissions(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		permissions, err := client.Permission.Query().All(context.Background())
		if err != nil {
			log.Printf("Error querying permissions: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying permissions from database"})
		}
		return c.JSON(http.StatusOK, permissions)
	}
}

//
//func GetAllPermissionsByUserByTenant(client *ent.Client) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		// 从请求上下文中获取租户ID
//		v := viewer.FromContext(c.Request().Context())
//		if v == nil {
//			log.Printf("No viewer found in context")
//			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No viewer found in context"})
//		}
//		userID := v.UserID
//		tenantID := v.TenantID
//
//		// 查询用户的所有角色
//		roles, err := client.User.Query().
//			Where(user.IDEQ(userID)).
//			QueryRoles().
//			WithPermissions().
//			//Where(role.HasTenantWith(tenant.IDEQ(tenantID))).
//			All(c.Request().Context())
//		if err != nil {
//			log.Printf("Error querying roles for user %s in tenant %s: %v", userID, tenantID, err)
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying roles from database"})
//		}
//
//		// 收集所有权限
//		var permissions []*ent.Permission
//		permissionSet := make(map[string]bool)
//		for _, role := range roles {
//			for _, perm := range role.Edges.Permissions {
//				if !permissionSet[perm.ID.String()] {
//					permissions = append(permissions, perm)
//					permissionSet[perm.ID.String()] = true
//				}
//			}
//		}
//
//		log.Printf("Permissions for user %s in tenant %s: %v", userID, tenantID, permissions)
//		return c.JSON(http.StatusOK, permissions)
//	}
//}

func GetAllPermissionsByTenant(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			log.Printf("No viewer found in context")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No viewer found in context"})
		}
		tenantID := v.TenantID
		roleName := v.RoleName

		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, tenantID, roleName)

		var permissions []*ent.Permission
		var err error
		if roleName == "admin" || roleName == "superadmin" {
			permissions, err = client.Permission.Query().Where(permission.HasRolesWith(role.HasTenantWith(tenant.IDEQ(tenantID)))).All(c.Request().Context())
			if err != nil {
				log.Printf("Error querying permissions for tenant %s: %v", tenantID, err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying permissions from database"})
			}
		} else {
			permissions, err = client.Permission.Query().Where(permission.HasRolesWith(role.HasUsersWith(user.IDEQ(v.UserID)))).All(c.Request().Context())
			if err != nil {
				log.Printf("Error querying permissions for user %s: %v", v.UserID, err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying permissions from database"})
			}
		}

		log.Printf("Permissions for tenant %s: %v", tenantID, permissions)
		return c.JSON(http.StatusOK, permissions)
	}
}

// CreatePermission  创建permission
func CreatePermission(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		v := viewer.FromContext(c.Request().Context())
		if v == nil || v.RoleName != "admin" {
			log.Printf("No viewer found in context or not authorized")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		var permissions []*PermissionDTO
		if err := c.Bind(&permissions); err != nil {
			log.Printf("Error binding permission: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		createdPermissions := make([]*ent.Permission, 0, len(permissions))

		for _, dto := range permissions {
			p, err := client.Permission.Create().
				SetName(dto.Name).
				SetDescription(dto.Description).
				//SetTenantID(v.TenantID).            // 关联到当前租户
				AddResourceIDs(dto.ResourceIDs...). // 可选的资源ID
				Save(context.Background())
			if err != nil {
				log.Printf("Error creating permission: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create permission"})
			}
			createdPermissions = append(createdPermissions, p)
		}
		return c.JSON(http.StatusCreated, map[string]string{"message": "Permission created successfully"})
	}
}

func DeletePermissionByName(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		dto := new(PermissionDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding permission: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		p, err := client.Permission.Query().Where(permission.NameEQ(dto.Name)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("Permission not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Permission not found"})
		}
		if err != nil {
			log.Printf("Error querying permission: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying role from database"})
		}

		err = client.Permission.DeleteOne(p).Exec(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("Permission not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Permission not found"})
		}
		if err != nil {
			log.Printf("Error deleting permission: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error deleting permission"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}
