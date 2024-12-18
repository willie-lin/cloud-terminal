package middlewarers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/willie-lin/cloud-terminal/app/viewer"
	"github.com/willie-lin/cloud-terminal/pkg/utils"
	"log"
	"net/http"
	"time"
)

// AuthenticateAndAuthorize 中间件用于验证和授权用户
func AuthenticateAndAuthorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken, err := c.Cookie("AccessToken")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing access token")
		}
		// 验证访问令牌
		config := utils.ValidAccessTokenConfig()
		token, err := jwt.ParseWithClaims(accessToken.Value, config.NewClaimsFunc(c), func(token *jwt.Token) (interface{}, error) {
			return config.SigningKey, nil
		})

		if err != nil || !token.Valid {
			// 如果访问令牌无效，使用刷新令牌生成新的访问令牌
			refreshToken, err := c.Cookie("RefreshToken")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing refresh token")
			}

			// 验证刷新令牌
			config = utils.ValidateRefreshTokenConfig()
			token, err = jwt.ParseWithClaims(refreshToken.Value, config.NewClaimsFunc(c), func(token *jwt.Token) (interface{}, error) {
				return config.SigningKey, nil
			})

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(*utils.JwtCustomClaims); ok {
					// 使用刷新令牌的声明生成新的访问令牌
					newAccessToken, err := utils.CreateAccessToken(claims.UserID, claims.TenantID, claims.Email, claims.Username, claims.RoleName)
					if err != nil {
						return err
					}
					// 将新的访问令牌设置在Cookie中
					c.SetCookie(&http.Cookie{
						Name:     "AccessToken",
						Value:    newAccessToken,
						Expires:  time.Now().Add(time.Hour * 1), // The cookie will expire in 1 hour
						SameSite: http.SameSiteNoneMode,
						HttpOnly: true, // The cookie is not accessible via JavaScript
						Secure:   true, // The cookie is sent only over HTTPS
						Path:     "/",  // The cookie is available within the entire domain
					})
					return c.JSON(http.StatusOK, map[string]string{"message": "Token refreshed successfully"})
				}
			} else {
				// If the RefreshToken is invalid, return an error and prompt the user to log in again
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token, please log in again")
			}
		}
		// 提取租户ID并设置到请求上下文中
		if claims, ok := token.Claims.(*utils.JwtCustomClaims); ok && token.Valid {
			v := &viewer.Viewer{
				TenantID: claims.TenantID,
				UserID:   claims.UserID,
				RoleName: claims.RoleName,
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
