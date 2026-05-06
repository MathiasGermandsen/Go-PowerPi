package middleware

import (
	"context"
	"net/http"
	"strings"

	"Power-Pi/auth"
	"Power-Pi/database"
)

type contextKey string

const ClaimsContextKey contextKey = "claims"

func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing Authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			claims, err := auth.ValidateToken(parts[1], secret)
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			var revoked database.RevokedToken
			if result := database.DB.Where("jti = ?", claims.ID).First(&revoked); result.Error == nil {
				http.Error(w, "token has been revoked", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireScope(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsContextKey).(*auth.Claims)
			if !ok || claims == nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			for _, s := range claims.Scopes {
				if s == scope {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "insufficient scope", http.StatusForbidden)
		})
	}
}
