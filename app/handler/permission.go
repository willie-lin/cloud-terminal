package handler

//
//import (
//	"fmt"
//	"github.com/labstack/echo/v4"
//	"github.com/labstack/gommon/log"
//	"github.com/willie-lin/cloud-terminal/app/database/ent"
//	"github.com/willie-lin/cloud-terminal/app/database/ent/permission"
//	"github.com/willie-lin/cloud-terminal/app/database/ent/role"
//	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
//	"github.com/willie-lin/cloud-terminal/app/viewer"
//	"net/http"
//)
//
//type PermissionDTO struct {
//	Name         string   `json:"name"`
//	Description  string   `json:"description"`
//	Actions      []string `json:"actions"`
//	ResourceType string   `json:"resource_type"` // 可选的资源ID
//}
//
//func CheckPermissionName(client *ent.Client) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		v := viewer.FromContext(c.Request().Context())
//		if v == nil {
//			log.Printf("No viewer found in context")
//			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No viewer found in context"})
//		}
//
//		dto := new(PermissionDTO)
//		if err := c.Bind(&dto); err != nil {
//			log.Printf("Error binding permission: %v", err)
//			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
//		}
//		// 检查permission是否已经存在
//		exists, err := client.Permission.Query().Where(permission.NameEQ(dto.Name)).Exist(c.Request().Context())
//		if err != nil {
//			log.Printf("Error checking name: %v", err)
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error checking name from database"})
//		}
//		return c.JSON(http.StatusOK, map[string]bool{"exists": exists})
//	}
//}
//
//// GetAllPermissions  获取所有permission
//func GetAllPermissions(client *ent.Client) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		permissions, err := client.Permission.Query().All(c.Request().Context())
//		if err != nil {
//			log.Printf("Error querying permissions: %v", err)
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying permissions from database"})
//		}
//		return c.JSON(http.StatusOK, permissions)
//	}
//}
//
//func GetAllPermissionsByTenant(client *ent.Client) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		v := viewer.FromContext(c.Request().Context())
//		if v == nil {
//			log.Printf("No viewer found in context")
//			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No viewer found in context"})
//		}
//		tenantID := v.TenantID
//		roleName := v.RoleName
//
//		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, tenantID, roleName)
//
//		var permissions []*ent.Permission
//		var err error
//		if roleName == "admin" || roleName == "superadmin" {
//			permissions, err = client.Permission.Query().
//				//Where(permission.HasRolesWith(role.HasTenantWith(tenant.IDEQ(tenantID)))).
//				All(c.Request().Context())
//			if err != nil {
//				log.Printf("Error querying permissions for tenant %s: %v", tenantID, err)
//				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying permissions from database"})
//			}
//		} else {
//			permissions, err = client.Permission.Query().
//				//Where(permission.HasRolesWith(role.HasUsersWith(user.IDEQ(v.UserID)))).
//				All(c.Request().Context())
//			if err != nil {
//				log.Printf("Error querying permissions for user %s: %v", v.UserID, err)
//				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying permissions from database"})
//			}
//		}
//
//		log.Printf("Permissions for tenant %s: %v", tenantID, permissions)
//		return c.JSON(http.StatusOK, permissions)
//	}
//}
//
//// CreatePermission  创建permission
//func CreatePermission(client *ent.Client) echo.HandlerFunc {
//	return func(c echo.Context) (err error) {
//		v := viewer.FromContext(c.Request().Context())
//		if v == nil || v.RoleName != "admin" {
//			log.Printf("No viewer found in context or not authorized")
//			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
//		}
//
//		var permissions []*PermissionDTO
//		if err := c.Bind(&permissions); err != nil {
//			log.Printf("Error binding permission: %v", err)
//			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
//		}
//		// 查询当前租户下的角色，通过边关系表进行关联
//		role, err := client.Role.Query().Where(role.HasUsersWith(user.IDEQ(v.UserID)), role.NameEQ(v.RoleName)).Only(c.Request().Context())
//		if err != nil {
//			log.Printf("Error querying role for tenant %s: %v", v.TenantID, err)
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying role from database"})
//		}
//
//		fmt.Println(role)
//		createdPermissions := make([]*ent.Permission, 0, len(permissions))
//
//		for _, dto := range permissions {
//			p, err := client.Permission.Create().
//				SetName(dto.Name).
//				SetDescription(dto.Description).
//				//SetTenantID(v.TenantID).            // 关联到当前租户
//				SetActions(dto.Actions). // 可选的资源ID
//				SetResourceType(dto.ResourceType).
//				//AddRoleIDs(role.ID).
//				Save(c.Request().Context())
//			if err != nil {
//				log.Printf("Error creating permission: %v", err)
//				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create permission"})
//			}
//			createdPermissions = append(createdPermissions, p)
//		}
//		return c.JSON(http.StatusCreated, map[string]string{"message": "Permission created successfully"})
//	}
//}
//
//// IsDefaultPermission 检查角色是否为默认角色
//func IsDefaultPermission(permissionName string) bool {
//	defaultPermissions := []string{"UserManagement", "RoleManagement", "LogAccess", "ApplicationDevelopment"}
//	for _, r := range defaultPermissions {
//		if r == permissionName {
//			return true
//		}
//	}
//	return false
//}
//
//func DeletePermissionByName(client *ent.Client) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		// 从请求上下文中获取租户ID
//		v := viewer.FromContext(c.Request().Context())
//		if v == nil || v.RoleName != "admin" {
//			log.Printf("No viewer found in context or not authorized")
//			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
//		}
//
//		dto := new(PermissionDTO)
//		if err := c.Bind(&dto); err != nil {
//			log.Printf("Error binding permission: %v", err)
//			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
//		}
//
//		if IsDefaultPermission(dto.Name) {
//			log.Printf("Attempt to delete default role: %s", dto.Name)
//			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot delete default Permission"})
//		}
//
//		p, err := client.Permission.Query().Where(permission.NameEQ(dto.Name)).Only(c.Request().Context())
//		if ent.IsNotFound(err) {
//			log.Printf("Permission not found: %v", err)
//			return c.JSON(http.StatusNotFound, map[string]string{"error": "Permission not found"})
//		}
//		if err != nil {
//			log.Printf("Error querying permission: %v", err)
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying role from database"})
//		}
//
//		err = client.Permission.DeleteOne(p).Exec(c.Request().Context())
//		if ent.IsNotFound(err) {
//			log.Printf("Permission not found: %v", err)
//			return c.JSON(http.StatusNotFound, map[string]string{"error": "Permission not found"})
//		}
//		if err != nil {
//			log.Printf("Error deleting permission: %v", err)
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error deleting permission"})
//		}
//		return c.NoContent(http.StatusNoContent)
//	}
//}
