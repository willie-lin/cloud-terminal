package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/privacy"
	"github.com/willie-lin/cloud-terminal/ent/role"
	"github.com/willie-lin/cloud-terminal/ent/group"
	"github.com/willie-lin/cloud-terminal/ent/tenant"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"github.com/willie-lin/cloud-terminal/viewer"
)

type TenantDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// 可选填：在创建租户时同步初始化该租户的管理员
	AdminEmail    string `json:"admin_email,omitempty"`
	AdminPassword string `json:"admin_password,omitempty"`
	AdminUsername string `json:"admin_username,omitempty"`
}

func CreateTenant(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		if v.RoleName != "super_admin" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins can create tenants"})
		}

		dto := new(TenantDTO)
		if err := c.Bind(dto); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		if strings.TrimSpace(dto.Name) == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tenant name is required"})
		}

		// 使用 Privacy 提权上下文，确保建租户操作不受拦截
		ctx := privacy.DecisionContext(c.Request().Context(), privacy.Allow)

		// 1. 创建租户
		createdTenant, err := client.Tenant.
			Create().
			SetName(strings.TrimSpace(dto.Name)).
			SetDescription(dto.Description).
			Save(ctx)

		if err != nil {
			log.Printf("Error creating tenant: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create tenant: %v", err)})
		}

		adminCreated := false
		var createdAdminID string

		// 2. 如果提供了初始管理员信息，校验必填与长度
		if dto.AdminEmail != "" || dto.AdminPassword != "" || dto.AdminUsername != "" {
			if strings.TrimSpace(dto.AdminEmail) == "" || dto.AdminPassword == "" || strings.TrimSpace(dto.AdminUsername) == "" {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "管理员邮箱、密码与用户名均为必填项"})
			}
			if len(strings.TrimSpace(dto.AdminUsername)) < 6 {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "管理员用户名至少需要 6 个字符"})
			}

			groupName := strings.TrimSpace(dto.Name) + "_Group"
			g, err := client.Group.Create().
				SetName(groupName).
				Save(ctx)
			if err != nil {
				log.Printf("Warning: failed to create tenant group: %v", err)
			} else {
				// 获取 tenant_admin 角色
				tenantRole, rErr := client.Role.Query().Where(role.NameEQ("tenant_admin")).Only(ctx)
				if rErr != nil {
					// 兜底退回
					tenantRole, _ = client.Role.Query().First(ctx)
				}

				hashedPassword, pErr := utils.GenerateFromPassword([]byte(dto.AdminPassword), utils.DefaultCost)
				if pErr == nil {
					username := strings.TrimSpace(dto.AdminUsername)

					uBuilder := client.User.Create().
						SetUsername(username).
						SetEmail(strings.TrimSpace(dto.AdminEmail)).
						SetPassword(string(hashedPassword)).
						SetGroup(g)

					if tenantRole != nil {
						uBuilder.AddRoles(tenantRole)
					}

					u, uErr := uBuilder.Save(ctx)
					if uErr == nil && u != nil {
						adminCreated = true
						createdAdminID = u.ID
						log.Printf("Successfully created tenant admin user: %s (Email: %s)", username, dto.AdminEmail)
					} else {
						log.Printf("Error: failed to create initial tenant admin user: %v", uErr)
						return c.JSON(http.StatusInternalServerError, map[string]string{
							"error": fmt.Sprintf("Tenant created, but admin user creation failed: %v", uErr),
						})
					}
				}
			}
		}

		return c.JSON(http.StatusCreated, map[string]any{
			"tenant":           createdTenant,
			"admin_created":    adminCreated,
			"admin_user_id":    createdAdminID,
		})
	}
}

func GetTenantByName(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		dto := new(TenantDTO)

		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		tant, err := client.Tenant.Query().Where(tenant.NameEQ(dto.Name)).Only(context.Background())
		if err != nil {
			if ent.IsNotFound(err) {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying Tenant from database"})
		}

		return c.JSON(http.StatusOK, tant)
	}

}

func DeleteTenantName(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {

		dto := new(TenantDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding Tenant: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		tent, err := client.Tenant.Query().Where(tenant.NameEQ(dto.Name)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("Tenant not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
		}
		if err != nil {
			log.Printf("Error querying Tenant: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying Tenant from database"})
		}

		err = client.Tenant.DeleteOne(tent).Exec(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("Tenant not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
		}
		if err != nil {
			log.Printf("Error deleting Tenant: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error deleting Tenant"})
		}
		return c.NoContent(http.StatusNoContent)
	}

}

func CheckTenantName(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {

		dto := new(TenantDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding tenant: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		//fmt.Println(dto.Name)
		// 检查tenant是否已经存在
		ctx := context.Background()
		exists, err := client.Tenant.Query().Where(tenant.NameEQ(dto.Name)).Exist(ctx)
		if err != nil {
			log.Printf("Error checking tenant: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error checking tenant from database"})
		}

		//fmt.Println(exists)
		return c.JSON(http.StatusOK, map[string]bool{"exists": exists})
	}
}


// ==================== RESTful Tenant CRUD ====================

type TenantUpdateDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

func ListTenants(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		if !isSuperAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins can list tenants"})
		}

		tenants, err := client.Tenant.Query().All(c.Request().Context())
		if err != nil {
			log.Printf("Error querying tenants: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying tenants"})
		}
		return c.JSON(http.StatusOK, tenants)
	}
}

func GetTenantByID(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tenant ID"})
		}
		t, err := client.Tenant.Query().
			Where(tenant.IDEQ(id)).
			Only(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
		}
		if err != nil {
			log.Printf("Error querying tenant: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying tenant"})
		}
		return c.JSON(http.StatusOK, t)
	}
}

func UpdateTenant(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		if !isSuperAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins can update tenants"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tenant ID"})
		}

		dto := new(TenantUpdateDTO)
		if err := c.Bind(dto); err != nil {
			log.Printf("Error binding tenant update: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		ctx := privacy.DecisionContext(c.Request().Context(), privacy.Allow)

		updater := client.Tenant.UpdateOneID(id)
		if dto.Name != nil {
			updater.SetName(*dto.Name)
		}
		if dto.Description != nil {
			updater.SetDescription(*dto.Description)
		}
		if dto.Status != nil {
			updater.SetStatus(tenant.Status(*dto.Status))
		}

		t, err := updater.Save(ctx)
		if err != nil {
			log.Printf("Error updating tenant: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update tenant"})
		}
		return c.JSON(http.StatusOK, t)
	}
}

func DeleteTenantByID(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		if !isSuperAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins can delete tenants"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tenant ID"})
		}

		ctx := privacy.DecisionContext(c.Request().Context(), privacy.Allow)

		err := client.Tenant.DeleteOneID(id).Exec(ctx)
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
		}
		if err != nil {
			log.Printf("Error deleting tenant: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete tenant"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}

// ==================== Tenant User Management ====================

// ListTenantUsers 查询指定租户 Group 下的所有用户
func ListTenantUsers(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		if v.RoleName != "super_admin" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins can view tenant users"})
		}

		tenantID := c.Param("id")
		if _, err := uuid.Parse(tenantID); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tenant ID"})
		}

		ctx := privacy.DecisionContext(c.Request().Context(), privacy.Allow)

		// 通过 tenant Group 查找用户（租户 Group 名约定为 <tenantName>_Group）
		t, err := client.Tenant.Query().Where(tenant.IDEQ(tenantID)).Only(ctx)
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying tenant"})
		}

		groupName := t.Name + "_Group"
		g, gErr := client.Group.Query().Where(group.NameEQ(groupName)).WithUsers(func(uq *ent.UserQuery) {
			uq.WithRoles()
		}).Only(ctx)
		if gErr != nil {
			// 该租户没有对应 group，返回空列表
			return c.JSON(http.StatusOK, []any{})
		}

		return c.JSON(http.StatusOK, g.Edges.Users)
	}
}

// AddUserToTenant 将已有用户加入指定租户 Group
func AddUserToTenant(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		if v.RoleName != "super_admin" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins can manage tenant users"})
		}

		tenantID := c.Param("id")
		if _, err := uuid.Parse(tenantID); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tenant ID"})
		}

		type DTO struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Username string `json:"username"`
			RoleName string `json:"role_name"`
		}
		dto := new(DTO)
		if err := c.Bind(dto); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		if strings.TrimSpace(dto.Email) == "" || dto.Password == "" || len(strings.TrimSpace(dto.Username)) < 6 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email、username（至少6字符）和 password 均为必填"})
		}

		ctx := privacy.DecisionContext(c.Request().Context(), privacy.Allow)

		t, err := client.Tenant.Query().Where(tenant.IDEQ(tenantID)).Only(ctx)
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying tenant"})
		}

		groupName := t.Name + "_Group"
		g, gErr := client.Group.Query().Where(group.NameEQ(groupName)).Only(ctx)
		if gErr != nil {
			// 该租户还没有 Group，自动创建
			g, gErr = client.Group.Create().SetName(groupName).Save(ctx)
			if gErr != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create tenant group: %v", gErr)})
			}
		}

		hashedPassword, pErr := utils.GenerateFromPassword([]byte(dto.Password), utils.DefaultCost)
		if pErr != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error hashing password"})
		}

		targetRole := strings.TrimSpace(dto.RoleName)
		if targetRole == "" {
			targetRole = "user"
		}

		uBuilder := client.User.Create().
			SetEmail(strings.TrimSpace(dto.Email)).
			SetUsername(strings.TrimSpace(dto.Username)).
			SetPassword(string(hashedPassword)).
			SetOnline(true).
			SetStatus(true).
			SetGroup(g)

		r, rErr := client.Role.Query().Where(role.NameEQ(targetRole)).Only(ctx)
		if rErr != nil {
			r, _ = client.Role.Query().Where(role.NameEQ("user")).Only(ctx)
		}
		if r != nil {
			uBuilder.AddRoles(r)
		}

		u, uErr := uBuilder.Save(ctx)
		if uErr != nil {
			log.Printf("Error creating tenant user: %v", uErr)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create user: %v", uErr)})
		}

		return c.JSON(http.StatusCreated, map[string]string{
			"user_id":  u.ID,
			"username": u.Username,
			"email":    u.Email,
		})
	}
}

// RemoveUserFromTenant 将指定用户从租户 Group 中移除（解绑，不删除用户账号）
func RemoveUserFromTenant(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		if v.RoleName != "super_admin" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins can manage tenant users"})
		}

		tenantID := c.Param("id")
		userID := c.Param("uid")
		if _, err := uuid.Parse(tenantID); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tenant ID"})
		}
		if _, err := uuid.Parse(userID); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		ctx := privacy.DecisionContext(c.Request().Context(), privacy.Allow)

		// 解除用户与 Group 的绑定（清除 group 外键，而不是删除用户）
		_, err := client.User.UpdateOneID(userID).ClearGroup().Save(ctx)
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		if err != nil {
			log.Printf("Error removing user from tenant: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to remove user: %v", err)})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "User removed from tenant successfully"})
	}
}

