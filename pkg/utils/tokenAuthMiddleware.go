package utils

import "github.com/labstack/echo/v4"

func TokenAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 获取token
			token := c.Request().Header.Get("Authorization")
			// 验证token，这里只是简单的示例，你应该使用你自己的验证逻辑
			if token != "your_token" {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}
