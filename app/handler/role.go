package handler

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"github.com/willie-lin/cloud-terminal/app/database/ent/role"
	"net/http"
)

type RoleDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

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

// CreateRole 创建role
func CreateRole(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) (err error) {

		var roles []*RoleDTO

		if err := c.Bind(&roles); err != nil {
			log.Printf("Error binding role: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		for _, dto := range roles {
			ro, err := client.Role.Create().
				SetName(dto.Name).
				SetDescription(dto.Description).
				Save(context.Background())
			if err != nil {
				log.Printf("Error creating role: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create role"})
			}
			fmt.Printf("Created role with ID: %s\n", ro.ID)
		}
		return c.JSON(http.StatusCreated, map[string]string{"message": "Roles created successfully"})
	}
}

func DeleteRoleByName(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
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
