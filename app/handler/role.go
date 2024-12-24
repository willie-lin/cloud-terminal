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
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			log.Printf("No viewer found in context")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No viewer found in context"})
		}
		userID := v.UserID
		tenantID := v.TenantID
		roleName := v.RoleName

		//// 查询用户的角色
		//userRoles, err := client.User.Query().
		//	Where(user.IDEQ(userID)).
		//	QueryRoles().
		//	All(c.Request().Context())
		//if err != nil {
		//	log.Printf("Error querying roles for user %s: %v", userID, err)
		//	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying roles from database"})
		//}

		// 查询该租户下的所有角色
		//allTenantRoles, err := client.Tenant.Query().Where(tenant.IDEQ(tenantID)).QueryRoles().All(c.Request().Context())
		//if err != nil {
		//	log.Printf("Error querying all roles for tenant %s: %v", tenantID, err)
		//	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying roles from database"})
		//}

		// 查询该租户下的所有角色
		allTenantRoles, err := client.Role.Query().Where(role.HasTenantWith(tenant.IDEQ(tenantID))).All(c.Request().Context())
		if err != nil {
			log.Printf("Error querying all roles for tenant %s: %v", tenantID, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying roles from database"})
		}
		// 如果是管理员，查询所有属于同一租户的角色
		var roles []*ent.Role
		if roleName == "Admin" { // 直接使用 rolename 判断是否是管理员
			roles = allTenantRoles // 管理员返回所有角色
		} else {
			// 非管理员需要过滤，只返回分配给用户的角色
			userRoles, err := client.User.Query().
				Where(user.IDEQ(userID)).
				QueryRoles().
				All(c.Request().Context())

			if err != nil {
				log.Printf("Error querying user roles for user %s: %v", v.UserID, err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying user roles from database"})
			}

			userRoleIDs := make(map[string]struct{})
			for _, ur := range userRoles {
				userRoleIDs[ur.ID.String()] = struct{}{}
			}

			for _, atr := range allTenantRoles {
				if _, ok := userRoleIDs[atr.ID.String()]; ok {
					roles = append(roles, atr)
				}
			}
		}

		log.Printf("Roles for user in tenant %s: %v", tenantID, roles)
		return c.JSON(http.StatusOK, roles)
	}
}

// CreateRole 创建role
func CreateRole(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 从请求上下文中获取租户ID
		v := viewer.FromContext(c.Request().Context())
		if v == nil || v.RoleName != "Tenant Admin" {
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
			// 从请求上下文中获取租户ID
			//tenantID := c.Get("tenant_id").(uuid.UUID)

			r, err := client.Role.Create().
				SetName(dto.Name).
				SetDescription(dto.Description).
				//SetTenantID(v.TenantID). // 关联到租户
				//AddTenantIDs(v.TenantID).
				AddPermissionIDs(dto.PermissionIDs...).
				Save(context.Background())
			if err != nil {
				log.Printf("Error creating role: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create role"})
			}
			fmt.Printf("Created role with ID: %s\n", role.ID)
			createdRoles = append(createdRoles, r)
			return c.JSON(http.StatusCreated, r)
		}
		fmt.Println(createdRoles)
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
