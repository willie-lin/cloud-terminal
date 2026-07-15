package middlewarers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"github.com/willie-lin/cloud-terminal/viewer"
	"log"
	"net/http"
	"strings"
	"time"
)

// AuthenticateAndAuthorize 中间件用于验证和授权用户
func AuthenticateAndAuthorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		token, err := parseAndValidateToken(c, "AccessToken", utils.AccessTokenSecret)
		if err != nil {
			// 尝试用刷新令牌
			token, err = parseAndValidateToken(c, "RefreshToken", utils.RefreshTokenSecret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or missing token")
			}
			// 刷新令牌有效，签发新的访问令牌
			if claims, ok := token.Claims.(*utils.JwtCustomClaims); ok {
				newAccessToken, err := utils.CreateAccessToken(claims.UserID, claims.TenantID, claims.AccountID, claims.Email, claims.Username, claims.RoleName)
				if err != nil {
					return err
				}
				c.SetCookie(&http.Cookie{
					Name:     "AccessToken",
					Value:    newAccessToken,
					Expires:  time.Now().Add(time.Hour * 1),
					SameSite: http.SameSiteNoneMode,
					HttpOnly: true,
					Secure:   true,
					Path:     "/",
				})
				return c.JSON(http.StatusOK, map[string]string{"message": "Token refreshed successfully"})
			}
		}

		if claims, ok := token.Claims.(*utils.JwtCustomClaims); ok && token.Valid {
			v := &viewer.Viewer{
				TenantID:  claims.TenantID,
				UserID:    claims.UserID,
				AccountID: claims.AccountID,
				RoleName:  strings.ToLower(claims.RoleName),
			}
			ctx := viewer.NewContext(c.Request().Context(), v)
			c.SetRequest(c.Request().WithContext(ctx))
		} else {
			log.Println("Invalid token claims or token not valid")
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid access token")
		}
		return next(c)
	}
}

func parseAndValidateToken(c *echo.Context, cookieName, secret string) (*jwt.Token, error) {
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		return nil, err
	}
	claims := &utils.JwtCustomClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token, nil
}
