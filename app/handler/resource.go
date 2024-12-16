package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/viewer"
	"log"
	"net/http"
)

func CreateResource(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil || v.RoleName != "Tenant Admin" {
			log.Printf("No viewer found in context or not authorized")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		type ResourceDTO struct {
			Name        string `json:"name"`
			Type        string `json:"type"`
			Value       string `json:"value"`
			Description string `json:"description"`
		}

		var dto ResourceDTO
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding resource: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		resource, err := client.Resource.Create().
			SetName(dto.Name).
			SetType(dto.Type).
			SetValue(dto.Value).
			SetDescription(dto.Description).
			SetTenantID(v.TenantID). // 关联到当前租户
			Save(context.Background())
		if err != nil {
			log.Printf("Error creating resource: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create resource"})
		}

		return c.JSON(http.StatusCreated, resource)
	}
}
