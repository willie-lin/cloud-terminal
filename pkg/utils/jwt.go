package utils

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

//var jwtKey = []byte("your_secret_key")
//
//type Claims struct {
//	Username string `json:"username"`
//	jwt.StandardClaims
//}
//
//// GenerateToken 生成JWT
//func GenerateToken(username string) (string, error) {
//	expirationTime := time.Now().Add(24 * time.Hour)
//	claims := &Claims{
//		Username: username,
//		StandardClaims: jwt.StandardClaims{
//			ExpiresAt: expirationTime.Unix(),
//		},
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	tokenString, err := token.SignedString(jwtKey)
//
//	if err != nil {
//		return "", err
//	}
//
//	return tokenString, nil
//}

var refreshKey = []byte("your_refresh_key")

type RefreshClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// GenerateRefreshToken 生成Refresh Token
func GenerateRefreshToken(email string) (string, error) {
	expirationTime := time.Now().Add(24 * 7 * time.Hour) // Refresh token的有效期通常比access token长
	claims := &RefreshClaims{
		Email: email,
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

//func ValidateToken(myToken string) (*jwt.Token, error) {
//	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//		}
//		return jwtKey, nil
//	})
//
//	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//		fmt.Println(claims["foo"], claims["nbf"])
//	} else {
//		fmt.Println(err)
//	}
//
//	return token, err
//}
//
//func ValidateRefreshToken(myToken string) (*jwt.Token, error) {
//	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//		}
//		return refreshKey, nil
//	})
//
//	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//		fmt.Println(claims["username"], claims["exp"])
//	} else {
//		fmt.Println(err)
//	}
//
//	return token, err
//}

//生成和存储 RefreshToken

//func GenerateRefreshToken(email string) (*ent.RefreshTokens, error) {
//	//claims := &jwt.StandardClaims{
//	//	Subject:   email,
//	//	IssuedAt:  time.Now().Unix(),
//	//	ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
//	//}
//	//
//	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	//return token.SignedString([]byte("your-secret-key"))
//
//	// 将 RefreshToken 存储到数据库中...
//
//	return nil, nil
//}

// GenerateAccessToken 生成 AccessToken：
func GenerateAccessToken(email string) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   email,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key"))
}

// ValidateAccessToken 验证 AccessToken
func ValidateAccessToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

//刷新 AccessToken

func RefreshAccessToken(refreshTokenString string) (string, error) {
	// 从数据库中获取 RefreshToken...

	//if refreshToken.ExpiresAt.Before(time.Now()) {
	//	return "", errors.New("refresh token expired")
	//}
	//
	//// 生成新的 AccessToken...
	//newAccessToken, err := GenerateAccessToken(refreshToken.User)
	//if err != nil {
	//	return "", err
	//}
	//
	//return newAccessToken, nil
	return "qqq", nil
}
