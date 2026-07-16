package audit

import (
	"context"
	"time"
)

// EventType 定义事件类型
type EventType string

const (
	// System Events
	EventSystemStartup  EventType = "system.startup"
	EventSystemShutdown EventType = "system.shutdown"

	// HTTP Events
	EventHTTPRequest EventType = "http.request"
)

// Event represents a single audit log entry (5W1H)
type Event struct {
	ID        string    `json:"id"`        // Unique Event ID
	Timestamp time.Time `json:"timestamp"` // When

	// Who
	ActorID   string `json:"actor_id"`   // User ID
	ActorType string `json:"actor_type"` // user, system, anonymous
	SourceIP  string `json:"source_ip"`  // Source IP address
	UserAgent string `json:"user_agent"` // Client User Agent

	// Should be where
	ResourceID   string `json:"resource_id"`   // Resource ID (e.g. holding_id)
	ResourceType string `json:"resource_type"` // holding, order, user_config
	Endpoint     string `json:"endpoint"`      // API Endpoint

	// What
	Action  EventType              `json:"action"`  // Action Type
	Details map[string]interface{} `json:"details"` // Detailed parameters
	Status  string                 `json:"status"`  // success, failure

	// Why & How
	RequestID string `json:"request_id"` // Trace ID
	Latency   int64  `json:"latency_ms"` // Duration in milliseconds
	Error     string `json:"error,omitempty"`
}

// Auditor is the interface for recording audit events
type Auditor interface {
	Log(ctx context.Context, event Event) error
	Close() error
}
