package web

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"triangleoflove/backend/internal/auth"
	"triangleoflove/backend/internal/domain"
)

type contextKey string

const callerKey contextKey = "caller"

// CallerAccount holds the identity claims extracted from a validated JWT.
type CallerAccount struct {
	ID   string
	Role domain.Role
}

// CallerFromContext returns the CallerAccount injected by Middleware.
func CallerFromContext(ctx context.Context) CallerAccount {
	v, _ := ctx.Value(callerKey).(CallerAccount)
	return v
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
}

// Middleware validates the Bearer token and injects a CallerAccount into context.
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
		caller := CallerAccount{
			ID:   claims.AccountID,
			Role: domain.Role(claims.Role),
		}
		ctx := context.WithValue(r.Context(), callerKey, caller)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
