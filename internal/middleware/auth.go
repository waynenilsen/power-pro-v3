// Package middleware provides HTTP middleware for the API.
package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
)

// Context keys for authentication data.
type contextKey string

const (
	// UserIDKey is the context key for the authenticated user ID.
	UserIDKey contextKey = "user_id"
	// IsAdminKey is the context key for the admin flag.
	IsAdminKey contextKey = "is_admin"
)

// GetUserID retrieves the authenticated user ID from the request context.
func GetUserID(r *http.Request) string {
	if id, ok := r.Context().Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}

// IsAdmin checks if the authenticated user is an admin.
func IsAdmin(r *http.Request) bool {
	if isAdmin, ok := r.Context().Value(IsAdminKey).(bool); ok {
		return isAdmin
	}
	return false
}

// AuthConfig holds configuration for the authentication middleware.
type AuthConfig struct {
	// WriteJSON is a function to write JSON error responses.
	WriteError func(w http.ResponseWriter, status int, message string)
}

// RequireAuth creates middleware that requires authentication.
// It extracts user ID from the Authorization header (Bearer token) or X-User-ID header.
// The X-Admin header is used to indicate admin status.
// This is temporary until a proper session-based auth system is implemented.
func RequireAuth(cfg AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := extractUserID(r)
			if userID == "" {
				log.Printf("AUTH: Unauthorized request to %s - no user context", r.URL.Path)
				cfg.WriteError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			isAdmin := r.Header.Get("X-Admin") == "true"

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, IsAdminKey, isAdmin)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdmin creates middleware that requires admin privileges.
// It must be used after RequireAuth middleware.
func RequireAdmin(cfg AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserID(r)
			if userID == "" {
				log.Printf("AUTH: Unauthorized request to %s - no user context", r.URL.Path)
				cfg.WriteError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			if !IsAdmin(r) {
				log.Printf("AUTH: Forbidden request to %s by user %s - admin required", r.URL.Path, userID)
				cfg.WriteError(w, http.StatusForbidden, "Admin privileges required")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ResourceOwnerFunc is a function type that extracts the owner ID from a request.
type ResourceOwnerFunc func(r *http.Request) (ownerID string, err error)

// RequireOwnerOrAdmin creates middleware that requires either resource ownership or admin privileges.
// It uses the provided function to determine the resource owner.
// If ownerFunc returns an empty string, it means the ownership check should be skipped (e.g., during resource creation).
func RequireOwnerOrAdmin(cfg AuthConfig, ownerFunc ResourceOwnerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserID(r)
			if userID == "" {
				log.Printf("AUTH: Unauthorized request to %s - no user context", r.URL.Path)
				cfg.WriteError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			// Admin can access any resource
			if IsAdmin(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Check ownership
			ownerID, err := ownerFunc(r)
			if err != nil {
				log.Printf("AUTH: Error checking ownership for %s: %v", r.URL.Path, err)
				cfg.WriteError(w, http.StatusInternalServerError, "Failed to verify resource ownership")
				return
			}

			// If ownerID is empty, skip ownership check (e.g., list operations where ownership is already in path)
			if ownerID == "" {
				next.ServeHTTP(w, r)
				return
			}

			if ownerID != userID {
				log.Printf("AUTH: Forbidden request to %s by user %s - not resource owner (owner: %s)", r.URL.Path, userID, ownerID)
				cfg.WriteError(w, http.StatusForbidden, "Access denied: you do not have permission to access this resource")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// extractUserID extracts the user ID from the request.
// It checks the Authorization header first (Bearer token), then X-User-ID header.
// This is temporary until proper session-based auth is implemented.
func extractUserID(r *http.Request) string {
	// Check Authorization header (Bearer token)
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		// For now, treat the bearer token as the user ID
		// This will be replaced with session lookup later
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != "" {
			return token
		}
	}

	// Fall back to X-User-ID header (for testing/development)
	return r.Header.Get("X-User-ID")
}

// ChainMiddleware chains multiple middleware functions together.
func ChainMiddleware(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}
