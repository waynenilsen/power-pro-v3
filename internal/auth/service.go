// Package auth provides authentication service layer implementations.
// This package handles user registration, login, logout, and session management.
package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"strings"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	// bcryptCost is the cost factor for bcrypt hashing.
	bcryptCost = 12
	// tokenBytes is the number of random bytes for session tokens.
	tokenBytes = 32
	// sessionDuration is the lifetime of a session.
	sessionDuration = 7 * 24 * time.Hour
	// minPasswordLength is the minimum required password length.
	minPasswordLength = 8
)

// User represents a user in the system.
type User struct {
	ID           string
	Email        string
	Name         string
	PasswordHash string
	IsAdmin      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Session represents an active user session.
type Session struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	EmailExists(ctx context.Context, email string) (bool, error)
}

// SessionRepository defines the interface for session persistence.
type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	GetByToken(ctx context.Context, token string) (*Session, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID string) (int64, error)
	CleanupExpired(ctx context.Context, before time.Time) (int64, error)
}

// Service provides authentication operations.
type Service struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
	now         func() time.Time
}

// NewService creates a new authentication service.
func NewService(userRepo UserRepository, sessionRepo SessionRepository) *Service {
	return &Service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		now:         time.Now,
	}
}

// RegisterRequest contains the data needed to register a new user.
type RegisterRequest struct {
	Email    string
	Password string
	Name     string
}

// RegisterResult contains the result of a successful registration.
type RegisterResult struct {
	User *User
}

// Register creates a new user account.
// It validates the email and password, hashes the password, and creates the user.
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*RegisterResult, error) {
	// Normalize email to lowercase
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Validate email
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	// Validate password
	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := s.userRepo.EmailExists(ctx, email)
	if err != nil {
		return nil, apperrors.NewInternal("failed to check email availability", err)
	}
	if exists {
		return nil, apperrors.NewConflict("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, apperrors.NewInternal("failed to hash password", err)
	}

	// Create user
	now := s.now()
	user := &User{
		ID:           uuid.New().String(),
		Email:        email,
		Name:         strings.TrimSpace(req.Name),
		PasswordHash: string(hashedPassword),
		IsAdmin:      false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperrors.NewInternal("failed to create user", err)
	}

	// Return user without password hash
	return &RegisterResult{
		User: &User{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			IsAdmin:   user.IsAdmin,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

// LoginRequest contains the credentials for login.
type LoginRequest struct {
	Email    string
	Password string
}

// LoginResult contains the result of a successful login.
type LoginResult struct {
	User  *User
	Token string
}

// Login authenticates a user and creates a session.
// Returns the same error for wrong email or password to prevent user enumeration.
func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResult, error) {
	// Normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Lookup user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if apperrors.IsNotFound(err) {
			return nil, apperrors.NewUnauthorized("invalid credentials")
		}
		return nil, apperrors.NewInternal("failed to lookup user", err)
	}
	if user == nil {
		return nil, apperrors.NewUnauthorized("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.NewUnauthorized("invalid credentials")
	}

	// Generate session token
	token, err := generateToken()
	if err != nil {
		return nil, apperrors.NewInternal("failed to generate session token", err)
	}

	// Create session
	now := s.now()
	session := &Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: now.Add(sessionDuration),
		CreatedAt: now,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, apperrors.NewInternal("failed to create session", err)
	}

	// Return user (without password hash) and token
	return &LoginResult{
		User: &User{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			IsAdmin:   user.IsAdmin,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: token,
	}, nil
}

// Logout invalidates a session by its token.
// This operation is idempotent - returns success even if session doesn't exist.
func (s *Service) Logout(ctx context.Context, token string) error {
	if err := s.sessionRepo.DeleteByToken(ctx, token); err != nil {
		// Log the error but don't expose it - logout is idempotent
		return nil
	}
	return nil
}

// ValidateSession checks if a session token is valid and not expired.
// Returns the user if the session is valid.
func (s *Service) ValidateSession(ctx context.Context, token string) (*User, error) {
	if token == "" {
		return nil, apperrors.NewUnauthorized("session token required")
	}

	session, err := s.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		if apperrors.IsNotFound(err) {
			return nil, apperrors.NewUnauthorized("invalid session")
		}
		return nil, apperrors.NewInternal("failed to lookup session", err)
	}
	if session == nil {
		return nil, apperrors.NewUnauthorized("invalid session")
	}

	// Check if session is expired
	if s.now().After(session.ExpiresAt) {
		return nil, apperrors.NewUnauthorized("session expired")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		if apperrors.IsNotFound(err) {
			return nil, apperrors.NewUnauthorized("user not found")
		}
		return nil, apperrors.NewInternal("failed to lookup user", err)
	}
	if user == nil {
		return nil, apperrors.NewUnauthorized("user not found")
	}

	// Return user without password hash
	return &User{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUserBySession retrieves the full user object from a session token.
// This is a convenience method that wraps ValidateSession.
func (s *Service) GetUserBySession(ctx context.Context, token string) (*User, error) {
	return s.ValidateSession(ctx, token)
}

// CleanupExpiredSessions removes all expired sessions from the database.
// Returns the number of sessions deleted.
func (s *Service) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	return s.sessionRepo.CleanupExpired(ctx, s.now())
}

// DeleteUserSessions removes all sessions for a specific user.
// Returns the number of sessions deleted.
func (s *Service) DeleteUserSessions(ctx context.Context, userID string) (int64, error) {
	return s.sessionRepo.DeleteByUserID(ctx, userID)
}

// validateEmail validates an email address.
func validateEmail(email string) error {
	if email == "" {
		return apperrors.NewValidation("email", "email is required")
	}
	if !strings.Contains(email, "@") {
		return apperrors.NewValidation("email", "invalid email format")
	}
	// Check for domain part after @
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return apperrors.NewValidation("email", "invalid email format")
	}
	// Check domain has at least one dot (basic check)
	if !strings.Contains(parts[1], ".") {
		return apperrors.NewValidation("email", "invalid email format")
	}
	return nil
}

// validatePassword validates a password.
func validatePassword(password string) error {
	if len(password) < minPasswordLength {
		return apperrors.NewValidation("password", "password must be at least 8 characters")
	}
	return nil
}

// generateToken generates a cryptographically secure random token.
func generateToken() (string, error) {
	bytes := make([]byte, tokenBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// SQLiteUserRepository implements UserRepository using SQLite.
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository creates a new SQLite-backed user repository.
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// GetByID retrieves a user by ID.
func (r *SQLiteUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User
	var email, name, passwordHash sql.NullString
	var isAdmin int64
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, name, password_hash, is_admin, created_at, updated_at
		FROM users WHERE id = ?
	`, id).Scan(&user.ID, &email, &name, &passwordHash, &isAdmin, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, apperrors.NewNotFound("user", id)
	}
	if err != nil {
		return nil, err
	}

	user.Email = email.String
	user.Name = name.String
	user.PasswordHash = passwordHash.String
	user.IsAdmin = isAdmin == 1
	user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	user.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &user, nil
}

// GetByEmail retrieves a user by email (case-insensitive).
func (r *SQLiteUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	var dbEmail, name, passwordHash sql.NullString
	var isAdmin int64
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, name, password_hash, is_admin, created_at, updated_at
		FROM users WHERE LOWER(email) = LOWER(?)
	`, email).Scan(&user.ID, &dbEmail, &name, &passwordHash, &isAdmin, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, apperrors.NewNotFound("user", email)
	}
	if err != nil {
		return nil, err
	}

	user.Email = dbEmail.String
	user.Name = name.String
	user.PasswordHash = passwordHash.String
	user.IsAdmin = isAdmin == 1
	user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	user.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &user, nil
}

// Create persists a new user.
func (r *SQLiteUserRepository) Create(ctx context.Context, user *User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, email, name, password_hash, is_admin, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, user.ID, user.Email, toNullString(user.Name), user.PasswordHash, boolToInt(user.IsAdmin), user.CreatedAt, user.UpdatedAt)
	return err
}

// EmailExists checks if an email is already registered.
func (r *SQLiteUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users WHERE LOWER(email) = LOWER(?)
	`, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SQLiteSessionRepository implements SessionRepository using SQLite.
type SQLiteSessionRepository struct {
	db *sql.DB
}

// NewSQLiteSessionRepository creates a new SQLite-backed session repository.
func NewSQLiteSessionRepository(db *sql.DB) *SQLiteSessionRepository {
	return &SQLiteSessionRepository{db: db}
}

// Create persists a new session.
func (r *SQLiteSessionRepository) Create(ctx context.Context, session *Session) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO sessions (id, user_id, token, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, session.ID, session.UserID, session.Token, session.ExpiresAt.Format(time.RFC3339), session.CreatedAt.Format(time.RFC3339))
	return err
}

// GetByToken retrieves a session by its token.
func (r *SQLiteSessionRepository) GetByToken(ctx context.Context, token string) (*Session, error) {
	var session Session
	var expiresAt, createdAt string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token, expires_at, created_at
		FROM sessions WHERE token = ?
	`, token).Scan(&session.ID, &session.UserID, &session.Token, &expiresAt, &createdAt)

	if err == sql.ErrNoRows {
		return nil, apperrors.NewNotFound("session", token)
	}
	if err != nil {
		return nil, err
	}

	session.ExpiresAt, _ = time.Parse(time.RFC3339, expiresAt)
	session.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)

	return &session, nil
}

// DeleteByToken deletes a session by its token.
func (r *SQLiteSessionRepository) DeleteByToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE token = ?`, token)
	return err
}

// DeleteByUserID deletes all sessions for a user.
// Returns the number of sessions deleted.
func (r *SQLiteSessionRepository) DeleteByUserID(ctx context.Context, userID string) (int64, error) {
	result, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE user_id = ?`, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// CleanupExpired deletes all sessions that expired before the given time.
// Returns the number of sessions deleted.
func (r *SQLiteSessionRepository) CleanupExpired(ctx context.Context, before time.Time) (int64, error) {
	result, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE expires_at < ?`, before.Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Helper functions

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
