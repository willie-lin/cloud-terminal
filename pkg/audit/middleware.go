package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

// Middleware creates an Echo middleware for audit logging
func Middleware(auditor Auditor) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// Skip health checks and non-critical endpoints
			if c.Path() == "/api/v1/health" || c.Path() == "/health" {
				return next(c)
			}
			// Only log mutations (POST/PUT/DELETE), skip GET page loads
			if c.Request().Method == "GET" {
				return next(c)
			}

			start := time.Now()
			requestID := c.Request().Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
				c.Response().Header().Set("X-Request-ID", requestID)
			}

			// Store request ID in context for handlers to use
			c.Set("request_id", requestID)

			endpoint := c.Request().Method + " " + c.Path()

			// 🔓 SECURITY: Propagate auditing details to standard Go context
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, "request_id", requestID)
			ctx = context.WithValue(ctx, "audit_endpoint", endpoint)
			c.SetRequest(c.Request().WithContext(ctx))

			// Pre-capture body for logging (bind/read next will consume it)
			var reqBody interface{}
			if c.Request().Method == "POST" || c.Request().Method == "PUT" {
				reqBody = getRedactedBody(c)
			}

			// Call next handler
			err := next(c)

			// Record audit event
			latency := time.Since(start).Milliseconds()
			status := "success"
			errMsg := ""
			if err != nil {
				status = "failure"
				errMsg = err.Error()
			}

			statusCode := http.StatusOK
			if sr, ok := c.Response().(interface{ Status() int }); ok {
				statusCode = sr.Status()
			}

			// Log the request
			details := map[string]interface{}{
				"method":      c.Request().Method,
				"path":        c.Path(),
				"status_code": statusCode,
			}
			if reqBody != nil {
				details["body"] = reqBody
			}

			_ = auditor.Log(c.Request().Context(), Event{
				RequestID:    requestID,
				Action:       EventHTTPRequest,
				ActorID:      getActorID(c),
				ActorType:    getActorType(c),
				SourceIP:     c.RealIP(),
				UserAgent:    c.Request().UserAgent(),
				Endpoint:     endpoint,
				ResourceType: "http",
				Status:       status,
				Latency:      latency,
				Error:        errMsg,
				Details:      details,
			})

			return err
		}
	}
}

// getRedactedBody reads, restores, and redactes the request body
func getRedactedBody(c *echo.Context) interface{} {
	req := c.Request()
	if req.Body == nil {
		return nil
	}

	// Read body
	bodyBytes, _ := io.ReadAll(req.Body)
	// Restore body for next handler
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if len(bodyBytes) == 0 {
		return nil
	}

	// Try to parse as JSON
	var payload map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		// Not JSON, return truncated string if text
		str := string(bodyBytes)
		if len(str) > 512 {
			return str[:512] + "..."
		}
		return str
	}

	// Redact sensitive fields
	return redactMap(payload)
}

func redactMap(m map[string]interface{}) map[string]interface{} {
	redacted := make(map[string]interface{})
	for k, v := range m {
		key := strings.ToLower(k)
		if isSensitive(key) {
			redacted[k] = "[REDACTED]"
		} else {
			redacted[k] = redactValue(v)
		}
	}
	return redacted
}

func redactValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		return redactMap(val)
	case []interface{}:
		list := make([]interface{}, len(val))
		for i, item := range val {
			list[i] = redactValue(item)
		}
		return list
	}
	return v
}

func isSensitive(key string) bool {
	sensitiveKeywords := []string{
		"password", "secret", "token", "private_key", "api_key",
		"auth", "credential", "mnemonic", "seed",
	}
	for _, kw := range sensitiveKeywords {
		if strings.Contains(key, kw) {
			return true
		}
	}
	return false
}

func getActorID(c *echo.Context) string {
	if name, ok := c.Get("user_name").(string); ok && name != "" {
		return name
	}
	if uid, ok := c.Get("user_id").(string); ok && uid != "" {
		return uid
	}
	return "anonymous"
}

func getActorType(c *echo.Context) string {
	if role, ok := c.Get("user_role").(string); ok && role != "" {
		return role
	}
	if _, ok := c.Get("auth_subject").(string); ok {
		return "user"
	}
	return "anonymous"
}
