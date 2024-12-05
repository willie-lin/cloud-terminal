package utils

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil || sess.IsNew {
			// 会话不存在或已过期，返回401状态码
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Session expired or not found"})
		}
		// 会话有效，继续处理请求
		return next(c)
	}
}
