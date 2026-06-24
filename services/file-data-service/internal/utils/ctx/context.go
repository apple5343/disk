package ctxutil

import "context"

type ContextKey string

const (
	UserIDCtxKey ContextKey = "user-id"
)

func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDCtxKey, userID)
}

func UserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDCtxKey).(string)
	if !ok {
		return ""
	}
	return userID
}
