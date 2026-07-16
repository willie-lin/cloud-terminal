package audit

import (
	"context"
	"sync"
	"time"
)

// InMemoryAuditor stores audit events in memory for the admin dashboard
type InMemoryAuditor struct {
	mu     sync.RWMutex
	events []Event
	cap    int
}

func NewInMemoryAuditor(capacity int) *InMemoryAuditor {
	return &InMemoryAuditor{events: make([]Event, 0, capacity), cap: capacity}
}

func (m *InMemoryAuditor) Log(ctx context.Context, event Event) error {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	m.mu.Lock()
	m.events = append(m.events, event)
	if len(m.events) > m.cap {
		m.events = m.events[len(m.events)-m.cap:]
	}
	m.mu.Unlock()
	return nil
}

// Recent returns the most recent n events (newest first)
func (m *InMemoryAuditor) Recent(n int) []Event {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if n > len(m.events) {
		n = len(m.events)
	}
	result := make([]Event, n)
	for i := 0; i < n; i++ {
		result[i] = m.events[len(m.events)-1-i]
	}
	return result
}

// PaginatedResult holds a paginated slice of events plus total count
type PaginatedResult struct {
	Data     []Event `json:"data"`
	Total    int     `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
}

// ListPaginated returns events with pagination (page starts at 1, pageSize defaults to 20)
func (m *InMemoryAuditor) ListPaginated(page, pageSize int) PaginatedResult {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := len(m.events)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	// Calculate slice bounds (events are stored oldest-first; we want newest-first)
	start := total - (page-1)*pageSize - pageSize
	end := total - (page-1)*pageSize

	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start >= end {
		return PaginatedResult{
			Data:     []Event{},
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		}
	}

	// Extract slice and reverse (newest first)
	slice := m.events[start:end]
	result := make([]Event, len(slice))
	for i, evt := range slice {
		result[len(slice)-1-i] = evt
	}

	return PaginatedResult{
		Data:     result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

func (m *InMemoryAuditor) Close() error { return nil }
