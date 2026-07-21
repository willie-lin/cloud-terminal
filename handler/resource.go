package handler

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/resource"
	"github.com/willie-lin/cloud-terminal/ent/tenant"
	"github.com/willie-lin/cloud-terminal/pkg/crypto"
	"github.com/willie-lin/cloud-terminal/viewer"
)

type ResourceCreateDTO struct {
	Urn         string                 `json:"urn"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	IP          string                 `json:"ip"`
	Port        int                    `json:"port"`
	Env         string                 `json:"env"`
	Region      string                 `json:"region"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	Details     map[string]interface{} `json:"details"`
	AuthData    map[string]interface{} `json:"auth_data"`
	HostKey     string                 `json:"host_key"`
}

type ResourceUpdateDTO struct {
	Urn         *string                 `json:"urn"`
	Name        *string                 `json:"name"`
	Type        *string                 `json:"type"`
	IP          *string                 `json:"ip"`
	Port        *int                    `json:"port"`
	Env         *string                 `json:"env"`
	Region      *string                 `json:"region"`
	Description *string                 `json:"description"`
	Status      *string                 `json:"status"`
	Details     *map[string]interface{} `json:"details"`
	AuthData    *map[string]interface{} `json:"auth_data"`
	HostKey     *string                `json:"host_key"`
}

// CreateResource creates a new resource under the viewer's tenant
func CreateResource(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			log.Printf("No viewer found in context")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can create resources"})
		}

		dto := new(ResourceCreateDTO)
		if err := c.Bind(dto); err != nil {
			log.Printf("Error binding resource: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		creator := client.Resource.Create().
			SetUrn(dto.Urn).
			SetName(dto.Name).
			SetIP(dto.IP).
			SetTenantID(v.TenantID.String())

		if dto.Type != "" {
			creator.SetType(resource.Type(dto.Type))
		}
		if dto.Port > 0 {
			creator.SetPort(dto.Port)
		}
		if dto.Env != "" {
			creator.SetEnv(resource.Env(dto.Env))
		}
		if dto.Region != "" {
			creator.SetRegion(dto.Region)
		}
		if dto.Description != "" {
			creator.SetDescription(dto.Description)
		}
		if dto.Status != "" {
			creator.SetStatus(resource.Status(dto.Status))
		}
		if dto.Details != nil {
			creator.SetDetails(dto.Details)
		}
		if dto.AuthData != nil {
			encAuth, err := crypto.EncryptAuthData(dto.AuthData)
			if err != nil {
				log.Printf("Error encrypting auth_data: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to encrypt sensitive data"})
			}
			creator.SetAuthData(encAuth)
		}
		if dto.HostKey != "" {
			creator.SetHostKey(dto.HostKey)
		}

		r, err := creator.Save(c.Request().Context())
		if err != nil {
			log.Printf("Error creating resource: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create resource"})
		}

		return c.JSON(http.StatusCreated, r)
	}
}

// ListResources lists resources — super admin sees all, tenant admin sees own tenant's
func ListResources(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			log.Printf("No viewer found in context")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can list resources"})
		}

		query := client.Resource.Query()
		if !isSuperAdmin {
			query = query.Where(resource.HasTenantWith(tenant.IDEQ(v.TenantID.String())))
		}
		resources, err := query.All(c.Request().Context())
		if err != nil {
			log.Printf("Error querying resources: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying resources"})
		}

		return c.JSON(http.StatusOK, resources)
	}
}

// GetResource gets a single resource by ID
func GetResource(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid resource ID"})
		}

		r, err := client.Resource.Query().
			Where(resource.IDEQ(id)).
			Only(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Resource not found"})
		}
		if err != nil {
			log.Printf("Error querying resource: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying resource"})
		}

		return c.JSON(http.StatusOK, r)
	}
}

// UpdateResource updates a resource by ID
func UpdateResource(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can update resources"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid resource ID"})
		}

		dto := new(ResourceUpdateDTO)
		if err := c.Bind(dto); err != nil {
			log.Printf("Error binding resource update: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		updater := client.Resource.UpdateOneID(id)
		if dto.Urn != nil {
			updater.SetUrn(*dto.Urn)
		}
		if dto.Name != nil {
			updater.SetName(*dto.Name)
		}
		if dto.Type != nil {
			updater.SetType(resource.Type(*dto.Type))
		}
		if dto.IP != nil {
			updater.SetIP(*dto.IP)
		}
		if dto.Port != nil {
			updater.SetPort(*dto.Port)
		}
		if dto.Env != nil {
			updater.SetEnv(resource.Env(*dto.Env))
		}
		if dto.Region != nil {
			updater.SetRegion(*dto.Region)
		}
		if dto.Description != nil {
			updater.SetDescription(*dto.Description)
		}
		if dto.Status != nil {
			updater.SetStatus(resource.Status(*dto.Status))
		}
		if dto.Details != nil {
			updater.SetDetails(*dto.Details)
		}
		if dto.AuthData != nil {
			encAuth, err := crypto.EncryptAuthData(*dto.AuthData)
			if err != nil {
				log.Printf("Error encrypting auth_data: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to encrypt sensitive data"})
			}
			updater.SetAuthData(encAuth)
		}
		if dto.HostKey != nil {
			updater.SetHostKey(*dto.HostKey)
		}

		r, err := updater.Save(c.Request().Context())
		if err != nil {
			log.Printf("Error updating resource: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update resource"})
		}

		return c.JSON(http.StatusOK, r)
	}
}

// DeleteResource deletes a resource by ID
func DeleteResource(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can delete resources"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid resource ID"})
		}

		err := client.Resource.DeleteOneID(id).Exec(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Resource not found"})
		}
		if err != nil {
			log.Printf("Error deleting resource: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete resource"})
		}

		return c.NoContent(http.StatusNoContent)
	}
}
