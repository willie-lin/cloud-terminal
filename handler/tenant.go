package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/uber/jaeger-client-go/crossdock/log"
	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/tenant"
	"github.com/willie-lin/cloud-terminal/viewer"
	"net/http"
)

type TenantDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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
		createdTenant, err := client.Tenant.
			Create().
			SetName(dto.Name).
			SetDescription(dto.Description).
			Save(context.Background())

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "创建租户失败")
		}
		return c.JSON(http.StatusCreated, createdTenant)
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

		t, err := updater.Save(c.Request().Context())
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

		err := client.Tenant.DeleteOneID(id).Exec(c.Request().Context())
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
