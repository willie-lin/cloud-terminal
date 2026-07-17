package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/group"
	"github.com/willie-lin/cloud-terminal/viewer"
	"net/http"
	"strings"
)

type GroupCreateDTO struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Attributes  map[string]interface{} `json:"attributes"`
	UserIDs     []string               `json:"user_ids"`
}

type GroupUpdateDTO struct {
	Description *string                `json:"description"`
	Status      *string                `json:"status"`
	Attributes  *map[string]interface{} `json:"attributes"`
	UserIDs     *[]string              `json:"user_ids"`
}

func CreateGroup(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can create groups"})
		}

		dto := new(GroupCreateDTO)
		if err := c.Bind(dto); err != nil {
			log.Printf("Error binding group: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		creator := client.Group.Create().
			SetName(dto.Name).
			SetDescription(dto.Description)

		if dto.Attributes != nil {
			creator.SetAttributes(dto.Attributes)
		}
		if len(dto.UserIDs) > 0 {
			creator.AddUserIDs(dto.UserIDs...)
		}

		g, err := creator.Save(c.Request().Context())
		if err != nil {
			log.Printf("Error creating group: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create group"})
		}

		return c.JSON(http.StatusCreated, g)
	}
}

func ListGroups(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		groups, err := client.Group.Query().All(c.Request().Context())
		if err != nil {
			log.Printf("Error querying groups: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying groups"})
		}
		return c.JSON(http.StatusOK, groups)
	}
}

func GetGroup(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid group ID"})
		}

		g, err := client.Group.Query().
			Where(group.IDEQ(id)).
			WithUsers().
			Only(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Group not found"})
		}
		if err != nil {
			log.Printf("Error querying group: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying group"})
		}

		return c.JSON(http.StatusOK, g)
	}
}

func UpdateGroup(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		if !isSuperAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins can update groups"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid group ID"})
		}

		dto := new(GroupUpdateDTO)
		if err := c.Bind(dto); err != nil {
			log.Printf("Error binding group update: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		updater := client.Group.UpdateOneID(id)
		if dto.Description != nil {
			updater.SetDescription(*dto.Description)
		}
		if dto.Status != nil {
			updater.SetStatus(group.Status(*dto.Status))
		}
		if dto.Attributes != nil {
			updater.SetAttributes(*dto.Attributes)
		}
		if dto.UserIDs != nil {
			updater.ClearUsers().AddUserIDs(*dto.UserIDs...)
		}

		g, err := updater.Save(c.Request().Context())
		if err != nil {
			log.Printf("Error updating group: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update group"})
		}

		return c.JSON(http.StatusOK, g)
	}
}

func DeleteGroup(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		if !isSuperAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins can delete groups"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid group ID"})
		}

		err := client.Group.DeleteOneID(id).Exec(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Group not found"})
		}
		if err != nil {
			log.Printf("Error deleting group: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete group"})
		}

		return c.NoContent(http.StatusNoContent)
	}
}
