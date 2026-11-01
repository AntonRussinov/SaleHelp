package auth_google

import (
	"context"
	"net/http"
)

type contextKey int

const (
	contextKeyUserID contextKey = iota + 1
	contextKeyEmail
)

// SetUserIDToContext сохраняет user_id в контексте (после проверки JWT).
func SetUserIDToContext(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, contextKeyUserID, userID)
}

// UserIDFromContext возвращает user_id из контекста (0 если не задан).
func UserIDFromContext(ctx context.Context) int {
	v, _ := ctx.Value(contextKeyUserID).(int)
	return v
}

// SetEmailToContext сохраняет email в контексте.
func SetEmailToContext(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, contextKeyEmail, email)
}

// EmailFromContext возвращает email из контекста.
func EmailFromContext(ctx context.Context) string {
	v, _ := ctx.Value(contextKeyEmail).(string)
	return v
}

// SetJWTContext обновляет контекст запроса после проверки JWT (для использования в middleware).
func SetJWTContext(r *http.Request, userID int, email string) *http.Request {
	ctx := SetUserIDToContext(r.Context(), userID)
	ctx = SetEmailToContext(ctx, email)
	return r.WithContext(ctx)
}
