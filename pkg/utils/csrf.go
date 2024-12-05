package utils

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

// SetCSRFToken

func SetCSRFToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method != http.MethodOptions {
			token := c.Get(middleware.DefaultCSRFConfig.ContextKey)
			if token == nil {
				// 首先检查Cookie中是否已经有CSRF令牌
				csrfCookie, err := c.Cookie("_csrf")
				if err != nil || csrfCookie == nil {
					// 如果cookie不存在或者获取失败，则生成新的CSRF令牌
					token, err = GenerateRandomKey(32)
					if err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, "生成 CSRF 令牌失败")
					}
				} else {
					// 如果cookie存在，则使用现有的CSRF令牌
					token = csrfCookie.Value
				}
			}

			csrfToken := token.(string)
			cookie := &http.Cookie{
				Name:     "_csrf",
				Value:    csrfToken,
				Path:     "/",
				Domain:   c.Request().Host,      // 使用请求的主机名作为Domain
				Secure:   true,                  // 如果在生产环境中使用HTTPS，确保设置为true
				HttpOnly: true,                  // 确保HttpOnly设置
				SameSite: http.SameSiteNoneMode, // 允许跨站请求
				MaxAge:   3600,                  // 1小时
			}
			c.SetCookie(cookie)
			c.Response().Header().Set("X-CSRF-Token", csrfToken)
			//fmt.Printf("使用或生成的 CSRF 令牌: %s\n", csrfToken)
		}
		return next(c)
	}
}
