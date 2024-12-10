package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/role"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"net/http"
)

func CheckRoleName(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		type RoleDTO struct {
			Name string `json:"name"`
		}

		dto := new(RoleDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding role: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		// 检查role是否已经存在
		exists, err := client.Role.Query().Where(role.NameEQ(dto.Name)).Exist(context.Background())
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
		roles, err := client.Role.Query().All(context.Background())
		if err != nil {
			log.Printf("Error querying roles: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying roles from database"})
		}
		return c.JSON(http.StatusOK, roles)
	}
}

func GetAllRolesByTenant(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 从请求上下文中获取租户ID
		//userID := c.Get("user_id").(uuid.UUID)
		tenantID := c.Get("tenant_id").(uuid.UUID)

		roles, err := client.Role.Query().Where(role.HasTenantWith(tenant.IDEQ(tenantID))).All(context.Background())
		if err != nil {
			log.Printf("Error querying roles: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying roles from database"})
		}
		return c.JSON(http.StatusOK, roles)
	}
}

//// CreateRole 创建一个新角色
//func CreateRole(client *ent.Client) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		type RoleDTO struct {
//			Name        string `json:"name"`
//			Description string `json:"description"`
//			TenantID    string `json:"tenant_id"` // 关联租户的ID
//		}
//
//		dto := new(RoleDTO)
//		if err := c.Bind(&dto); err != nil {
//			log.Printf("Error binding role: %v", err)
//			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
//		}
//
//		tenantID, err := uuid.Parse(dto.TenantID)
//		if err != nil {
//			log.Printf("Invalid tenant ID: %v", err)
//			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tenant ID"})
//		}
//
//		role, err := client.Role.Create().
//			SetName(dto.Name).
//			SetDescription(dto.Description).
//			SetTenantID(tenantID).
//			Save(context.Background())
//		if err != nil {
//			log.Printf("Error creating role: %v", err)
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating role in database"})
//		}
//
//		return c.JSON(http.StatusCreated, map[string]string{"roleID": role.ID.String()})
//	}
//}

// CreateRole 创建role
func CreateRole(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		type RoleDTO struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			TenantID    string `json:"tenant_id"` // 新增租户ID字段
		}

		var roles []RoleDTO
		if err := c.Bind(&roles); err != nil {
			log.Printf("Error binding role: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		for _, dto := range roles {
			// 从请求上下文中获取租户ID
			tenantID := c.Get("tenant_id").(uuid.UUID)

			role, err := client.Role.Create().
				SetName(dto.Name).
				SetDescription(dto.Description).
				SetTenantID(tenantID). // 关联到租户
				Save(context.Background())
			if err != nil {
				log.Printf("Error creating role: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create role"})
			}
			fmt.Printf("Created role with ID: %s\n", role.ID)
			return c.JSON(http.StatusCreated, role)
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "Roles created successfully"})
	}
}

func DeleteRoleByName(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		type RoleDTO struct {
			Name string `json:"name"`
		}

		dto := new(RoleDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding role: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		ro, err := client.Role.Query().Where(role.NameEQ(dto.Name)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("Role not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
		}
		if err != nil {
			log.Printf("Error querying role: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying role from database"})
		}

		err = client.Role.DeleteOne(ro).Exec(context.Background())
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
