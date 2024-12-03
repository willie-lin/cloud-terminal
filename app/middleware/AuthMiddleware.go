package middleware

import (
	"github.com/willie-lin/cloud-terminal/app/auth"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(authEnforcer *auth.Enforcer, client *ent.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*ent.User)
			tenant, err := user.QueryTenant().Only(c.Request().Context())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user tenant")
			}

			allowed, err := authEnforcer.Enforce(user.Username, tenant.Name, c.Path(), c.Request().Method)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check permissions")
			}

			if !allowed {
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
			}

			return next(c)
		}
	}
}
