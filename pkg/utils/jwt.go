package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// JwtCustomClaims 在全局范围内定义你的jwtCustomClaims类型
type JwtCustomClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

// 在代码中直接定义密钥（使用随机密钥）
var (
	AccessTokenSecret  string
	RefreshTokenSecret string
)

func init() {
	var err error
	AccessTokenSecret, err = GenerateRandomKey(32) //32字节 = 256位
	if err != nil {
		panic(fmt.Sprintf("Failed to generate access token secret: %v", err))
	}
	RefreshTokenSecret, err = GenerateRandomKey(32)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate refresh token secret: %v", err))
	}
}

// CreateAccessToken 创建一个函数来生成JWT
func CreateAccessToken(userID uuid.UUID, email, username string) (string, error) {
	claims := &JwtCustomClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(AccessTokenSecret))
}

// ValidAccessTokenConfig valid access token
func ValidAccessTokenConfig() echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey:  []byte(AccessTokenSecret),
		TokenLookup: "cookie:AccessToken",
	}
}

// CreateRefreshToken 创建一个函数来生成RefreshToken
func CreateRefreshToken(userID uuid.UUID, email, username string) (string, error) {
	claims := &JwtCustomClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 144)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(RefreshTokenSecret))
}

// ValidateRefreshTokenConfig  创建一个函数来验证RefreshToken
func ValidateRefreshTokenConfig() echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey:  []byte(RefreshTokenSecret),
		TokenLookup: "cookie:RefreshToken",
	}
}

// CheckAccessToken Middleware for checking the AccessToken in the cookie
func CheckAccessToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		accessToken, err := c.Cookie("AccessToken")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing access token")
		}

		// Validate the AccessToken
		config := ValidAccessTokenConfig()
		token, err := jwt.ParseWithClaims(accessToken.Value, config.NewClaimsFunc(c), func(token *jwt.Token) (interface{}, error) {
			return config.SigningKey, nil
		})

		if err != nil || !token.Valid {
			// If the AccessToken is invalid, use the RefreshToken to generate a new AccessToken
			refreshToken, err := c.Cookie("RefreshToken")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing refresh token")
			}

			// Validate the RefreshToken
			config = ValidateRefreshTokenConfig()
			token, err = jwt.ParseWithClaims(refreshToken.Value, config.NewClaimsFunc(c), func(token *jwt.Token) (interface{}, error) {
				return config.SigningKey, nil
			})

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(*JwtCustomClaims); ok {
					// Use the RefreshToken's claims to generate a new AccessToken
					newAccessToken, err := CreateAccessToken(claims.UserID, claims.Email, claims.Username)
					if err != nil {
						return err
					}

					// Set the new AccessToken in the cookie
					c.SetCookie(&http.Cookie{
						Name:     "AccessToken",
						Value:    newAccessToken,
						Expires:  time.Now().Add(time.Hour * 1), // The cookie will expire in 1 hour
						SameSite: http.SameSiteNoneMode,
						HttpOnly: false, // The cookie is not accessible via JavaScript
						Secure:   true,  // The cookie is sent only over HTTPS
						Path:     "/",   // The cookie is available within the entire domain
					})

					// Send the HTTP status code
					//c.Response().WriteHeader(http.StatusOK)
					//return c.Redirect(http.StatusTemporaryRedirect, c.Request().URL.String())
					return c.JSON(http.StatusOK, map[string]string{"message": "Token refreshed successfully"})
				}
			} else {
				// If the RefreshToken is invalid, return an error and prompt the user to log in again
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token, please log in again")
			}
		}
		return next(c)
	}
}

//func SetCSRFToken(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		if c.Request().Method != http.MethodOptions {
//			token := c.Get(middleware.DefaultCSRFConfig.ContextKey)
//			if token == nil {
//				var err error
//				token, err = GenerateRandomKey(32)
//				if err != nil {
//					return echo.NewHTTPError(http.StatusInternalServerError, "生成 CSRF 令牌失败")
//				}
//			}
//			fmt.Println("6666666666666666666666")
//			csrfToken := token.(string)
//			fmt.Println("Generated CSRF Token:", csrfToken) // 调试输出
//			cookie := &http.Cookie{
//				Name:   "_csrf",
//				Value:  csrfToken,
//				Path:   "/",
//				Domain: c.Request().Host,
//				//Domain: "localhost",
//				//Secure:   c.IsTLS(),
//				Secure:   true,
//				HttpOnly: true, // 确保 HttpOnly 设置为 false 以允许前端访问
//				//SameSite: http.SameSiteLaxMode,
//				SameSite: http.SameSiteNoneMode,
//				MaxAge:   3600,
//			}
//
//			c.SetCookie(cookie)
//			c.Response().Header().Set("X-CSRF-Token", csrfToken)
//			fmt.Printf("CSRF 令牌: %s\n", csrfToken)
//
//			// 打印所有响应头
//			for name, values := range c.Response().Header() {
//				for _, value := range values {
//					fmt.Printf("%s: %s\n", name, value)
//				}
//			}
//		}
//		fmt.Println("777777777777777")
//
//		return next(c)
//	}
//}
//
