package handler

import (
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/session"
	"github.com/willie-lin/cloud-terminal/viewer"
	"net/http"
	"strings"
	"time"
)

var sessionStore *sessions.CookieStore

func InitSessionStore() {
	sessionStore = sessions.NewCookieStore([]byte("cloud-terminal-secret-key-change-in-production"))
}

func GetSessionStore() *sessions.CookieStore {
	return sessionStore
}

func CheckSession(c *echo.Context) error {
	sess, err := sessionStore.Get(c.Request(), "session")
	if err != nil || sess.IsNew {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Session expired or not found"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Session is active"})
}

// ========== ent.Session CRUD ==========

type SessionUpdateDTO struct {
	Status         *string    `json:"status"`
	EndedAt        *time.Time `json:"ended_at"`
	RemoteAddress  *string    `json:"remote_address"`
}

func ListSessions(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can list sessions"})
		}

		query := client.Session.Query()
		if status := c.QueryParam("status"); status != "" {
			query = query.Where(session.StatusEQ(session.Status(status)))
		}
		if principal := c.QueryParam("principal_urn"); principal != "" {
			query = query.Where(session.PrincipalUrnEQ(principal))
		}
		sessions, err := query.Order(ent.Desc(session.FieldStartedAt)).All(c.Request().Context())
		if err != nil {
			log.Printf("Error querying sessions: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying sessions"})
		}
		return c.JSON(http.StatusOK, sessions)
	}
}

func GetSession(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid session ID"})
		}
		s, err := client.Session.Query().
			Where(session.IDEQ(id)).
			Only(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Session not found"})
		}
		if err != nil {
			log.Printf("Error querying session: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying session"})
		}
		return c.JSON(http.StatusOK, s)
	}
}

func UpdateSession(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		if !isSuperAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins can update sessions"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid session ID"})
		}

		dto := new(SessionUpdateDTO)
		if err := c.Bind(dto); err != nil {
			log.Printf("Error binding session update: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
		}

		updater := client.Session.UpdateOneID(id)
		if dto.Status != nil {
			updater.SetStatus(session.Status(*dto.Status))
		}
		if dto.EndedAt != nil {
			updater.SetEndedAt(*dto.EndedAt)
		}
		if dto.RemoteAddress != nil {
			updater.SetRemoteAddress(*dto.RemoteAddress)
		}

		s, err := updater.Save(c.Request().Context())
		if err != nil {
			log.Printf("Error updating session: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update session"})
		}
		return c.JSON(http.StatusOK, s)
	}
}

func DeleteSession(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		if !isSuperAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only super admins can delete sessions"})
		}

		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid session ID"})
		}

		err := client.Session.DeleteOneID(id).Exec(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Session not found"})
		}
		if err != nil {
			log.Printf("Error deleting session: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete session"})
		}
		return c.NoContent(http.StatusNoContent)
	}
}
