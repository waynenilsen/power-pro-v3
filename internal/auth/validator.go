package auth

import (
	"context"

	"github.com/waynenilsen/power-pro-v3/internal/middleware"
)

// SessionValidatorAdapter adapts the auth.Service to implement middleware.SessionValidator.
// This allows the auth service to be used with the auth middleware.
type SessionValidatorAdapter struct {
	service *Service
}

// NewSessionValidatorAdapter creates a new adapter that wraps the auth service.
func NewSessionValidatorAdapter(service *Service) *SessionValidatorAdapter {
	return &SessionValidatorAdapter{service: service}
}

// ValidateSession implements middleware.SessionValidator.
// It validates the token and returns user information for the middleware context.
func (a *SessionValidatorAdapter) ValidateSession(ctx context.Context, token string) (*middleware.AuthUser, error) {
	user, err := a.service.ValidateSession(ctx, token)
	if err != nil {
		return nil, err
	}

	return &middleware.AuthUser{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		IsAdmin: user.IsAdmin,
	}, nil
}
