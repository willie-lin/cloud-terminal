package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
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
		TokenLookup: "cookie:token",
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

// ValidateRefreshToken 创建一个函数来验证RefreshToken
func ValidateRefreshToken(refreshToken string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("refresh_secret"), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	return claims, nil
}
