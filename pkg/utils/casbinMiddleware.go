package utils

import (
	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CasbinMiddleware(enforcer *casbin.Enforcer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Request().Header.Get("user") // 从请求头中获取用户信息，具体实现根据您的需求
			path := c.Path()                       // 获取请求路径
			method := c.Request().Method           // 获取请求方法

			allowed, err := enforcer.Enforce(user, path, method)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Error enforcing policy")
			}
			if !allowed {
				return echo.NewHTTPError(http.StatusForbidden, "You don't have permission to access this resource")
			}
			return next(c)
		}
	}
}
