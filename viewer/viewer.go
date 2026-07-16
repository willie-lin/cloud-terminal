package viewer

import (
	"context"
	"github.com/google/uuid"
)

// Viewer holds the current request context identity info.
type Viewer struct {
	UserID   uuid.UUID
	TenantID uuid.UUID
	GroupID  uuid.UUID
	RoleName string
}

type viewerKey struct{}

func FromContext(ctx context.Context) *Viewer {
	v, _ := ctx.Value(viewerKey{}).(*Viewer)
	return v
}

func NewContext(ctx context.Context, v *Viewer) context.Context {
	return context.WithValue(ctx, viewerKey{}, v)
}
