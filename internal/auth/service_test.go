package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepo is a mock implementation of UserRepository for testing.
type mockUserRepo struct {
	users       map[string]*User
	emailIndex  map[string]string // email -> user ID
	createErr   error
	getByIDErr  error
	getByEmailErr error
	emailExistsErr error
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:      make(map[string]*User),
		emailIndex: make(map[string]string),
	}
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	user, ok := m.users[id]
	if !ok {
		return nil, apperrors.NewNotFound("user", id)
	}
	return user, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	if m.getByEmailErr != nil {
		return nil, m.getByEmailErr
	}
	id, ok := m.emailIndex[email]
	if !ok {
		return nil, apperrors.NewNotFound("user", email)
	}
	return m.users[id], nil
}

func (m *mockUserRepo) Create(ctx context.Context, user *User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.users[user.ID] = user
	m.emailIndex[user.Email] = user.ID
	return nil
}

func (m *mockUserRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	if m.emailExistsErr != nil {
		return false, m.emailExistsErr
	}
	_, exists := m.emailIndex[email]
	return exists, nil
}

// mockSessionRepo is a mock implementation of SessionRepository for testing.
type mockSessionRepo struct {
	sessions    map[string]*Session // token -> session
	createErr   error
	getByTokenErr error
	deleteErr   error
}

func newMockSessionRepo() *mockSessionRepo {
	return &mockSessionRepo{
		sessions: make(map[string]*Session),
	}
}

func (m *mockSessionRepo) Create(ctx context.Context, session *Session) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.sessions[session.Token] = session
	return nil
}

func (m *mockSessionRepo) GetByToken(ctx context.Context, token string) (*Session, error) {
	if m.getByTokenErr != nil {
		return nil, m.getByTokenErr
	}
	session, ok := m.sessions[token]
	if !ok {
		return nil, apperrors.NewNotFound("session", token)
	}
	return session, nil
}

func (m *mockSessionRepo) DeleteByToken(ctx context.Context, token string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.sessions, token)
	return nil
}

func TestService_Register(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		request     RegisterRequest
		setupMock   func(*mockUserRepo, *mockSessionRepo)
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *RegisterResult)
	}{
		{
			name: "successful registration",
			request: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:   false,
			checkResult: func(t *testing.T, result *RegisterResult) {
				assert.NotEmpty(t, result.User.ID)
				assert.Equal(t, "test@example.com", result.User.Email)
				assert.Equal(t, "Test User", result.User.Name)
				assert.False(t, result.User.IsAdmin)
				assert.Empty(t, result.User.PasswordHash) // Should not be exposed
			},
		},
		{
			name: "email normalized to lowercase",
			request: RegisterRequest{
				Email:    "TEST@EXAMPLE.COM",
				Password: "password123",
				Name:     "Test User",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:   false,
			checkResult: func(t *testing.T, result *RegisterResult) {
				assert.Equal(t, "test@example.com", result.User.Email)
			},
		},
		{
			name: "email with whitespace trimmed",
			request: RegisterRequest{
				Email:    "  test@example.com  ",
				Password: "password123",
				Name:     "Test User",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:   false,
			checkResult: func(t *testing.T, result *RegisterResult) {
				assert.Equal(t, "test@example.com", result.User.Email)
			},
		},
		{
			name: "name with whitespace trimmed",
			request: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "  Test User  ",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:   false,
			checkResult: func(t *testing.T, result *RegisterResult) {
				assert.Equal(t, "Test User", result.User.Name)
			},
		},
		{
			name: "empty email",
			request: RegisterRequest{
				Email:    "",
				Password: "password123",
				Name:     "Test User",
			},
			setupMock:   func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:     true,
			errContains: "email is required",
		},
		{
			name: "email without @",
			request: RegisterRequest{
				Email:    "testexample.com",
				Password: "password123",
				Name:     "Test User",
			},
			setupMock:   func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:     true,
			errContains: "invalid email format",
		},
		{
			name: "email with empty local part",
			request: RegisterRequest{
				Email:    "@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			setupMock:   func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:     true,
			errContains: "invalid email format",
		},
		{
			name: "email with empty domain",
			request: RegisterRequest{
				Email:    "test@",
				Password: "password123",
				Name:     "Test User",
			},
			setupMock:   func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:     true,
			errContains: "invalid email format",
		},
		{
			name: "email without domain dot",
			request: RegisterRequest{
				Email:    "test@examplecom",
				Password: "password123",
				Name:     "Test User",
			},
			setupMock:   func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:     true,
			errContains: "invalid email format",
		},
		{
			name: "password too short",
			request: RegisterRequest{
				Email:    "test@example.com",
				Password: "1234567",
				Name:     "Test User",
			},
			setupMock:   func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:     true,
			errContains: "password must be at least 8 characters",
		},
		{
			name: "password exactly 8 characters",
			request: RegisterRequest{
				Email:    "test@example.com",
				Password: "12345678",
				Name:     "Test User",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:   false,
		},
		{
			name: "email already exists",
			request: RegisterRequest{
				Email:    "existing@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				ur.users["existing-id"] = &User{ID: "existing-id", Email: "existing@example.com"}
				ur.emailIndex["existing@example.com"] = "existing-id"
			},
			wantErr:     true,
			errContains: "email already registered",
		},
		{
			name: "email check fails",
			request: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				ur.emailExistsErr = errors.New("database error")
			},
			wantErr:     true,
			errContains: "failed to check email availability",
		},
		{
			name: "create user fails",
			request: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				ur.createErr = errors.New("database error")
			},
			wantErr:     true,
			errContains: "failed to create user",
		},
		{
			name: "name is optional",
			request: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:   false,
			checkResult: func(t *testing.T, result *RegisterResult) {
				assert.Equal(t, "", result.User.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := newMockUserRepo()
			sessionRepo := newMockSessionRepo()
			tt.setupMock(userRepo, sessionRepo)

			svc := NewService(userRepo, sessionRepo)
			svc.now = func() time.Time { return fixedTime }

			result, err := svc.Register(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}

			// Verify password was hashed
			user := userRepo.users[result.User.ID]
			require.NotNil(t, user)
			require.NotEmpty(t, user.PasswordHash)
			err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(tt.request.Password))
			require.NoError(t, err, "password should be properly hashed")
		})
	}
}

func TestService_Login(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcryptCost)

	tests := []struct {
		name        string
		request     LoginRequest
		setupMock   func(*mockUserRepo, *mockSessionRepo)
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *LoginResult, *mockSessionRepo)
	}{
		{
			name: "successful login",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				user := &User{
					ID:           "user-123",
					Email:        "test@example.com",
					Name:         "Test User",
					PasswordHash: string(hashedPassword),
					CreatedAt:    fixedTime.Add(-24 * time.Hour),
					UpdatedAt:    fixedTime.Add(-24 * time.Hour),
				}
				ur.users[user.ID] = user
				ur.emailIndex[user.Email] = user.ID
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *LoginResult, sr *mockSessionRepo) {
				assert.Equal(t, "user-123", result.User.ID)
				assert.Equal(t, "test@example.com", result.User.Email)
				assert.Equal(t, "Test User", result.User.Name)
				assert.NotEmpty(t, result.Token)
				assert.Empty(t, result.User.PasswordHash) // Should not be exposed

				// Verify session was created
				session := sr.sessions[result.Token]
				require.NotNil(t, session)
				assert.Equal(t, "user-123", session.UserID)
				assert.Equal(t, fixedTime.Add(7*24*time.Hour), session.ExpiresAt)
			},
		},
		{
			name: "email case insensitive",
			request: LoginRequest{
				Email:    "TEST@EXAMPLE.COM",
				Password: "password123",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				user := &User{
					ID:           "user-123",
					Email:        "test@example.com",
					Name:         "Test User",
					PasswordHash: string(hashedPassword),
				}
				ur.users[user.ID] = user
				ur.emailIndex["test@example.com"] = user.ID
			},
			wantErr: false,
		},
		{
			name: "wrong email",
			request: LoginRequest{
				Email:    "wrong@example.com",
				Password: "password123",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				user := &User{
					ID:           "user-123",
					Email:        "test@example.com",
					PasswordHash: string(hashedPassword),
				}
				ur.users[user.ID] = user
				ur.emailIndex[user.Email] = user.ID
			},
			wantErr:     true,
			errContains: "invalid credentials",
		},
		{
			name: "wrong password",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				user := &User{
					ID:           "user-123",
					Email:        "test@example.com",
					PasswordHash: string(hashedPassword),
				}
				ur.users[user.ID] = user
				ur.emailIndex[user.Email] = user.ID
			},
			wantErr:     true,
			errContains: "invalid credentials",
		},
		{
			name: "same error for wrong email and wrong password",
			request: LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock:   func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:     true,
			errContains: "invalid credentials",
		},
		{
			name: "session creation fails",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				user := &User{
					ID:           "user-123",
					Email:        "test@example.com",
					PasswordHash: string(hashedPassword),
				}
				ur.users[user.ID] = user
				ur.emailIndex[user.Email] = user.ID
				sr.createErr = errors.New("database error")
			},
			wantErr:     true,
			errContains: "failed to create session",
		},
		{
			name: "user lookup database error",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				ur.getByEmailErr = errors.New("database error")
			},
			wantErr:     true,
			errContains: "failed to lookup user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := newMockUserRepo()
			sessionRepo := newMockSessionRepo()
			tt.setupMock(userRepo, sessionRepo)

			svc := NewService(userRepo, sessionRepo)
			svc.now = func() time.Time { return fixedTime }

			result, err := svc.Login(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			if tt.checkResult != nil {
				tt.checkResult(t, result, sessionRepo)
			}
		})
	}
}

func TestService_Logout(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		setupMock func(*mockUserRepo, *mockSessionRepo)
		wantErr   bool
	}{
		{
			name:  "successful logout",
			token: "valid-token",
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				sr.sessions["valid-token"] = &Session{
					ID:     "session-123",
					Token:  "valid-token",
					UserID: "user-123",
				}
			},
			wantErr: false,
		},
		{
			name:      "logout with nonexistent token - idempotent",
			token:     "nonexistent-token",
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:   false,
		},
		{
			name:  "logout with delete error - still returns success",
			token: "error-token",
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				sr.deleteErr = errors.New("database error")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := newMockUserRepo()
			sessionRepo := newMockSessionRepo()
			tt.setupMock(userRepo, sessionRepo)

			svc := NewService(userRepo, sessionRepo)

			err := svc.Logout(context.Background(), tt.token)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestService_ValidateSession(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		token       string
		setupMock   func(*mockUserRepo, *mockSessionRepo)
		wantErr     bool
		errContains string
		checkUser   func(*testing.T, *User)
	}{
		{
			name:  "valid session",
			token: "valid-token",
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				user := &User{
					ID:        "user-123",
					Email:     "test@example.com",
					Name:      "Test User",
					IsAdmin:   false,
					CreatedAt: fixedTime.Add(-24 * time.Hour),
					UpdatedAt: fixedTime.Add(-24 * time.Hour),
				}
				ur.users[user.ID] = user

				sr.sessions["valid-token"] = &Session{
					ID:        "session-123",
					UserID:    "user-123",
					Token:     "valid-token",
					ExpiresAt: fixedTime.Add(24 * time.Hour), // expires in future
					CreatedAt: fixedTime.Add(-1 * time.Hour),
				}
			},
			wantErr: false,
			checkUser: func(t *testing.T, user *User) {
				assert.Equal(t, "user-123", user.ID)
				assert.Equal(t, "test@example.com", user.Email)
				assert.Equal(t, "Test User", user.Name)
				assert.Empty(t, user.PasswordHash) // Should not be exposed
			},
		},
		{
			name:        "empty token",
			token:       "",
			setupMock:   func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:     true,
			errContains: "session token required",
		},
		{
			name:        "nonexistent session",
			token:       "nonexistent-token",
			setupMock:   func(ur *mockUserRepo, sr *mockSessionRepo) {},
			wantErr:     true,
			errContains: "invalid session",
		},
		{
			name:  "expired session",
			token: "expired-token",
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				sr.sessions["expired-token"] = &Session{
					ID:        "session-123",
					UserID:    "user-123",
					Token:     "expired-token",
					ExpiresAt: fixedTime.Add(-1 * time.Hour), // expired
					CreatedAt: fixedTime.Add(-8 * 24 * time.Hour),
				}
			},
			wantErr:     true,
			errContains: "session expired",
		},
		{
			name:  "user not found for session",
			token: "orphan-token",
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				sr.sessions["orphan-token"] = &Session{
					ID:        "session-123",
					UserID:    "nonexistent-user",
					Token:     "orphan-token",
					ExpiresAt: fixedTime.Add(24 * time.Hour),
					CreatedAt: fixedTime.Add(-1 * time.Hour),
				}
			},
			wantErr:     true,
			errContains: "user not found",
		},
		{
			name:  "session lookup error",
			token: "error-token",
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				sr.getByTokenErr = errors.New("database error")
			},
			wantErr:     true,
			errContains: "failed to lookup session",
		},
		{
			name:  "user lookup error",
			token: "user-error-token",
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				sr.sessions["user-error-token"] = &Session{
					ID:        "session-123",
					UserID:    "user-123",
					Token:     "user-error-token",
					ExpiresAt: fixedTime.Add(24 * time.Hour),
					CreatedAt: fixedTime.Add(-1 * time.Hour),
				}
				ur.getByIDErr = errors.New("database error")
			},
			wantErr:     true,
			errContains: "failed to lookup user",
		},
		{
			name:  "admin user session",
			token: "admin-token",
			setupMock: func(ur *mockUserRepo, sr *mockSessionRepo) {
				user := &User{
					ID:        "admin-123",
					Email:     "admin@example.com",
					Name:      "Admin User",
					IsAdmin:   true,
					CreatedAt: fixedTime.Add(-24 * time.Hour),
					UpdatedAt: fixedTime.Add(-24 * time.Hour),
				}
				ur.users[user.ID] = user

				sr.sessions["admin-token"] = &Session{
					ID:        "session-456",
					UserID:    "admin-123",
					Token:     "admin-token",
					ExpiresAt: fixedTime.Add(24 * time.Hour),
					CreatedAt: fixedTime.Add(-1 * time.Hour),
				}
			},
			wantErr: false,
			checkUser: func(t *testing.T, user *User) {
				assert.True(t, user.IsAdmin)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := newMockUserRepo()
			sessionRepo := newMockSessionRepo()
			tt.setupMock(userRepo, sessionRepo)

			svc := NewService(userRepo, sessionRepo)
			svc.now = func() time.Time { return fixedTime }

			user, err := svc.ValidateSession(context.Background(), tt.token)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)

			if tt.checkUser != nil {
				tt.checkUser(t, user)
			}
		})
	}
}

func TestService_GetUserBySession(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	userRepo := newMockUserRepo()
	sessionRepo := newMockSessionRepo()

	user := &User{
		ID:        "user-123",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: fixedTime.Add(-24 * time.Hour),
		UpdatedAt: fixedTime.Add(-24 * time.Hour),
	}
	userRepo.users[user.ID] = user

	sessionRepo.sessions["valid-token"] = &Session{
		ID:        "session-123",
		UserID:    "user-123",
		Token:     "valid-token",
		ExpiresAt: fixedTime.Add(24 * time.Hour),
		CreatedAt: fixedTime.Add(-1 * time.Hour),
	}

	svc := NewService(userRepo, sessionRepo)
	svc.now = func() time.Time { return fixedTime }

	// Test that GetUserBySession is an alias for ValidateSession
	result, err := svc.GetUserBySession(context.Background(), "valid-token")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "user-123", result.ID)
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email   string
		wantErr bool
	}{
		{"test@example.com", false},
		{"user.name@domain.org", false},
		{"user+tag@example.co.uk", false},
		{"", true},
		{"testexample.com", true},
		{"@example.com", true},
		{"test@", true},
		{"test@domain", true},
		{"test@@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			err := validateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		wantErr  bool
	}{
		{"12345678", false},         // exactly 8 chars
		{"password123", false},      // more than 8 chars
		{"a very long password", false},
		{"1234567", true},           // 7 chars
		{"short", true},             // too short
		{"", true},                  // empty
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			err := validatePassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	// Generate multiple tokens and ensure they're all unique
	tokens := make(map[string]bool)
	for i := 0; i < 100; i++ {
		token, err := generateToken()
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// Token should be base64 encoded 32 bytes = 44 chars (with padding) or 43 (without)
		assert.GreaterOrEqual(t, len(token), 43)

		// Should be unique
		assert.False(t, tokens[token], "token should be unique")
		tokens[token] = true
	}
}

func TestBcryptCost(t *testing.T) {
	// Verify that bcrypt cost is set correctly
	assert.Equal(t, 12, bcryptCost)
}

func TestSessionDuration(t *testing.T) {
	// Verify session duration is 7 days
	assert.Equal(t, 7*24*time.Hour, sessionDuration)
}

func TestTokenBytes(t *testing.T) {
	// Verify token is 32 bytes
	assert.Equal(t, 32, tokenBytes)
}

func TestMinPasswordLength(t *testing.T) {
	// Verify minimum password length is 8
	assert.Equal(t, 8, minPasswordLength)
}
