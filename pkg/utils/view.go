package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/app/viewer"
	"log"
	"net/http"
)

func WithViewer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 从cookie中获取token
		accessToken, err := c.Cookie("AccessToken")
		if err != nil {
			log.Printf("Missing token: %v", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
		}

		log.Printf("AccessToken=%s", accessToken.Value)

		config := ValidAccessTokenConfig()
		token, err := jwt.ParseWithClaims(accessToken.Value, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return config.SigningKey, nil
		})

		if err != nil {
			log.Printf("Error parsing token: %v", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		if !token.Valid {
			log.Printf("Token is not valid")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		fmt.Println("zzzzzzzzzz") // Debug log to confirm execution point

		claims, ok := token.Claims.(*JwtCustomClaims)
		if !ok {
			log.Printf("Invalid token claims")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
		}
		fmt.Printf("Token claims: Email=%s, UserID=%s, TenantID=%s, RoleID=%s\n", claims.Email, claims.UserID, claims.TenantID, claims.RoleName)

		view := &viewer.Viewer{
			UserID:   claims.UserID,
			TenantID: claims.TenantID,
			//RoleID:   claims.RoleID,
			//Admin:    claims.RoleID.String() == "admin" || claims.RoleID.String() == "superAdmin",
		}

		ctx := viewer.NewContext(c.Request().Context(), view)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
