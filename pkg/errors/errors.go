package errors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

// OpenAIErrorFormat represents the standard error response matching OpenAI API
type OpenAIErrorFormat struct {
	Error OpenAIErrorDetail `json:"error"`
}

type OpenAIErrorDetail struct {
	Message string  `json:"message"`
	Type    string  `json:"type"`
	Param   *string `json:"param"`
	Code    *string `json:"code"`
}

// APIError represents a custom error that can be converted to HTTP status and OpenAI error format
type APIError struct {
	HTTPCode int
	Message  string
	Type     string
	Code     string
}

func (e *APIError) Error() string {
	return e.Message
}

// Common pre-defined errors
var (
	ErrUnauthorized    = &APIError{HTTPCode: http.StatusUnauthorized, Message: "Invalid API Key provided", Type: "invalid_request_error", Code: "invalid_api_key"}
	ErrRateLimit       = &APIError{HTTPCode: http.StatusTooManyRequests, Message: "Rate limit exceeded", Type: "requests", Code: "rate_limit_exceeded"}
	ErrQuotaExceeded   = &APIError{HTTPCode: http.StatusPaymentRequired, Message: "You exceeded your current quota, please check your plan and billing details.", Type: "insufficient_quota", Code: "insufficient_quota"}
	ErrInternal        = &APIError{HTTPCode: http.StatusInternalServerError, Message: "Internal server error while processing the request", Type: "server_error", Code: "internal_error"}
	ErrModelNotFound   = &APIError{HTTPCode: http.StatusNotFound, Message: "The model does not exist or you do not have access to it.", Type: "invalid_request_error", Code: "model_not_found"}
	ErrBackendUpstream = &APIError{HTTPCode: 502, Message: "The upstream AI provider returned an error.", Type: "server_error", Code: "bad_gateway"}
)

// New creates a custom API Error
func New(httpCode int, msg, errType, errCode string) *APIError {
	return &APIError{
		HTTPCode: httpCode,
		Message:  msg,
		Type:     errType,
		Code:     errCode,
	}
}

// EchoHTTPErrorHandler is a custom HTTP error handler for Echo
func EchoHTTPErrorHandler(c *echo.Context, err error) {
	if c.Response().(*echo.Response).Committed {
		return
	}

	var apiErr *APIError

	switch e := err.(type) {
	case *APIError:
		apiErr = e
	case *echo.HTTPError:
		msg := fmt.Sprintf("%v", e.Message)
		apiErr = &APIError{
			HTTPCode: e.Code,
			Message:  msg,
			Type:     "invalid_request_error",
		}
	default:
		msg := err.Error()
		// 敏感错误信息过滤屏蔽，避免泄漏底层技术栈（如 ent / sqlite / privacy）
		if strings.Contains(msg, "ent") || strings.Contains(msg, "privacy") || strings.Contains(msg, "sql") {
			msg = "Internal server error"
		}
		apiErr = &APIError{
			HTTPCode: http.StatusInternalServerError,
			Message:  msg,
			Type:     "server_error",
		}
	}

	// Always respond in OpenAI compatible JSON format
	res := OpenAIErrorFormat{
		Error: OpenAIErrorDetail{
			Message: apiErr.Message,
			Type:    apiErr.Type,
		},
	}

	if apiErr.Code != "" {
		res.Error.Code = &apiErr.Code
	}

	_ = c.JSON(apiErr.HTTPCode, res)
}
