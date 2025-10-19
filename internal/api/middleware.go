package api

import (
	"ai-assistant/internal/service"
	"context"
	"net/http"
	"strconv"
	"strings"
)

type contextKey string
const UserIDKey contextKey = "userID"

func AuthMiddleware(authSvc service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := authSvc.ParseToken(r.Context(), tokenString)
			if err != nil {
				http.Error(w, "Invalid Token", http.StatusUnauthorized)
				return
			}

			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				http.Error(w, "Invalid User ID", http.StatusUnauthorized)
				return
			}
			
			// Кладём ID юзера в контекст
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}