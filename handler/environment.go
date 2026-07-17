package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/environment"
	"github.com/willie-lin/cloud-terminal/ent/tenant"
	"github.com/willie-lin/cloud-terminal/viewer"
	"net/http"
	"strings"
)

type EnvironmentCreateDTO struct {
	Name          string                   `json:"name"`
	Description   string                   `json:"description"`
	Image         string                   `json:"image"`
	Port          int                      `json:"port"`
	ResourceLimit map[string]interface{}   `json:"resource_limit"`
	EnvVars       map[string]interface{}   `json:"env_vars"`
	Volumes       []map[string]interface{} `json:"volumes"`
}

type EnvironmentUpdateDTO struct {
	Description   *string                   `json:"description"`
	Image         *string                   `json:"image"`
	Port          *int                      `json:"port"`
	Status        *string                   `json:"status"`
	ResourceLimit *map[string]interface{}   `json:"resource_limit"`
	EnvVars       *map[string]interface{}   `json:"env_vars"`
	Volumes       *[]map[string]interface{} `json:"volumes"`
}

func CreateEnvironment(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can create environments"})
		}

		dto := new(EnvironmentCreateDTO)
		if err := c.Bind(dto); err != nil {
			log.Printf("Error binding environment: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		creator := client.Environment.Create().
			SetName(dto.Name).
			SetDescription(dto.Description).
			SetImage(dto.Image).
			SetTenantID(v.TenantID.String())

		if dto.Port > 0 {
			creator.SetPort(dto.Port)
		}
		if dto.ResourceLimit != nil {
			creator.SetResourceLimit(dto.ResourceLimit)
		}
		if dto.EnvVars != nil {
			creator.SetEnvVars(dto.EnvVars)
		}
		if dto.Volumes != nil {
			creator.SetVolumes(dto.Volumes)
		}

		env, err := creator.Save(c.Request().Context())
		if err != nil {
			log.Printf("Error creating environment: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create environment"})
		}
		return c.JSON(http.StatusCreated, env)
	}
}

func ListEnvironments(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can list environments"})
		}

		query := client.Environment.Query()
		if !isSuperAdmin {
			query = query.Where(environment.HasTenantWith(tenant.IDEQ(v.TenantID.String())))
		}
		envs, err := query.All(c.Request().Context())
		if err != nil {
			log.Printf("Error querying environments: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying environments"})
		}
		return c.JSON(http.StatusOK, envs)
	}
}

func GetEnvironment(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid environment ID"})
		}
		env, err := client.Environment.Query().
			Where(environment.IDEQ(id)).
			Only(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Environment not found"})
		}
		if err != nil {
			log.Printf("Error querying environment: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying environment"})
		}
		return c.JSON(http.StatusOK, env)
	}
}

func UpdateEnvironment(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can update environments"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid environment ID"})
		}

		dto := new(EnvironmentUpdateDTO)
		if err := c.Bind(dto); err != nil {
			log.Printf("Error binding environment update: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		updater := client.Environment.UpdateOneID(id)
		if dto.Description != nil {
			updater.SetDescription(*dto.Description)
		}
		if dto.Image != nil {
			updater.SetImage(*dto.Image)
		}
		if dto.Port != nil {
			updater.SetPort(*dto.Port)
		}
		if dto.Status != nil {
			updater.SetStatus(environment.Status(*dto.Status))
		}
		if dto.ResourceLimit != nil {
			updater.SetResourceLimit(*dto.ResourceLimit)
		}
		if dto.EnvVars != nil {
			updater.SetEnvVars(*dto.EnvVars)
		}
		if dto.Volumes != nil {
			updater.SetVolumes(*dto.Volumes)
		}

		env, err := updater.Save(c.Request().Context())
		if err != nil {
			log.Printf("Error updating environment: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update environment"})
		}
		return c.JSON(http.StatusOK, env)
	}
}

func DeleteEnvironment(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		if !isSuperAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins can delete environments"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid environment ID"})
		}

		err := client.Environment.DeleteOneID(id).Exec(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Environment not found"})
		}
		if err != nil {
			log.Printf("Error deleting environment: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete environment"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}
