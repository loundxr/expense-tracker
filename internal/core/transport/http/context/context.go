package core_http_context

import "context"

type contextKey string

const (
	userIDKey   contextKey = "user_id"
	userRoleKey contextKey = "user_role"
)

func SetUserID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func GetUserID(ctx context.Context) int {
	uid, ok := ctx.Value(userIDKey).(int)
	if !ok {
		return 0
	}
	return uid
}

func SetUserRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, userRoleKey, role)
}

func GetUserRole(ctx context.Context) string {
	role := ctx.Value(userRoleKey).(string)
	return role
}
