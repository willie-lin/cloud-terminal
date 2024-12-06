package middleware

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/app/database/ent"
	"net/http"
	"strings"
)

// AuthMiddleware 检查 JWT 并从中提取用户信息

//func AuthMiddleware(client *ent.Client) echo.MiddlewareFunc {
//	return func(next echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) error {
//			authHeader := c.Request().Header.Get("Authorization")
//			if authHeader == "" {
//				return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
//			}
//
//			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
//			token, err := jwt.ParseWithClaims(tokenString, &utils.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//				return []byte(utils.AccessTokenSecret), nil
//			})
//
//			if err != nil || !token.Valid {
//				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
//			}
//
//			claims, ok := token.Claims.(*utils.JwtCustomClaims)
//			if !ok || claims.Email == "" {
//				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
//			}
//
//			// 假设 UserID 本身就是 UUID 类型，不需要转换
//			//userID := uuid.MustParse(claims.UserID)
//
//			user, err := client.User.Get(context.Background(), claims.UserID)
//			if err != nil {
//				return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
//			}
//
//			c.Set("user", user)
//			return next(c)
//		}
//	}
//}

func RoleMiddleware(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*ent.User)

			userRoles, err := user.QueryRoles().All(context.Background())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user roles")
			}

			for _, userRole := range userRoles {
				for _, role := range roles {
					if userRole.Name == role {
						return next(c)
					}
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
		}
	}
}

func PermissionMiddleware(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*ent.User)

			roles, err := user.QueryRoles().All(context.Background())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user roles")
			}

			for _, role := range roles {
				permissions, err := role.QueryPermissions().All(context.Background())
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch role permissions")
				}

				for _, p := range permissions {
					if p.Name == permission {
						return next(c)
					}
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
		}
	}
}

//

func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// 这里应该返回你的JWT密钥
				return []byte("your-secret-key"), nil
			})

			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			claims := token.Claims.(jwt.MapClaims)
			c.Set("user", claims["username"])
			return next(c)
		}
	}
}

func Authorize(enforcer *casbin.Enforcer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(string)     // Assume user is set in context after authentication
			domain := c.Get("domain").(string) // Assume domain is set in context
			path := c.Request().URL.Path
			method := c.Request().Method

			allowed, err := enforcer.Enforce(user, domain, path, method)
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
