package utils

import (
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// SetCSRFToken sets or generates the _csrf cookie and X-CSRF-Token response header dynamically.
func SetCSRFToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if c.Request().Method != http.MethodOptions {
			var csrfToken string
			token := c.Get(middleware.DefaultCSRFConfig.ContextKey)
			if token != nil {
				if t, ok := token.(string); ok {
					csrfToken = t
				}
			}

			if csrfToken == "" {
				csrfCookie, err := c.Cookie("_csrf")
				if err == nil && csrfCookie != nil && csrfCookie.Value != "" {
					csrfToken = csrfCookie.Value
				} else {
					genToken, err := GenerateRandomKey(32)
					if err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, "生成 CSRF 令牌失败")
					}
					csrfToken = genToken
				}
			}

			// 动态解析 Domain 与 Protocol
			host := c.Request().Host
			if h, _, err := net.SplitHostPort(host); err == nil {
				host = h
			}

			isProd := os.Getenv("ENV") == "production"
			isTLS := c.Request().TLS != nil || c.Scheme() == "https" || strings.EqualFold(c.Request().Header.Get("X-Forwarded-Proto"), "https")
			isSecure := isProd || isTLS

			domain := ""
			if isProd && host != "" && host != "localhost" && host != "127.0.0.1" {
				domain = host
			}

			sameSite := http.SameSiteLaxMode
			if isSecure {
				sameSite = http.SameSiteNoneMode
			}

			cookie := &http.Cookie{
				Name:     "_csrf",
				Value:    csrfToken,
				Path:     "/",
				Domain:   domain,
				Secure:   isSecure,
				HttpOnly: false, // 允许前端 JS 读取
				SameSite: sameSite,
				MaxAge:   3600,
			}
			c.SetCookie(cookie)
			c.Response().Header().Set("X-CSRF-Token", csrfToken)
			c.Set("csrf_token", csrfToken)
		}
		return next(c)
	}
}

//func SetCSRFToken(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c *echo.Context) error {
//		if c.Request().Method != http.MethodOptions {
//			csrfToken := ""
//
//			// 检查 Cookie 中是否已经有 CSRF 令牌
//			csrfCookie, err := c.Cookie("_csrf")
//			if err != nil || csrfCookie == nil {
//				// 如果 Cookie 不存在或者获取失败，则生成新的 CSRF 令牌
//				token, err := GenerateRandomKey(32)
//				if err != nil {
//					return echo.NewHTTPError(http.StatusInternalServerError, "生成 CSRF 令牌失败")
//				}
//				csrfToken = token
//			} else {
//				// 如果 Cookie 存在，则使用现有的 CSRF 令牌
//				csrfToken = csrfCookie.Value
//			}
//
//			cookie := &http.Cookie{
//				Name:  "_csrf",
//				Value: csrfToken,
//				Path:  "/",
//				//Domain:   c.Request().Host,      // 使用请求的主机名作为 Domain
//				Domain:   "localhost",           // 使用请求的主机名作为 Domain
//				Secure:   true,                  // 如果在生产环境中使用 HTTPS，确保设置为 true
//				HttpOnly: true,                  // 确保 HttpOnly 设置
//				SameSite: http.SameSiteNoneMode, // 允许跨站请求
//				MaxAge:   3600,                  // 1 小时
//			}
//			c.SetCookie(cookie)
//			c.Response().Header().Set("X-CSRF-Token", csrfToken)
//		}
//		return next(c)
//	}
//}
