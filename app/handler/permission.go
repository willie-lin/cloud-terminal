package handler

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"net/http"
)

// GetAllPermissions  获取所有permission
func GetAllPermissions(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		permissions, err := client.Permission.Query().All(context.Background())
		if err != nil {
			log.Printf("Error querying permissions: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying permissions from database"})
		}
		return c.JSON(http.StatusOK, permissions)
	}
}

// CreatePermission  创建permission
func CreatePermission(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) (err error) {

		type PermissionDTO struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		var permissions []*PermissionDTO

		if err := c.Bind(&permissions); err != nil {
			log.Printf("Error binding permission: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		for _, dto := range permissions {
			ro, err := client.Permission.Create().
				SetName(dto.Name).
				SetDescription(dto.Description).
				Save(context.Background())
			if err != nil {
				log.Printf("Error creating permission: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create permission"})
			}
			fmt.Printf("Created permission with ID: %s\n", ro.ID)
		}
		return c.JSON(http.StatusCreated, map[string]string{"message": "Permission created successfully"})
	}
}
