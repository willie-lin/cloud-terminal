package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/uber/jaeger-client-go/crossdock/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/tenant"
	"net/http"
)

type TenantDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateTenant(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		dto := new(TenantDTO)

		if err := c.Bind(dto); err != nil {
			//return echo.NewHTTPError(http.StatusBadRequest, "无效的请求数据")
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
	return func(c echo.Context) error {
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
	return func(c echo.Context) error {

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
