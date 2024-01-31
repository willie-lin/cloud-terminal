package handler

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/permission"
	"net/http"
)

func CheckPermissionName(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		type PermissionDTO struct {
			Name string `json:"name"`
		}

		dto := new(PermissionDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding permission: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}
		// 检查permission是否已经存在
		exists, err := client.Permission.Query().Where(permission.NameEQ(dto.Name)).Exist(context.Background())
		if err != nil {
			log.Printf("Error checking name: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error checking name from database"})
		}
		return c.JSON(http.StatusOK, map[string]bool{"exists": exists})
	}
}

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

func DeletePermissionByName(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// RoleDTO
		type PermissionDTO struct {
			Name string `json:"name"`
		}

		dto := new(PermissionDTO)
		if err := c.Bind(&dto); err != nil {
			log.Printf("Error binding permission: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		p, err := client.Permission.Query().Where(permission.NameEQ(dto.Name)).Only(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("Permission not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Permission not found"})
		}
		if err != nil {
			log.Printf("Error querying permission: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying role from database"})
		}

		err = client.Permission.DeleteOne(p).Exec(context.Background())
		if ent.IsNotFound(err) {
			log.Printf("Permission not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Permission not found"})
		}
		if err != nil {
			log.Printf("Error deleting permission: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error deleting permission"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}
