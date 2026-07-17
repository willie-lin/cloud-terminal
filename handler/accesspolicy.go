package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/accesspolicy"
	"github.com/willie-lin/cloud-terminal/ent/schema"
	"github.com/willie-lin/cloud-terminal/viewer"
	"net/http"
	"strings"
	"time"
)

type AccessPolicyDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CheckAccessPolicyName(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// 从请求上下文中获取 viewer 信息
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			log.Printf("No viewer found in context")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No viewer found in context"})
		}

		dto := new(AccessPolicyDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding  accessPolicy: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		// 在 viewer 上下文中执行查询
		exists, err := client.AccessPolicy.Query().Where(accesspolicy.NameEQ(dto.Name)).Exist(c.Request().Context())
		if err != nil {
			log.Printf("Error checking accessPolicy name: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error checking accessPolicy  name from database"})
		}
		return c.JSON(http.StatusOK, map[string]bool{"exists": exists})
	}
}

// GetAllAccessPolicyByAccountByTenant   查询当前租户下的用户，管理员登陆时查询所有
func GetAllAccessPolicyByAccountByTenant(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			log.Printf("No viewer found in context")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No viewer found in context"})
		}
		tenantID := v.TenantID
		userID := v.UserID
		roleName := v.RoleName
		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", userID, tenantID, roleName)

		isSuperAdmin := roleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(roleName), "tenant_admin") // Or use a more precise matching logic

		var accessPolicies []*ent.AccessPolicy
		var err error

		if isSuperAdmin || isTenantAdmin {
			accessPolicies, err = client.AccessPolicy.Query().All(c.Request().Context())
		} else {
			accessPolicies, err = client.AccessPolicy.Query().All(c.Request().Context())
		}

		if err != nil {
			log.Printf("Error querying roles: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying roles from database"})
		}

		log.Printf("Roles for tenant %s: %v", tenantID, accessPolicies)
		return c.JSON(http.StatusOK, accessPolicies)
	}
}

//
//// CreateRole 创建role
//func CreateRole(client *ent.Client) echo.HandlerFunc {
//	return func(c *echo.Context) error {
//		// 从请求上下文中获取租户ID
//		v := viewer.FromContext(c.Request().Context())
//		if v == nil {
//			log.Printf("No viewer found in context or not authorized")
//			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
//		}
//		tenantID := v.TenantID
////		userID := v.UserID
//		roleName := v.RoleName
//		log.Printf("Viewer info: UserID=%s, TenantID=%s, RoleName=%s", userID, tenantID, roleName)
//
//		isSuperAdmin := roleName == "super_admin"
//		isTenantAdmin := strings.Contains(strings.ToLower(roleName), "tenant_admin") // Or use a more precise matching logic
//
//		//if isSuperAdmin || isTenantAdmin {
//		if !isSuperAdmin && !isTenantAdmin {
//			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins and tenant admins can create roles"})
//		}
//
//		var roles []*RoleDTO
//		if err := c.Bind(&roles); err != nil {
//			log.Printf("Error binding role: %v", err)
//			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
//		}
//
//		createdRoles := make([]*ent.Role, 0, len(roles))
//
//		for _, dto := range roles {
//			r, err := client.Role.Create().
//				SetName(dto.Name).
//				SetDescription(dto.Description).
//				
//				Save(c.Request().Context())
//			if err != nil {
//				log.Printf("Error creating role: %v", err)
//				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create role"})
//			}
//			createdRoles = append(createdRoles, r)
//			//return c.JSON(http.StatusCreated, r)
//		}
//		// 返回创建的所有角色
//		//return c.JSON(http.StatusCreated, createdRoles)
//		return c.JSON(http.StatusCreated, map[string]string{"message": "Roles created successfully"})
//	}
//}
//
//// IsDefaultRole 检查角色是否为默认角色
//func IsDefaultRole(roleName string) bool {
//	defaultRoles := []string{"Admin", "Developer", "Auditor", "User"}
//	for _, r := range defaultRoles {
//		if r == roleName {
//			return true
//		}
//	}
//	return false
//}
//
//func DeleteRoleByName(client *ent.Client) echo.HandlerFunc {
//	return func(c *echo.Context) error {
//
//		//// 从请求上下文中获取租户ID
//		//v := viewer.FromContext(c.Request().Context())
//		//if v == nil || v.RoleName != "admin" {
//		//	log.Printf("No viewer found in context or not authorized")
//		//	return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
//		//}
//		type RoleDTO struct {
//			Name string `json:"name"`
//		}
//		dto := new(RoleDTO)
//		if err := c.Bind(&dto); err != nil {
//			log.Printf("Error binding role: %v", err)
//			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
//		}
//		if IsDefaultRole(dto.Name) {
//			log.Printf("Attempt to delete default role: %s", dto.Name)
//			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cannot delete default role"})
//		}
//
//		ro, err := client.Role.Query().Where(role.NameEQ(dto.Name)).Only(c.Request().Context())
//		if ent.IsNotFound(err) {
//			log.Printf("Role not found: %v", err)
//			return c.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
//		}
//		if err != nil {
//			log.Printf("Error querying role: %v", err)
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying role from database"})
//		}
//
//		err = client.Role.DeleteOne(ro).Exec(c.Request().Context())
//		if ent.IsNotFound(err) {
//			log.Printf("Role not found: %v", err)
//			return c.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
//		}
//		if err != nil {
//			log.Printf("Error deleting role: %v", err)
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error deleting role"})
//		}
//		return c.NoContent(http.StatusNoContent)
//	}
//}

// ==================== AccessPolicy CRUD ====================

type AccessPolicyCreateDTO struct {
	Name          string                     `json:"name"`
	Description   string                     `json:"description"`
	Version       string                     `json:"version"`
	Statements    []schema.PolicyStatement   `json:"statements"`
	Immutable     bool                       `json:"immutable"`
	Priority      int                        `json:"priority"`
	EffectiveDate *time.Time                 `json:"effective_date"`
	ExpiryDate    *time.Time                 `json:"expiry_date"`
}

type AccessPolicyUpdateDTO struct {
	Description   *string                    `json:"description"`
	Version       *string                    `json:"version"`
	Statements    *[]schema.PolicyStatement  `json:"statements"`
	Immutable     *bool                      `json:"immutable"`
	Priority      *int                       `json:"priority"`
	EffectiveDate *time.Time                 `json:"effective_date"`
	ExpiryDate    *time.Time                 `json:"expiry_date"`
}

// CreateAccessPolicy creates a new access policy
func CreateAccessPolicy(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can create access policies"})
		}

		dto := new(AccessPolicyCreateDTO)
		if err := c.Bind(dto); err != nil {
			log.Printf("Error binding access policy: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		if dto.Version == "" {
			dto.Version = "v1"
		}
		if dto.Statements == nil {
			dto.Statements = []schema.PolicyStatement{}
		}

		creator := client.AccessPolicy.Create().
			SetName(dto.Name).
			SetDescription(dto.Description).
			SetVersion(dto.Version).
			SetStatements(dto.Statements).
			SetImmutable(dto.Immutable).
			SetPriority(dto.Priority)

		if dto.EffectiveDate != nil {
			creator.SetEffectiveDate(*dto.EffectiveDate)
		}
		if dto.ExpiryDate != nil {
			creator.SetExpiryDate(*dto.ExpiryDate)
		}

		p, err := creator.Save(c.Request().Context())
		if err != nil {
			log.Printf("Error creating access policy: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create access policy"})
		}

		return c.JSON(http.StatusCreated, p)
	}
}

// GetAccessPolicy gets a single access policy by ID
func GetAccessPolicy(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid access policy ID"})
		}

		p, err := client.AccessPolicy.Query().
			Where(accesspolicy.IDEQ(id)).
			Only(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Access policy not found"})
		}
		if err != nil {
			log.Printf("Error querying access policy: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying access policy"})
		}

		return c.JSON(http.StatusOK, p)
	}
}

// UpdateAccessPolicy updates an access policy by ID
func UpdateAccessPolicy(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can update access policies"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid access policy ID"})
		}

		dto := new(AccessPolicyUpdateDTO)
		if err := c.Bind(dto); err != nil {
			log.Printf("Error binding access policy update: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		updater := client.AccessPolicy.UpdateOneID(id)
		if dto.Description != nil {
			updater.SetDescription(*dto.Description)
		}
		if dto.Version != nil {
			updater.SetVersion(*dto.Version)
		}
		if dto.Statements != nil {
			updater.SetStatements(*dto.Statements)
		}
		if dto.Immutable != nil {
			updater.SetImmutable(*dto.Immutable)
		}
		if dto.Priority != nil {
			updater.SetPriority(*dto.Priority)
		}
		if dto.EffectiveDate != nil {
			updater.SetEffectiveDate(*dto.EffectiveDate)
		}
		if dto.ExpiryDate != nil {
			updater.SetExpiryDate(*dto.ExpiryDate)
		}

		p, err := updater.Save(c.Request().Context())
		if err != nil {
			log.Printf("Error updating access policy: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update access policy"})
		}

		return c.JSON(http.StatusOK, p)
	}
}

// DeleteAccessPolicy deletes an access policy by ID
func DeleteAccessPolicy(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can delete access policies"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid access policy ID"})
		}

		err := client.AccessPolicy.DeleteOneID(id).Exec(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Access policy not found"})
		}
		if err != nil {
			log.Printf("Error deleting access policy: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete access policy"})
		}

		return c.NoContent(http.StatusNoContent)
	}
}
