package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/willie-lin/cloud-terminal/ent"
	"github.com/willie-lin/cloud-terminal/ent/auditlog"
	"github.com/willie-lin/cloud-terminal/viewer"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ListAuditLogs(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		v := viewer.FromContext(c.Request().Context())
		if v == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		isSuperAdmin := v.RoleName == "super_admin"
		isTenantAdmin := strings.Contains(strings.ToLower(v.RoleName), "tenant_admin")
		if !isSuperAdmin && !isTenantAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admins can list audit logs"})
		}

		query := client.AuditLog.Query()

		if username := c.QueryParam("username"); username != "" {
			query = query.Where(auditlog.UsernameEQ(username))
		}
		if action := c.QueryParam("action"); action != "" {
			query = query.Where(auditlog.ActionEQ(action))
		}
		if result := c.QueryParam("result"); result != "" {
			query = query.Where(auditlog.ResultEQ(result))
		}
		if sessionID := c.QueryParam("session_id"); sessionID != "" {
			query = query.Where(auditlog.SessionIDEQ(sessionID))
		}
		if startStr := c.QueryParam("started_at_from"); startStr != "" {
			if t, err := time.Parse(time.RFC3339, startStr); err == nil {
				query = query.Where(auditlog.StartedAtGTE(t))
			}
		}
		if endStr := c.QueryParam("started_at_to"); endStr != "" {
			if t, err := time.Parse(time.RFC3339, endStr); err == nil {
				query = query.Where(auditlog.StartedAtLTE(t))
			}
		}

		page, _ := strconv.Atoi(c.QueryParam("page"))
		pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
		if pageSize <= 0 || pageSize > 100 {
			pageSize = 20
		}
		if page < 1 {
			page = 1
		}

		total, err := query.Count(c.Request().Context())
		if err != nil {
			log.Printf("Error counting audit logs: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying audit logs"})
		}

		logs, err := query.
			Order(ent.Desc(auditlog.FieldStartedAt)).
			Limit(pageSize).
			Offset((page - 1) * pageSize).
			All(c.Request().Context())
		if err != nil {
			log.Printf("Error querying audit logs: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying audit logs"})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		})
	}
}

func GetAuditLog(client *ent.Client) echo.HandlerFunc {
	return func(c *echo.Context) error {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid audit log ID"})
		}
		logEntry, err := client.AuditLog.Query().
			Where(auditlog.IDEQ(id)).
			WithUser().
			WithResource().
			Only(c.Request().Context())
		if ent.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Audit log not found"})
		}
		if err != nil {
			log.Printf("Error querying audit log: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error querying audit log"})
		}
		return c.JSON(http.StatusOK, logEntry)
	}
}
