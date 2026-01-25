package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/waynenilsen/power-pro-v3/internal/database"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
)

func setupTestDB(t *testing.T) (*SQLiteUserRepository, *SQLiteSessionRepository, func()) {
	db, cleanup, err := database.OpenTemp("../../migrations")
	require.NoError(t, err)

	userRepo := NewSQLiteUserRepository(db)
	sessionRepo := NewSQLiteSessionRepository(db)

	return userRepo, sessionRepo, cleanup
}

func TestSQLiteUserRepository_Create(t *testing.T) {
	userRepo, _, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	user := &User{
		ID:           "test-user-create-001",
		Email:        "create@example.com",
		Name:         "Create Test",
		PasswordHash: "hashed-password",
		IsAdmin:      false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Verify user was created
	retrieved, err := userRepo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, retrieved.ID)
	assert.Equal(t, user.Email, retrieved.Email)
	assert.Equal(t, user.Name, retrieved.Name)
	assert.Equal(t, user.PasswordHash, retrieved.PasswordHash)
	assert.Equal(t, user.IsAdmin, retrieved.IsAdmin)
}

func TestSQLiteUserRepository_CreateWithEmptyName(t *testing.T) {
	userRepo, _, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	user := &User{
		ID:           "test-user-empty-name",
		Email:        "emptyname@example.com",
		Name:         "", // empty name
		PasswordHash: "hashed-password",
		IsAdmin:      false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	retrieved, err := userRepo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "", retrieved.Name)
}

func TestSQLiteUserRepository_CreateAdmin(t *testing.T) {
	userRepo, _, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	user := &User{
		ID:           "test-admin-create-001",
		Email:        "admin@example.com",
		Name:         "Admin User",
		PasswordHash: "hashed-password",
		IsAdmin:      true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	retrieved, err := userRepo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.True(t, retrieved.IsAdmin)
}

func TestSQLiteUserRepository_GetByID_NotFound(t *testing.T) {
	userRepo, _, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	_, err := userRepo.GetByID(ctx, "nonexistent-user")
	require.Error(t, err)
	assert.True(t, apperrors.IsNotFound(err))
}

func TestSQLiteUserRepository_GetByEmail(t *testing.T) {
	userRepo, _, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	user := &User{
		ID:           "test-user-email-001",
		Email:        "findme@example.com",
		Name:         "Find Me",
		PasswordHash: "hashed-password",
		IsAdmin:      false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Find by exact email
	retrieved, err := userRepo.GetByEmail(ctx, "findme@example.com")
	require.NoError(t, err)
	assert.Equal(t, user.ID, retrieved.ID)

	// Find by email with different case
	retrieved, err = userRepo.GetByEmail(ctx, "FINDME@EXAMPLE.COM")
	require.NoError(t, err)
	assert.Equal(t, user.ID, retrieved.ID)
}

func TestSQLiteUserRepository_GetByEmail_NotFound(t *testing.T) {
	userRepo, _, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	_, err := userRepo.GetByEmail(ctx, "nonexistent@example.com")
	require.Error(t, err)
	assert.True(t, apperrors.IsNotFound(err))
}

func TestSQLiteUserRepository_EmailExists(t *testing.T) {
	userRepo, _, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	// Initially email should not exist
	exists, err := userRepo.EmailExists(ctx, "exists@example.com")
	require.NoError(t, err)
	assert.False(t, exists)

	// Create user
	user := &User{
		ID:           "test-user-exists-001",
		Email:        "exists@example.com",
		Name:         "Exists Test",
		PasswordHash: "hashed-password",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Now email should exist
	exists, err = userRepo.EmailExists(ctx, "exists@example.com")
	require.NoError(t, err)
	assert.True(t, exists)

	// Case insensitive check
	exists, err = userRepo.EmailExists(ctx, "EXISTS@EXAMPLE.COM")
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestSQLiteSessionRepository_Create(t *testing.T) {
	_, sessionRepo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	// Need to use an existing user ID from the test database
	session := &Session{
		ID:        "test-session-001",
		UserID:    "test-user-001", // This user exists from migration seed data
		Token:     "test-token-12345",
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	err := sessionRepo.Create(ctx, session)
	require.NoError(t, err)

	// Verify session was created
	retrieved, err := sessionRepo.GetByToken(ctx, session.Token)
	require.NoError(t, err)
	assert.Equal(t, session.ID, retrieved.ID)
	assert.Equal(t, session.UserID, retrieved.UserID)
	assert.Equal(t, session.Token, retrieved.Token)
}

func TestSQLiteSessionRepository_GetByToken(t *testing.T) {
	_, sessionRepo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	session := &Session{
		ID:        "test-session-get-001",
		UserID:    "test-user-001",
		Token:     "get-token-12345",
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	err := sessionRepo.Create(ctx, session)
	require.NoError(t, err)

	retrieved, err := sessionRepo.GetByToken(ctx, "get-token-12345")
	require.NoError(t, err)
	assert.Equal(t, session.ID, retrieved.ID)
	assert.Equal(t, session.UserID, retrieved.UserID)
}

func TestSQLiteSessionRepository_GetByToken_NotFound(t *testing.T) {
	_, sessionRepo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	_, err := sessionRepo.GetByToken(ctx, "nonexistent-token")
	require.Error(t, err)
	assert.True(t, apperrors.IsNotFound(err))
}

func TestSQLiteSessionRepository_DeleteByToken(t *testing.T) {
	_, sessionRepo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	session := &Session{
		ID:        "test-session-delete-001",
		UserID:    "test-user-001",
		Token:     "delete-token-12345",
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	err := sessionRepo.Create(ctx, session)
	require.NoError(t, err)

	// Verify session exists
	_, err = sessionRepo.GetByToken(ctx, session.Token)
	require.NoError(t, err)

	// Delete session
	err = sessionRepo.DeleteByToken(ctx, session.Token)
	require.NoError(t, err)

	// Verify session is gone
	_, err = sessionRepo.GetByToken(ctx, session.Token)
	require.Error(t, err)
	assert.True(t, apperrors.IsNotFound(err))
}

func TestSQLiteSessionRepository_DeleteByToken_NonExistent(t *testing.T) {
	_, sessionRepo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Deleting non-existent token should not error
	err := sessionRepo.DeleteByToken(ctx, "nonexistent-token")
	require.NoError(t, err)
}

func TestIntegration_RegisterAndLogin(t *testing.T) {
	userRepo, sessionRepo, cleanup := setupTestDB(t)
	defer cleanup()

	svc := NewService(userRepo, sessionRepo)
	ctx := context.Background()

	// Register a new user
	registerResult, err := svc.Register(ctx, RegisterRequest{
		Email:    "integration@example.com",
		Password: "securepassword123",
		Name:     "Integration Test",
	})
	require.NoError(t, err)
	require.NotNil(t, registerResult)
	assert.NotEmpty(t, registerResult.User.ID)

	// Login with the registered user
	loginResult, err := svc.Login(ctx, LoginRequest{
		Email:    "integration@example.com",
		Password: "securepassword123",
	})
	require.NoError(t, err)
	require.NotNil(t, loginResult)
	assert.Equal(t, registerResult.User.ID, loginResult.User.ID)
	assert.NotEmpty(t, loginResult.Token)

	// Validate the session
	user, err := svc.ValidateSession(ctx, loginResult.Token)
	require.NoError(t, err)
	assert.Equal(t, registerResult.User.ID, user.ID)

	// Logout
	err = svc.Logout(ctx, loginResult.Token)
	require.NoError(t, err)

	// Session should be invalid after logout
	_, err = svc.ValidateSession(ctx, loginResult.Token)
	require.Error(t, err)
	assert.True(t, apperrors.IsUnauthorized(err))
}

func TestIntegration_DuplicateEmail(t *testing.T) {
	userRepo, sessionRepo, cleanup := setupTestDB(t)
	defer cleanup()

	svc := NewService(userRepo, sessionRepo)
	ctx := context.Background()

	// Register first user
	_, err := svc.Register(ctx, RegisterRequest{
		Email:    "duplicate@example.com",
		Password: "password123",
		Name:     "First User",
	})
	require.NoError(t, err)

	// Try to register with same email
	_, err = svc.Register(ctx, RegisterRequest{
		Email:    "duplicate@example.com",
		Password: "differentpassword",
		Name:     "Second User",
	})
	require.Error(t, err)
	assert.True(t, apperrors.IsConflict(err))
}

func TestIntegration_WrongPassword(t *testing.T) {
	userRepo, sessionRepo, cleanup := setupTestDB(t)
	defer cleanup()

	svc := NewService(userRepo, sessionRepo)
	ctx := context.Background()

	// Register
	_, err := svc.Register(ctx, RegisterRequest{
		Email:    "wrongpass@example.com",
		Password: "correctpassword",
		Name:     "Test User",
	})
	require.NoError(t, err)

	// Try to login with wrong password
	_, err = svc.Login(ctx, LoginRequest{
		Email:    "wrongpass@example.com",
		Password: "wrongpassword",
	})
	require.Error(t, err)
	assert.True(t, apperrors.IsUnauthorized(err))
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestHelperFunctions(t *testing.T) {
	t.Run("toNullString with empty string", func(t *testing.T) {
		ns := toNullString("")
		assert.False(t, ns.Valid)
	})

	t.Run("toNullString with value", func(t *testing.T) {
		ns := toNullString("hello")
		assert.True(t, ns.Valid)
		assert.Equal(t, "hello", ns.String)
	})

	t.Run("boolToInt true", func(t *testing.T) {
		assert.Equal(t, 1, boolToInt(true))
	})

	t.Run("boolToInt false", func(t *testing.T) {
		assert.Equal(t, 0, boolToInt(false))
	})
}
