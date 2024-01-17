package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// JwtCustomClaims 在全局范围内定义你的jwtCustomClaims类型
type JwtCustomClaims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// CreateAccessToken 创建一个函数来生成JWT
func CreateAccessToken(email, username string) (string, error) {
	claims := &JwtCustomClaims{
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

// ValidAccessTokenConfig valid access token
func ValidAccessTokenConfig() echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey:  []byte("secret"),
		TokenLookup: "cookie:AccessToken",
	}
}

// CreateRefreshToken 创建一个函数来生成RefreshToken
func CreateRefreshToken(email, username string) (string, error) {
	claims := &JwtCustomClaims{
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 144)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("refresh_secret"))
}

// ValidateRefreshTokenConfig  创建一个函数来验证RefreshToken
func ValidateRefreshTokenConfig() echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey:  []byte("refresh_secret"),
		TokenLookup: "cookie:RefreshToken",
	}
}

// CheckAccessToken Middleware for checking the AccessToken in the cookie
func CheckAccessTokens(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken, err := c.Cookie("AccessToken")
		fmt.Println(accessToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing access token")
		}

		// Validate the AccessToken
		accessTokenConfig := ValidAccessTokenConfig()
		acctoken, err := jwt.ParseWithClaims(accessToken.Value, accessTokenConfig.NewClaimsFunc(c), func(token *jwt.Token) (interface{}, error) {
			return accessTokenConfig.SigningKey, nil
		})
		fmt.Println(acctoken)

		fmt.Println(11111111)
		if err != nil || !acctoken.Valid {
			// If the AccessToken is invalid, use the RefreshToken to generate a new AccessToken
			refreshToken, err := c.Cookie("RefreshToken")
			fmt.Println(refreshToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing refresh token")
			}

			fmt.Println(222222)
			// Validate the RefreshToken
			refreshTokenConfig := ValidateRefreshTokenConfig()
			newtoken, err := jwt.ParseWithClaims(refreshToken.Value, refreshTokenConfig.NewClaimsFunc(c), func(token *jwt.Token) (interface{}, error) {
				return refreshTokenConfig.SigningKey, nil
			})
			fmt.Println(newtoken)
			fmt.Println(333333)

			if err == nil && newtoken.Valid {
				if claims, ok := newtoken.Claims.(*JwtCustomClaims); ok {
					// Use the RefreshToken's claims to generate a new AccessToken
					newAccessToken, err := CreateAccessToken(claims.Email, claims.Username)
					fmt.Println(newAccessToken)
					if err != nil {
						return err
					}

					// Set the new AccessToken in the cookie
					c.SetCookie(&http.Cookie{
						Name:     "AccessToken",
						Value:    newAccessToken,
						Expires:  time.Now().Add(time.Hour * 1), // The cookie will expire in 1 hour
						SameSite: http.SameSiteNoneMode,
						HttpOnly: true, // The cookie is not accessible via JavaScript
						Secure:   true, // The cookie is sent only over HTTPS
						Path:     "/",  // The cookie is available within the entire domain
					})
				}
			} else {
				// If the RefreshToken is invalid, return an error and prompt the user to log in again
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token, please log in again")
			}
		}
		return next(c)
	}
}

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
					newAccessToken, err := CreateAccessToken(claims.Email, claims.Username)
					if err != nil {
						return err
					}

					// Set the new AccessToken in the cookie
					c.SetCookie(&http.Cookie{
						Name:     "AccessToken",
						Value:    newAccessToken,
						Expires:  time.Now().Add(time.Hour * 1), // The cookie will expire in 1 hour
						SameSite: http.SameSiteNoneMode,
						HttpOnly: true, // The cookie is not accessible via JavaScript
						Secure:   true, // The cookie is sent only over HTTPS
						Path:     "/",  // The cookie is available within the entire domain
					})

					// Send the HTTP status code
					c.Response().WriteHeader(http.StatusOK)
				}
			} else {
				// If the RefreshToken is invalid, return an error and prompt the user to log in again
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token, please log in again")
			}
		}

		return next(c)
	}
}
