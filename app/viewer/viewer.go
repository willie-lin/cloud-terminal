package viewer

import (
	"context"
	"github.com/google/uuid"
)

type Viewer struct {
	UserID    uuid.UUID
	TenantID  uuid.UUID
	AccountID uuid.UUID
	RoleName  string // 使用切片存储多个角色
}

type viewerKey struct{}

func FromContext(ctx context.Context) *Viewer {
	v, _ := ctx.Value(viewerKey{}).(*Viewer)
	return v
}

func NewContext(ctx context.Context, v *Viewer) context.Context {
	return context.WithValue(ctx, viewerKey{}, v)
}
