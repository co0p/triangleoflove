package web

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"triangleoflove/backend/internal/auth"
)

type contextKey string

// AccountIDKey is the context key for the authenticated account ID.
const AccountIDKey contextKey = "accountID"

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
}

// Middleware validates the Bearer token and injects the account ID into context.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			writeUnauthorized(w)
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := auth.ParseToken(tokenStr)
		if err != nil {
			writeUnauthorized(w)
			return
		}
		ctx := context.WithValue(r.Context(), AccountIDKey, claims.AccountID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
