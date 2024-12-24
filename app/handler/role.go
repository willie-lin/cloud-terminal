package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/role"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"github.com/willie-lin/cloud-terminal/app/database/ent/user"
	"github.com/willie-lin/cloud-terminal/app/viewer"
	"net/http"
)

type RoleDTO struct {
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
}

func CheckRoleName(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 从请求上下文中获取 viewer 信息
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			log.Printf("No viewer found in context")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No viewer found in context"})
		}

		dto := new(RoleDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding role: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		// 在 viewer 上下文中执行查询
		exists, err := client.Role.Query().Where(role.NameEQ(dto.Name)).Exist(c.Request().Context())
		if err != nil {
			log.Printf("Error checking name: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error checking name from database"})
		}
		return c.JSON(http.StatusOK, map[string]bool{"exists": exists})
	}
}

// GetAllRoles 获取所有role
func GetAllRoles(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		roles, err := client.Role.Query().All(c.Request().Context())
		if err != nil {
			log.Printf("Error querying roles: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying roles from database"})
		}
		return c.JSON(http.StatusOK, roles)
	}
}

// GetAllRolesByTenant 查询当前租户下的用户，管理员登陆时查询所有
func GetAllRolesByTenant(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			log.Printf("No viewer found in context")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No viewer found in context"})
		}
		tenantID := v.TenantID
		roleName := v.RoleName

		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", v.UserID, tenantID, roleName)

		var roles []*ent.Role
		var err error
		if roleName == "admin" || roleName == "superadmin" {
			roles, err = client.Role.Query().Where(role.HasTenantWith(tenant.IDEQ(tenantID))).All(c.Request().Context())
			if err != nil {
				log.Printf("Error querying roles for tenant %s: %v", tenantID, err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying roles from database"})
			}
		} else {
			roles, err = client.Role.Query().Where(role.HasUsersWith(user.IDEQ(v.UserID))).All(c.Request().Context())
			if err != nil {
				log.Printf("Error querying roles for user %s: %v", v.UserID, err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying roles from database"})
			}
		}

		log.Printf("Roles for tenant %s: %v", tenantID, roles)
		return c.JSON(http.StatusOK, roles)
	}
}

// CreateRole 创建role
func CreateRole(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 从请求上下文中获取租户ID
		v := viewer.FromContext(c.Request().Context())
		if v == nil || v.RoleName != "admin" {
			log.Printf("No viewer found in context or not authorized")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		var roles []*RoleDTO
		if err := c.Bind(&roles); err != nil {
			log.Printf("Error binding role: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		createdRoles := make([]*ent.Role, 0, len(roles))

		for _, dto := range roles {

			r, err := client.Role.Create().
				SetName(dto.Name).
				SetDescription(dto.Description).
				AddTenantIDs(v.TenantID). // 关联到租户
				AddPermissionIDs(dto.PermissionIDs...).
				Save(c.Request().Context())
			if err != nil {
				log.Printf("Error creating role: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create role"})
			}
			createdRoles = append(createdRoles, r)
			return c.JSON(http.StatusCreated, r)
		}
		return c.JSON(http.StatusCreated, map[string]string{"message": "Roles created successfully"})
	}
}

// IsDefaultRole 检查角色是否为默认角色
func IsDefaultRole(roleName string) bool {
	defaultRoles := []string{"Admin", "Developer", "Auditor", "User"}
	for _, r := range defaultRoles {
		if r == roleName {
			return true
		}
	}
	return false
}

func DeleteRoleByName(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		// 从请求上下文中获取租户ID
		v := viewer.FromContext(c.Request().Context())
		if v == nil || v.RoleName != "admin" {
			log.Printf("No viewer found in context or not authorized")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		type RoleDTO struct {
			Name string `json:"name"`
		}
		dto := new(RoleDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding role: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		if IsDefaultRole(dto.Name) {
			log.Printf("Attempt to delete default role: %s", dto.Name)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot delete default role"})
		}

		ro, err := client.Role.Query().Where(role.NameEQ(dto.Name)).Only(c.Request().Context())
		if ent.IsNotFound(err) {
			log.Printf("Role not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
		}
		if err != nil {
			log.Printf("Error querying role: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying role from database"})
		}

		err = client.Role.DeleteOne(ro).Exec(c.Request().Context())
		if ent.IsNotFound(err) {
			log.Printf("Role not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
		}
		if err != nil {
			log.Printf("Error deleting role: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error deleting role"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}
