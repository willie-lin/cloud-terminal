package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

var refreshKey = []byte("your_refresh_key")

type RefreshClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateRefreshToken 生成Refresh Token
func GenerateRefreshToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * 7 * time.Hour) // Refresh token的有效期通常比access token长
	claims := &RefreshClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(refreshKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
