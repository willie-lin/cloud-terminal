package middlewarers

import (
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
)

func RBAC(enforcer *casbin.Enforcer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*ent.User)
			tenant := c.Get("tenant").(*ent.Tenant)

			// 检查权限
			allowed, err := enforcer.Enforce(user.ID, tenant.ID, c.Request().URL.Path, c.Request().Method)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Error checking permissions")
			}
			if !allowed {
				return echo.NewHTTPError(http.StatusForbidden, "Permission denied")
			}

			return next(c)
		}
	}
}
