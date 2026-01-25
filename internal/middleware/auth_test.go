package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockSessionValidator is a mock implementation of SessionValidator for testing.
type mockSessionValidator struct {
	users map[string]*AuthUser
	err   error
}

func newMockSessionValidator() *mockSessionValidator {
	return &mockSessionValidator{
		users: make(map[string]*AuthUser),
	}
}

func (m *mockSessionValidator) ValidateSession(ctx context.Context, token string) (*AuthUser, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, ok := m.users[token]
	if !ok {
		return nil, errors.New("invalid session")
	}
	return user, nil
}

func (m *mockSessionValidator) addUser(token string, user *AuthUser) {
	m.users[token] = user
}

// mockErrorWriter tracks error responses.
type mockErrorWriter struct {
	status  int
	message string
}

func (m *mockErrorWriter) writeError(w http.ResponseWriter, status int, message string) {
	m.status = status
	m.message = message
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func TestRequireAuth_WithSessionValidator(t *testing.T) {
	validator := newMockSessionValidator()
	validator.addUser("valid-token", &AuthUser{
		ID:      "user-123",
		Email:   "test@example.com",
		Name:    "Test User",
		IsAdmin: false,
	})
	validator.addUser("admin-token", &AuthUser{
		ID:      "admin-123",
		Email:   "admin@example.com",
		Name:    "Admin User",
		IsAdmin: true,
	})

	errWriter := &mockErrorWriter{}
	cfg := AuthConfig{
		WriteError:       errWriter.writeError,
		SessionValidator: validator,
	}

	tests := []struct {
		name           string
		authHeader     string
		wantStatus     int
		wantUserID     string
		wantAdmin      bool
		wantUserInCtx  bool
		testMode       bool
		xUserID        string
		xAdmin         string
	}{
		{
			name:          "valid bearer token authenticates user",
			authHeader:    "Bearer valid-token",
			wantStatus:    http.StatusOK,
			wantUserID:    "user-123",
			wantAdmin:     false,
			wantUserInCtx: true,
		},
		{
			name:          "valid bearer token authenticates admin",
			authHeader:    "Bearer admin-token",
			wantStatus:    http.StatusOK,
			wantUserID:    "admin-123",
			wantAdmin:     true,
			wantUserInCtx: true,
		},
		{
			name:       "invalid bearer token returns 401",
			authHeader: "Bearer invalid-token",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "empty bearer token returns 401",
			authHeader: "Bearer ",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "no auth header returns 401",
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "non-bearer auth header returns 401 (without test mode)",
			authHeader: "Basic something",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:          "test mode with X-User-ID header authenticates",
			authHeader:    "",
			xUserID:       "test-user",
			wantStatus:    http.StatusOK,
			wantUserID:    "test-user",
			wantAdmin:     false,
			wantUserInCtx: true,
			testMode:      true,
		},
		{
			name:          "test mode with X-User-ID and X-Admin headers",
			authHeader:    "",
			xUserID:       "test-admin",
			xAdmin:        "true",
			wantStatus:    http.StatusOK,
			wantUserID:    "test-admin",
			wantAdmin:     true,
			wantUserInCtx: true,
			testMode:      true,
		},
		{
			name:       "test mode without headers returns 401",
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
			testMode:   true,
		},
		{
			name:          "valid session takes precedence over test mode headers",
			authHeader:    "Bearer valid-token",
			xUserID:       "different-user",
			xAdmin:        "true",
			wantStatus:    http.StatusOK,
			wantUserID:    "user-123",
			wantAdmin:     false,
			wantUserInCtx: true,
			testMode:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set or unset test mode
			if tt.testMode {
				os.Setenv("POWERPRO_TEST_MODE", "true")
				defer os.Unsetenv("POWERPRO_TEST_MODE")
			} else {
				os.Unsetenv("POWERPRO_TEST_MODE")
			}

			errWriter.status = 0
			errWriter.message = ""

			var capturedUserID string
			var capturedAdmin bool
			var capturedUser *AuthUser

			handler := RequireAuth(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedUserID = GetUserID(r)
				capturedAdmin = IsAdmin(r)
				capturedUser = UserFromContext(r.Context())
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			if tt.xUserID != "" {
				req.Header.Set("X-User-ID", tt.xUserID)
			}
			if tt.xAdmin != "" {
				req.Header.Set("X-Admin", tt.xAdmin)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, tt.wantUserID, capturedUserID)
				assert.Equal(t, tt.wantAdmin, capturedAdmin)
				if tt.wantUserInCtx {
					require.NotNil(t, capturedUser)
					assert.Equal(t, tt.wantUserID, capturedUser.ID)
					assert.Equal(t, tt.wantAdmin, capturedUser.IsAdmin)
				}
			} else {
				assert.Equal(t, tt.wantStatus, rr.Code)
			}
		})
	}
}

func TestRequireAuth_WithoutSessionValidator(t *testing.T) {
	// When no SessionValidator is provided, middleware falls back to test mode behavior
	errWriter := &mockErrorWriter{}
	cfg := AuthConfig{
		WriteError:       errWriter.writeError,
		SessionValidator: nil,
	}

	tests := []struct {
		name       string
		authHeader string
		xUserID    string
		xAdmin     string
		wantStatus int
		wantUserID string
		wantAdmin  bool
	}{
		{
			name:       "X-User-ID header authenticates",
			xUserID:    "test-user",
			wantStatus: http.StatusOK,
			wantUserID: "test-user",
			wantAdmin:  false,
		},
		{
			name:       "Bearer token used as user ID",
			authHeader: "Bearer test-user-id",
			wantStatus: http.StatusOK,
			wantUserID: "test-user-id",
			wantAdmin:  false,
		},
		{
			name:       "X-User-ID takes precedence over Bearer",
			authHeader: "Bearer bearer-user",
			xUserID:    "header-user",
			wantStatus: http.StatusOK,
			wantUserID: "header-user",
			wantAdmin:  false,
		},
		{
			name:       "no auth returns 401",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errWriter.status = 0
			errWriter.message = ""

			var capturedUserID string
			var capturedAdmin bool

			handler := RequireAuth(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedUserID = GetUserID(r)
				capturedAdmin = IsAdmin(r)
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			if tt.xUserID != "" {
				req.Header.Set("X-User-ID", tt.xUserID)
			}
			if tt.xAdmin != "" {
				req.Header.Set("X-Admin", tt.xAdmin)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tt.wantStatus == http.StatusOK {
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, tt.wantUserID, capturedUserID)
				assert.Equal(t, tt.wantAdmin, capturedAdmin)
			} else {
				assert.Equal(t, tt.wantStatus, rr.Code)
			}
		})
	}
}

func TestRequireAuth_SessionValidatorError(t *testing.T) {
	validator := newMockSessionValidator()
	validator.err = errors.New("database error")

	errWriter := &mockErrorWriter{}
	cfg := AuthConfig{
		WriteError:       errWriter.writeError,
		SessionValidator: validator,
	}

	handler := RequireAuth(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer some-token")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Invalid or expired session", errWriter.message)
}

func TestRequireAdmin(t *testing.T) {
	errWriter := &mockErrorWriter{}
	cfg := AuthConfig{
		WriteError: errWriter.writeError,
	}

	tests := []struct {
		name       string
		userID     string
		isAdmin    bool
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "admin user allowed",
			userID:     "admin-123",
			isAdmin:    true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-admin user forbidden",
			userID:     "user-123",
			isAdmin:    false,
			wantStatus: http.StatusForbidden,
			wantMsg:    "Admin privileges required",
		},
		{
			name:       "no user unauthorized",
			userID:     "",
			isAdmin:    false,
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "Authentication required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errWriter.status = 0
			errWriter.message = ""

			handler := RequireAdmin(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			ctx := req.Context()
			if tt.userID != "" {
				ctx = context.WithValue(ctx, UserIDKey, tt.userID)
				ctx = context.WithValue(ctx, IsAdminKey, tt.isAdmin)
			}
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantMsg != "" {
				assert.Equal(t, tt.wantMsg, errWriter.message)
			}
		})
	}
}

func TestRequireOwnerOrAdmin(t *testing.T) {
	errWriter := &mockErrorWriter{}
	cfg := AuthConfig{
		WriteError: errWriter.writeError,
	}

	tests := []struct {
		name       string
		userID     string
		isAdmin    bool
		ownerID    string
		ownerErr   error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "owner can access",
			userID:     "user-123",
			isAdmin:    false,
			ownerID:    "user-123",
			wantStatus: http.StatusOK,
		},
		{
			name:       "admin can access any resource",
			userID:     "admin-123",
			isAdmin:    true,
			ownerID:    "user-456",
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-owner forbidden",
			userID:     "user-123",
			isAdmin:    false,
			ownerID:    "user-456",
			wantStatus: http.StatusForbidden,
			wantMsg:    "Access denied: you do not have permission to access this resource",
		},
		{
			name:       "empty owner ID skips check",
			userID:     "user-123",
			isAdmin:    false,
			ownerID:    "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "no user unauthorized",
			userID:     "",
			isAdmin:    false,
			wantStatus: http.StatusUnauthorized,
			wantMsg:    "Authentication required",
		},
		{
			name:       "owner check error",
			userID:     "user-123",
			isAdmin:    false,
			ownerErr:   errors.New("database error"),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "Failed to verify resource ownership",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errWriter.status = 0
			errWriter.message = ""

			ownerFunc := func(r *http.Request) (string, error) {
				return tt.ownerID, tt.ownerErr
			}

			handler := RequireOwnerOrAdmin(cfg, ownerFunc)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			ctx := req.Context()
			if tt.userID != "" {
				ctx = context.WithValue(ctx, UserIDKey, tt.userID)
				ctx = context.WithValue(ctx, IsAdminKey, tt.isAdmin)
			}
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantMsg != "" {
				assert.Equal(t, tt.wantMsg, errWriter.message)
			}
		})
	}
}

func TestContextHelpers(t *testing.T) {
	t.Run("GetUserID from request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), UserIDKey, "user-123")
		req = req.WithContext(ctx)

		assert.Equal(t, "user-123", GetUserID(req))
	})

	t.Run("GetUserID returns empty for no user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		assert.Equal(t, "", GetUserID(req))
	})

	t.Run("UserIDFromContext", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "user-123")
		assert.Equal(t, "user-123", UserIDFromContext(ctx))
	})

	t.Run("UserIDFromContext returns empty for no user", func(t *testing.T) {
		assert.Equal(t, "", UserIDFromContext(context.Background()))
	})

	t.Run("IsAdmin from request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), IsAdminKey, true)
		req = req.WithContext(ctx)

		assert.True(t, IsAdmin(req))
	})

	t.Run("IsAdmin returns false for no admin flag", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		assert.False(t, IsAdmin(req))
	})

	t.Run("IsAdminFromContext", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), IsAdminKey, true)
		assert.True(t, IsAdminFromContext(ctx))
	})

	t.Run("IsAdminFromContext returns false for no admin flag", func(t *testing.T) {
		assert.False(t, IsAdminFromContext(context.Background()))
	})

	t.Run("UserFromContext", func(t *testing.T) {
		user := &AuthUser{ID: "user-123", Email: "test@example.com", Name: "Test", IsAdmin: true}
		ctx := context.WithValue(context.Background(), UserKey, user)

		retrieved := UserFromContext(ctx)
		require.NotNil(t, retrieved)
		assert.Equal(t, "user-123", retrieved.ID)
		assert.Equal(t, "test@example.com", retrieved.Email)
		assert.Equal(t, "Test", retrieved.Name)
		assert.True(t, retrieved.IsAdmin)
	})

	t.Run("UserFromContext returns nil for no user", func(t *testing.T) {
		assert.Nil(t, UserFromContext(context.Background()))
	})
}

func TestChainMiddleware(t *testing.T) {
	var order []string

	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m1-before")
			next.ServeHTTP(w, r)
			order = append(order, "m1-after")
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m2-before")
			next.ServeHTTP(w, r)
			order = append(order, "m2-after")
		})
	}

	handler := ChainMiddleware(middleware1, middleware2)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}, order)
}

func TestIsTestMode(t *testing.T) {
	t.Run("returns true when POWERPRO_TEST_MODE is true", func(t *testing.T) {
		os.Setenv("POWERPRO_TEST_MODE", "true")
		defer os.Unsetenv("POWERPRO_TEST_MODE")

		assert.True(t, isTestMode())
	})

	t.Run("returns false when POWERPRO_TEST_MODE is not set", func(t *testing.T) {
		os.Unsetenv("POWERPRO_TEST_MODE")
		assert.False(t, isTestMode())
	})

	t.Run("returns false when POWERPRO_TEST_MODE is false", func(t *testing.T) {
		os.Setenv("POWERPRO_TEST_MODE", "false")
		defer os.Unsetenv("POWERPRO_TEST_MODE")

		assert.False(t, isTestMode())
	})

	t.Run("returns false when POWERPRO_TEST_MODE is other value", func(t *testing.T) {
		os.Setenv("POWERPRO_TEST_MODE", "yes")
		defer os.Unsetenv("POWERPRO_TEST_MODE")

		assert.False(t, isTestMode())
	})
}
