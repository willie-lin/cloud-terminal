package viewer

import (
	"context"
)

type Viewer struct {
	UserID   string
	TenantID string
	Roles    []string
	Admin    bool
}

func FromContext(ctx context.Context) *Viewer {
	view, _ := ctx.Value(viewerContextKey).(*Viewer)
	return view
}

func (v *Viewer) HasRole(role string) bool {
	for _, r := range v.Roles {
		if r == role {
			return true
		}
	}
	return false
}

type viewerKey string

const viewerContextKey viewerKey = "viewer"

// NewContext returns a new context with the given viewer.
func NewContext(ctx context.Context, v *Viewer) context.Context {
	return context.WithValue(ctx, viewerContextKey, v)
}
