package response

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

// Response 标准响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 返回成功响应 (HTTP 200)
func Success(c *echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// Error 返回错误响应
func Error(c *echo.Context, httpCode int, msg string) error {
	return c.JSON(httpCode, Response{
		Code: httpCode,
		Msg:  msg,
	})
}

// ServerError 返回服务器内部错误 (HTTP 500)
func ServerError(c *echo.Context, err error) error {
	msg := err.Error()
	// 敏感错误信息过滤屏蔽，避免泄漏底层技术栈
	if strings.Contains(msg, "ent") || strings.Contains(msg, "privacy") || strings.Contains(msg, "sql") {
		msg = "Internal server error"
	}
	return c.JSON(http.StatusInternalServerError, Response{
		Code: http.StatusInternalServerError,
		Msg:  msg,
	})
}
