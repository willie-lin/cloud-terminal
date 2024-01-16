package utils

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// CreateJWTConfig 创建一个函数来生成JWT中间件的配置
func CreateJWTConfig() echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey:  []byte("secret"),
		TokenLookup: "cookie:token",
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.ErrUnauthorized
		},
	}
}

//func CreateJWTConfigs() middleware.JWTConfig {
//	return middleware.JWTConfig{
//		Claims:     &jwtCustomClaims{},
//		SigningKey: []byte("secret"),
//		Extractor:  middleware.JWTFromCookie("token"), // 从cookie中提取token
//		ErrorHandler: func(err error) error {
//			return echo.ErrUnauthorized
//		},
//	}
//}
//
//// extractTokenFromCookie 是一个自定义的函数，用于从cookie中提取token
//func extractTokenFromCookie(c echo.Context) (string, error) {
//	cookie, err := c.Cookie("token")
//	if err != nil {
//		if errors.Is(err, http.ErrNoCookie) {
//			return "", echo.ErrUnauthorized
//		}
//		return "", echo.ErrInternalServerError
//	}
//	return cookie.Value, nil
//}
