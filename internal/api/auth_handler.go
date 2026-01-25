package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/auth"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
)

// AuthHandler handles HTTP requests for authentication operations.
type AuthHandler struct {
	service *auth.Service
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

// RegisterRequest represents the request body for user registration.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

// LoginRequest represents the request body for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse represents the API response format for a user.
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// LoginResponse represents the API response for a successful login.
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expiresAt"`
	User      UserResponse `json:"user"`
}

// sessionDuration must match auth.sessionDuration (7 days)
const sessionDuration = 7 * 24 * time.Hour

func userToResponse(u *auth.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Basic validation for missing fields
	if strings.TrimSpace(req.Email) == "" {
		writeDomainError(w, apperrors.NewValidation("email", "email is required"))
		return
	}
	if req.Password == "" {
		writeDomainError(w, apperrors.NewValidation("password", "password is required"))
		return
	}

	result, err := h.service.Register(r.Context(), auth.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeData(w, http.StatusCreated, userToResponse(result.User))
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Basic validation for missing fields
	if strings.TrimSpace(req.Email) == "" {
		writeDomainError(w, apperrors.NewValidation("email", "email is required"))
		return
	}
	if req.Password == "" {
		writeDomainError(w, apperrors.NewValidation("password", "password is required"))
		return
	}

	result, err := h.service.Login(r.Context(), auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(sessionDuration)

	writeData(w, http.StatusOK, LoginResponse{
		Token:     result.Token,
		ExpiresAt: expiresAt,
		User:      userToResponse(result.User),
	})
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header
	token := extractBearerToken(r)
	if token == "" {
		writeDomainError(w, apperrors.NewUnauthorized("authentication required"))
		return
	}

	// Logout is idempotent - always succeeds
	_ = h.service.Logout(r.Context(), token)

	w.WriteHeader(http.StatusNoContent)
}

// Me handles GET /auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// User should be set in context by auth middleware
	user := middleware.UserFromContext(r.Context())
	if user == nil {
		writeDomainError(w, apperrors.NewUnauthorized("authentication required"))
		return
	}

	// Get full user details from token
	token := extractBearerToken(r)
	if token == "" {
		writeDomainError(w, apperrors.NewUnauthorized("authentication required"))
		return
	}

	fullUser, err := h.service.GetUserBySession(r.Context(), token)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeData(w, http.StatusOK, userToResponse(fullUser))
}

// extractBearerToken extracts the token from the Authorization header.
func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}
