package ctxutil

import "context"

type ContextKey string

const (
	UserIDCtxKey ContextKey = "user-id"
	TokenCtxKey  ContextKey = "token"
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

func ContextWithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, TokenCtxKey, token)
}

func TokenFromContext(ctx context.Context) string {
	token, ok := ctx.Value(TokenCtxKey).(string)
	if !ok {
		return ""
	}
	return token
}
