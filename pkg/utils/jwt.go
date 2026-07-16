package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"net/http"
	"time"
)

// JwtCustomClaims JWT 自定义声明
type JwtCustomClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	TenantID uuid.UUID `json:"tenant_id"`
	RoleName string    `json:"role_name"`
	GroupID  uuid.UUID `json:"group_id"`
	jwt.RegisteredClaims
}

var (
	AccessTokenSecret  string
	RefreshTokenSecret string
)

func init() {
	var err error
	AccessTokenSecret, err = GenerateRandomKey(32)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate access token secret: %v", err))
	}
	RefreshTokenSecret, err = GenerateRandomKey(32)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate refresh token secret: %v", err))
	}
}

// CreateAccessToken 创建JWT访问令牌
func CreateAccessToken(userID, tenantID, groupID uuid.UUID, email, username, roleName string) (string, error) {
	claims := &JwtCustomClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		TenantID: tenantID,
		RoleName: roleName,
		GroupID:  groupID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(AccessTokenSecret))
}

// CreateRefreshToken 创建刷新令牌
func CreateRefreshToken(userID, tenantID, groupID uuid.UUID, email, username, roleName string) (string, error) {
	claims := &JwtCustomClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		TenantID: tenantID,
		RoleName: roleName,
		GroupID:  groupID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 144)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(RefreshTokenSecret))
}

// CheckAccessToken 中间件用于检查访问令牌
func CheckAccessToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		accessToken, err := c.Cookie("AccessToken")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing access token")
		}

		claims := &JwtCustomClaims{}
		token, err := jwt.ParseWithClaims(accessToken.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(AccessTokenSecret), nil
		})

		if err != nil || !token.Valid {
			refreshToken, err := c.Cookie("RefreshToken")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing refresh token")
			}

			refreshClaims := &JwtCustomClaims{}
			token, err = jwt.ParseWithClaims(refreshToken.Value, refreshClaims, func(token *jwt.Token) (interface{}, error) {
				return []byte(RefreshTokenSecret), nil
			})

			if err == nil && token.Valid {
				newAccessToken, err := CreateAccessToken(refreshClaims.UserID, refreshClaims.TenantID, refreshClaims.GroupID, refreshClaims.Email, refreshClaims.Username, refreshClaims.RoleName)
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
			} else {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token, please log in again")
			}
		}

		if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
			c.Set("user_id", claims.UserID)
			c.Set("tenant_id", claims.TenantID)
		}
		return next(c)
	}
}
