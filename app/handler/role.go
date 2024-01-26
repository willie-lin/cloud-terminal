package handler

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"net/http"
)

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

// CreateRole 创建role
func CreateRole(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) (err error) {

		type RoleDTO struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		var roles []*RoleDTO

		if err := c.Bind(&roles); err != nil {
			log.Printf("Error binding user: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		for _, dto := range roles {
			role, err := client.Role.Create().
				SetName(dto.Name).
				SetDescription(dto.Description).
				Save(context.Background())
			if err != nil {
				log.Printf("Error creating role: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create role"})
			}
			fmt.Printf("Created role with ID: %s\n", role.ID)
		}
		return c.JSON(http.StatusCreated, map[string]string{"message": "Roles created successfully"})
	}
}
