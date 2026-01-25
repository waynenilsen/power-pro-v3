// Package middleware provides HTTP middleware for the API.
package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
)

// Context keys for authentication data.
type contextKey string

const (
	// UserIDKey is the context key for the authenticated user ID.
	UserIDKey contextKey = "user_id"
	// IsAdminKey is the context key for the admin flag.
	IsAdminKey contextKey = "is_admin"
	// UserKey is the context key for the full user object.
	UserKey contextKey = "user"
)

// AuthUser represents a user in the context for middleware purposes.
// This is separate from the auth.User to avoid circular dependencies.
type AuthUser struct {
	ID      string
	Email   string
	Name    string
	IsAdmin bool
}

// SessionValidator validates session tokens and returns user information.
type SessionValidator interface {
	ValidateSession(ctx context.Context, token string) (*AuthUser, error)
}

// GetUserID retrieves the authenticated user ID from the request context.
func GetUserID(r *http.Request) string {
	if id, ok := r.Context().Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}

// UserIDFromContext retrieves the authenticated user ID from a context.
func UserIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
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

// IsAdminFromContext checks if the authenticated user is an admin from context.
func IsAdminFromContext(ctx context.Context) bool {
	if isAdmin, ok := ctx.Value(IsAdminKey).(bool); ok {
		return isAdmin
	}
	return false
}

// UserFromContext retrieves the full user from a context.
func UserFromContext(ctx context.Context) *AuthUser {
	if user, ok := ctx.Value(UserKey).(*AuthUser); ok {
		return user
	}
	return nil
}

// AuthConfig holds configuration for the authentication middleware.
type AuthConfig struct {
	// WriteJSON is a function to write JSON error responses.
	WriteError func(w http.ResponseWriter, status int, message string)
	// SessionValidator validates session tokens.
	SessionValidator SessionValidator
}

// isTestMode checks if the application is running in test mode.
func isTestMode() bool {
	return os.Getenv("POWERPRO_TEST_MODE") == "true"
}

// RequireAuth creates middleware that requires authentication.
// It extracts user ID from the Authorization header (Bearer token) and validates via the auth service.
// In test mode (POWERPRO_TEST_MODE=true), falls back to X-User-ID and X-Admin headers.
func RequireAuth(cfg AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Try Bearer token authentication first
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")
				if token != "" && cfg.SessionValidator != nil {
					user, err := cfg.SessionValidator.ValidateSession(ctx, token)
					if err != nil {
						log.Printf("AUTH: Invalid session token for %s: %v", r.URL.Path, err)
						cfg.WriteError(w, http.StatusUnauthorized, "Invalid or expired session")
						return
					}

					// Set user info in context
					ctx = context.WithValue(ctx, UserIDKey, user.ID)
					ctx = context.WithValue(ctx, IsAdminKey, user.IsAdmin)
					ctx = context.WithValue(ctx, UserKey, user)

					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			// Fall back to test mode headers if no valid Bearer token
			if isTestMode() || cfg.SessionValidator == nil {
				userID := r.Header.Get("X-User-ID")
				if userID == "" {
					// Also check Bearer token as user ID for backward compatibility in tests
					if strings.HasPrefix(authHeader, "Bearer ") {
						userID = strings.TrimPrefix(authHeader, "Bearer ")
					}
				}

				if userID != "" {
					isAdmin := r.Header.Get("X-Admin") == "true"

					ctx = context.WithValue(ctx, UserIDKey, userID)
					ctx = context.WithValue(ctx, IsAdminKey, isAdmin)
					// Create a minimal user for test mode
					ctx = context.WithValue(ctx, UserKey, &AuthUser{
						ID:      userID,
						IsAdmin: isAdmin,
					})

					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			log.Printf("AUTH: Unauthorized request to %s - no user context", r.URL.Path)
			cfg.WriteError(w, http.StatusUnauthorized, "Authentication required")
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

// ChainMiddleware chains multiple middleware functions together.
func ChainMiddleware(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}
