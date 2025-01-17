package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/account"
	"github.com/willie-lin/cloud-terminal/app/database/ent/role"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"github.com/willie-lin/cloud-terminal/app/viewer"
	"net/http"
	"strings"
)

type RoleDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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

// GetAllRolesByAccountByTenant  查询当前租户下的用户，管理员登陆时查询所有
func GetAllRolesByAccountByTenant(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			log.Printf("No viewer found in context")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No viewer found in context"})
		}
		tenantID := v.TenantID
		accountID := v.AccountID
		userID := v.UserID
		roleName := v.RoleName
		log.Printf("Viewer info: UserID=%s, AccountID=%s,TenantID=%s, RoleName=%s", userID, accountID, tenantID, roleName)

		isSuperAdmin := roleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(roleName), "tenant_admin") // Or use a more precise matching logic

		var roles []*ent.Role
		var err error

		if isSuperAdmin || isTenantAdmin {
			roles, err = client.Role.Query().Where(role.HasAccountWith(account.And(account.ID(accountID), account.HasTenantWith(tenant.IDEQ(tenantID))))).All(c.Request().Context())
		} else {
			roles, err = client.Role.Query().Where(role.NameEQ(roleName)).All(c.Request().Context())
		}

		if err != nil {
			log.Printf("Error querying roles: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying roles from database"})
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
		if v == nil {
			log.Printf("No viewer found in context or not authorized")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		tenantID := v.TenantID
		accountID := v.AccountID
		userID := v.UserID
		roleName := v.RoleName
		log.Printf("Viewer info: UserID=%s, AccountID=%s,TenantID=%s, RoleName=%s", userID, accountID, tenantID, roleName)

		isSuperAdmin := roleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(roleName), "tenant_admin") // Or use a more precise matching logic

		//if isSuperAdmin || isTenantAdmin {
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins and tenant admins can create roles"})
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
				SetAccountID(accountID).
				Save(c.Request().Context())
			if err != nil {
				log.Printf("Error creating role: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create role"})
			}
			createdRoles = append(createdRoles, r)
			//return c.JSON(http.StatusCreated, r)
		}
		// 返回创建的所有角色
		//return c.JSON(http.StatusCreated, createdRoles)
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

		//// 从请求上下文中获取租户ID
		//v := viewer.FromContext(c.Request().Context())
		//if v == nil || v.RoleName != "admin" {
		//	log.Printf("No viewer found in context or not authorized")
		//	return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		//}
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
