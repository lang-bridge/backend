package metactx

import (
	"context"
	"platform/internal/pkg/types"
)

type ctxKey uint

const (
	userIDKey ctxKey = iota + 1
)

func WithUserID(ctx context.Context, userID types.UserID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserID(ctx context.Context) types.UserID {
	value, ok := ctx.Value(userIDKey).(types.UserID)
	if !ok {
		panic("userID is absent")
	}
	return value
}
