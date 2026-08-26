package middleware

import (
	"context"
	"liftwork/internal/security"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "user_id"

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerAuth := r.Header.Get("Authorization")

			parts := strings.Fields(headerAuth)

			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			userID, err := security.ValidateToken(token, secret)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				userIDKey,
				userID,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}
