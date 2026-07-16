package audit

import (
	"context"
	"fmt"
	"strings"
)

// ExtractActor extracts the actor from the context
func ExtractActor(ctx context.Context) string {
	if val := ctx.Value("auth_subject"); val != nil {
		if s, ok := val.(string); ok && s != "" {
			return s
		}
	}
	// Also check general "user" key (e.g. from JWT middleware standard)
	if val := ctx.Value("user"); val != nil {
		// handle *jwt.Token or string
		return "authenticated_user" // simplified, usually need type assertion
	}
	return "unknown" // Use "unknown" for generic non-auth, "system" reserved for internal
}

// ExtractRequestID extracts the request ID from the context
func ExtractRequestID(ctx context.Context) string {
	if val := ctx.Value("request_id"); val != nil {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

// ExtractEndpoint extracts the endpoint name from the context
func ExtractEndpoint(ctx context.Context) string {
	if val := ctx.Value("audit_endpoint"); val != nil {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

// MultiAuditor 分发审计日志到多个 Auditor
type MultiAuditor struct {
	auditors []Auditor
}

func NewMultiAuditor(auditors ...Auditor) *MultiAuditor {
	return &MultiAuditor{auditors: auditors}
}

func (m *MultiAuditor) Log(ctx context.Context, event Event) error {
	var errs []string
	for _, a := range m.auditors {
		if err := a.Log(ctx, event); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("audit log errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

func (m *MultiAuditor) Close() error {
	var errs []string
	for _, a := range m.auditors {
		if err := a.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("audit close errors: %s", strings.Join(errs, "; "))
	}
	return nil
}
