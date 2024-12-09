package middlewares

import (
	"context"
	"fmt"
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

func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse token claims")
			}

			c.Set("user", claims["username"])
			c.Set("roles", claims["roles"])

			return next(c)
		}
	}
}

//func Authorize(enforcer *casbin.Enforcer) echo.MiddlewareFunc {
//	return func(next echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) error {
//			user := c.Get("user").(string) // Assume user is set in context after authentication
//			fmt.Println(user)
//			domain := c.Get("domain").(string) // Assume domain is set in context
//			fmt.Println(domain)
//			path := c.Request().URL.Path
//			fmt.Println(path)
//			method := c.Request().Method
//			fmt.Println(method)
//
//			allowed, err := enforcer.Enforce(user, domain, path, method)
//			if err != nil {
//				return echo.NewHTTPError(http.StatusInternalServerError, "Error checking permissions")
//			}
//			if !allowed {
//				return echo.NewHTTPError(http.StatusForbidden, "Permission denied")
//			}
//			return next(c)
//		}
//	}
//}

func Authorize(enforcer *casbin.Enforcer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("Authorize middleware is running")

			user, ok := c.Get("user").(string)
			if !ok {
				fmt.Println("User not found in context")
				return echo.NewHTTPError(http.StatusUnauthorized, "用户未认证")
			}
			fmt.Println("User:", user)

			domain, ok := c.Get("domain").(string)
			if !ok {
				fmt.Println("Domain not found in context")
				return echo.NewHTTPError(http.StatusUnauthorized, "域未设置")
			}
			fmt.Println("Domain:", domain)

			path := c.Request().URL.Path
			fmt.Println("Path:", path)

			method := c.Request().Method
			fmt.Println("Method:", method)

			allowed, err := enforcer.Enforce(user, domain, path, method)
			if err != nil {
				fmt.Println("Error enforcing policy:", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "检查权限时出错")
			}
			if !allowed {
				fmt.Println("Permission denied")
				return echo.NewHTTPError(http.StatusForbidden, "没有权限")
			}

			fmt.Println("Permission granted")
			return next(c)
		}
	}
}
