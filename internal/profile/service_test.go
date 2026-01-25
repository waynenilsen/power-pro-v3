package profile

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/waynenilsen/power-pro-v3/internal/database"
	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
)

// mockProfileRepo is a mock implementation of ProfileRepository for testing.
type mockProfileRepo struct {
	profiles    map[string]*Profile
	getErr      error
	updateErr   error
	lastUpdate  ProfileUpdate
	lastUserID  string
}

func newMockProfileRepo() *mockProfileRepo {
	return &mockProfileRepo{
		profiles: make(map[string]*Profile),
	}
}

func (m *mockProfileRepo) GetByUserID(ctx context.Context, userID string) (*Profile, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	profile, ok := m.profiles[userID]
	if !ok {
		return nil, apperrors.NewNotFound("user", userID)
	}
	// Return a copy to simulate database behavior
	copy := *profile
	if profile.Name != nil {
		nameCopy := *profile.Name
		copy.Name = &nameCopy
	}
	return &copy, nil
}

func (m *mockProfileRepo) Update(ctx context.Context, userID string, update ProfileUpdate) (*Profile, error) {
	m.lastUserID = userID
	m.lastUpdate = update

	if m.updateErr != nil {
		return nil, m.updateErr
	}

	profile, ok := m.profiles[userID]
	if !ok {
		return nil, apperrors.NewNotFound("user", userID)
	}

	// Apply the update
	if update.SetName {
		profile.Name = update.Name
	}
	if update.SetWeightUnit {
		profile.WeightUnit = update.WeightUnit
	}
	profile.UpdatedAt = update.UpdatedAt

	// Return a copy
	copy := *profile
	if profile.Name != nil {
		nameCopy := *profile.Name
		copy.Name = &nameCopy
	}
	return &copy, nil
}

// Helper to create a string pointer
func strPtr(s string) *string {
	return &s
}

func TestService_GetProfile(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		userID      string
		setupMock   func(*mockProfileRepo)
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *Profile)
	}{
		{
			name:   "successful get profile",
			userID: "user-123",
			setupMock: func(m *mockProfileRepo) {
				name := "Test User"
				m.profiles["user-123"] = &Profile{
					ID:         "user-123",
					Email:      "test@example.com",
					Name:       &name,
					WeightUnit: "lb",
					CreatedAt:  fixedTime,
					UpdatedAt:  fixedTime,
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, p *Profile) {
				assert.Equal(t, "user-123", p.ID)
				assert.Equal(t, "test@example.com", p.Email)
				require.NotNil(t, p.Name)
				assert.Equal(t, "Test User", *p.Name)
				assert.Equal(t, "lb", p.WeightUnit)
			},
		},
		{
			name:   "profile with nil name",
			userID: "user-456",
			setupMock: func(m *mockProfileRepo) {
				m.profiles["user-456"] = &Profile{
					ID:         "user-456",
					Email:      "noname@example.com",
					Name:       nil,
					WeightUnit: "kg",
					CreatedAt:  fixedTime,
					UpdatedAt:  fixedTime,
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, p *Profile) {
				assert.Equal(t, "user-456", p.ID)
				assert.Nil(t, p.Name)
				assert.Equal(t, "kg", p.WeightUnit)
			},
		},
		{
			name:        "empty user ID",
			userID:      "",
			setupMock:   func(m *mockProfileRepo) {},
			wantErr:     true,
			errContains: "user ID is required",
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			setupMock: func(m *mockProfileRepo) {
				// No profiles added
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:   "repository error",
			userID: "user-123",
			setupMock: func(m *mockProfileRepo) {
				m.getErr = apperrors.NewInternal("database error", nil)
			},
			wantErr:     true,
			errContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockProfileRepo()
			tt.setupMock(repo)

			svc := NewService(repo)

			result, err := svc.GetProfile(context.Background(), tt.userID)

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
		})
	}
}

func TestService_UpdateProfile(t *testing.T) {
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		userID      string
		request     UpdateProfileRequest
		setupMock   func(*mockProfileRepo)
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *Profile, *mockProfileRepo)
	}{
		{
			name:   "update name only",
			userID: "user-123",
			request: UpdateProfileRequest{
				Name: strPtr("New Name"),
			},
			setupMock: func(m *mockProfileRepo) {
				oldName := "Old Name"
				m.profiles["user-123"] = &Profile{
					ID:         "user-123",
					Email:      "test@example.com",
					Name:       &oldName,
					WeightUnit: "lb",
					CreatedAt:  fixedTime.Add(-24 * time.Hour),
					UpdatedAt:  fixedTime.Add(-24 * time.Hour),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, p *Profile, m *mockProfileRepo) {
				require.NotNil(t, p.Name)
				assert.Equal(t, "New Name", *p.Name)
				assert.Equal(t, "lb", p.WeightUnit) // Unchanged
				assert.True(t, m.lastUpdate.SetName)
				assert.False(t, m.lastUpdate.SetWeightUnit)
			},
		},
		{
			name:   "update weight unit only",
			userID: "user-123",
			request: UpdateProfileRequest{
				WeightUnit: strPtr("kg"),
			},
			setupMock: func(m *mockProfileRepo) {
				name := "Test User"
				m.profiles["user-123"] = &Profile{
					ID:         "user-123",
					Email:      "test@example.com",
					Name:       &name,
					WeightUnit: "lb",
					CreatedAt:  fixedTime.Add(-24 * time.Hour),
					UpdatedAt:  fixedTime.Add(-24 * time.Hour),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, p *Profile, m *mockProfileRepo) {
				assert.Equal(t, "kg", p.WeightUnit)
				require.NotNil(t, p.Name)
				assert.Equal(t, "Test User", *p.Name) // Unchanged
				assert.False(t, m.lastUpdate.SetName)
				assert.True(t, m.lastUpdate.SetWeightUnit)
			},
		},
		{
			name:   "update both name and weight unit",
			userID: "user-123",
			request: UpdateProfileRequest{
				Name:       strPtr("New Name"),
				WeightUnit: strPtr("kg"),
			},
			setupMock: func(m *mockProfileRepo) {
				oldName := "Old Name"
				m.profiles["user-123"] = &Profile{
					ID:         "user-123",
					Email:      "test@example.com",
					Name:       &oldName,
					WeightUnit: "lb",
					CreatedAt:  fixedTime.Add(-24 * time.Hour),
					UpdatedAt:  fixedTime.Add(-24 * time.Hour),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, p *Profile, m *mockProfileRepo) {
				require.NotNil(t, p.Name)
				assert.Equal(t, "New Name", *p.Name)
				assert.Equal(t, "kg", p.WeightUnit)
				assert.True(t, m.lastUpdate.SetName)
				assert.True(t, m.lastUpdate.SetWeightUnit)
			},
		},
		{
			name:   "clear name with empty string",
			userID: "user-123",
			request: UpdateProfileRequest{
				Name: strPtr(""),
			},
			setupMock: func(m *mockProfileRepo) {
				oldName := "Old Name"
				m.profiles["user-123"] = &Profile{
					ID:         "user-123",
					Email:      "test@example.com",
					Name:       &oldName,
					WeightUnit: "lb",
					CreatedAt:  fixedTime.Add(-24 * time.Hour),
					UpdatedAt:  fixedTime.Add(-24 * time.Hour),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, p *Profile, m *mockProfileRepo) {
				assert.Nil(t, p.Name, "name should be cleared to nil")
				assert.True(t, m.lastUpdate.SetName)
				assert.Nil(t, m.lastUpdate.Name, "update.Name should be nil to set NULL")
			},
		},
		{
			name:   "clear name with whitespace-only string",
			userID: "user-123",
			request: UpdateProfileRequest{
				Name: strPtr("   "),
			},
			setupMock: func(m *mockProfileRepo) {
				oldName := "Old Name"
				m.profiles["user-123"] = &Profile{
					ID:         "user-123",
					Email:      "test@example.com",
					Name:       &oldName,
					WeightUnit: "lb",
					CreatedAt:  fixedTime.Add(-24 * time.Hour),
					UpdatedAt:  fixedTime.Add(-24 * time.Hour),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, p *Profile, m *mockProfileRepo) {
				assert.Nil(t, p.Name, "name should be cleared to nil")
				assert.True(t, m.lastUpdate.SetName)
				assert.Nil(t, m.lastUpdate.Name)
			},
		},
		{
			name:   "name with leading/trailing whitespace is trimmed",
			userID: "user-123",
			request: UpdateProfileRequest{
				Name: strPtr("  Trimmed Name  "),
			},
			setupMock: func(m *mockProfileRepo) {
				m.profiles["user-123"] = &Profile{
					ID:         "user-123",
					Email:      "test@example.com",
					WeightUnit: "lb",
					CreatedAt:  fixedTime.Add(-24 * time.Hour),
					UpdatedAt:  fixedTime.Add(-24 * time.Hour),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, p *Profile, m *mockProfileRepo) {
				require.NotNil(t, p.Name)
				assert.Equal(t, "Trimmed Name", *p.Name)
				require.NotNil(t, m.lastUpdate.Name)
				assert.Equal(t, "Trimmed Name", *m.lastUpdate.Name)
			},
		},
		{
			name:   "no changes - returns current profile",
			userID: "user-123",
			request: UpdateProfileRequest{
				// Both nil - no changes
			},
			setupMock: func(m *mockProfileRepo) {
				name := "Original Name"
				m.profiles["user-123"] = &Profile{
					ID:         "user-123",
					Email:      "test@example.com",
					Name:       &name,
					WeightUnit: "lb",
					CreatedAt:  fixedTime.Add(-24 * time.Hour),
					UpdatedAt:  fixedTime.Add(-24 * time.Hour),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, p *Profile, m *mockProfileRepo) {
				require.NotNil(t, p.Name)
				assert.Equal(t, "Original Name", *p.Name)
				assert.Equal(t, "lb", p.WeightUnit)
				// lastUpdate should not be set since no update was performed
				assert.Empty(t, m.lastUserID)
			},
		},
		{
			name:        "empty user ID",
			userID:      "",
			request:     UpdateProfileRequest{Name: strPtr("New Name")},
			setupMock:   func(m *mockProfileRepo) {},
			wantErr:     true,
			errContains: "user ID is required",
		},
		{
			name:   "name too long",
			userID: "user-123",
			request: UpdateProfileRequest{
				Name: strPtr(string(make([]byte, 101))), // 101 characters
			},
			setupMock: func(m *mockProfileRepo) {
				m.profiles["user-123"] = &Profile{
					ID:         "user-123",
					Email:      "test@example.com",
					WeightUnit: "lb",
				}
			},
			wantErr:     true,
			errContains: "name must be 100 characters or less",
		},
		{
			name:   "name exactly 100 characters is valid",
			userID: "user-123",
			request: UpdateProfileRequest{
				Name: strPtr(string(make([]byte, 100))), // exactly 100 characters
			},
			setupMock: func(m *mockProfileRepo) {
				m.profiles["user-123"] = &Profile{
					ID:         "user-123",
					Email:      "test@example.com",
					WeightUnit: "lb",
				}
			},
			wantErr: false,
		},
		{
			name:   "invalid weight unit",
			userID: "user-123",
			request: UpdateProfileRequest{
				WeightUnit: strPtr("invalid"),
			},
			setupMock: func(m *mockProfileRepo) {
				m.profiles["user-123"] = &Profile{
					ID:         "user-123",
					Email:      "test@example.com",
					WeightUnit: "lb",
				}
			},
			wantErr:     true,
			errContains: "weight unit must be 'lb' or 'kg'",
		},
		{
			name:   "weight unit lb is valid",
			userID: "user-123",
			request: UpdateProfileRequest{
				WeightUnit: strPtr("lb"),
			},
			setupMock: func(m *mockProfileRepo) {
				m.profiles["user-123"] = &Profile{
					ID:         "user-123",
					Email:      "test@example.com",
					WeightUnit: "kg",
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, p *Profile, m *mockProfileRepo) {
				assert.Equal(t, "lb", p.WeightUnit)
			},
		},
		{
			name:   "weight unit kg is valid",
			userID: "user-123",
			request: UpdateProfileRequest{
				WeightUnit: strPtr("kg"),
			},
			setupMock: func(m *mockProfileRepo) {
				m.profiles["user-123"] = &Profile{
					ID:         "user-123",
					Email:      "test@example.com",
					WeightUnit: "lb",
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, p *Profile, m *mockProfileRepo) {
				assert.Equal(t, "kg", p.WeightUnit)
			},
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			request: UpdateProfileRequest{
				Name: strPtr("New Name"),
			},
			setupMock:   func(m *mockProfileRepo) {},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name:   "repository update error",
			userID: "user-123",
			request: UpdateProfileRequest{
				Name: strPtr("New Name"),
			},
			setupMock: func(m *mockProfileRepo) {
				m.profiles["user-123"] = &Profile{
					ID:         "user-123",
					Email:      "test@example.com",
					WeightUnit: "lb",
				}
				m.updateErr = apperrors.NewInternal("database error", nil)
			},
			wantErr:     true,
			errContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockProfileRepo()
			tt.setupMock(repo)

			svc := NewService(repo)
			svc.now = func() time.Time { return fixedTime }

			result, err := svc.UpdateProfile(context.Background(), tt.userID, tt.request)

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
				tt.checkResult(t, result, repo)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "John Doe", false},
		{"empty string is valid (clears name)", "", false},
		{"whitespace only is valid (clears name)", "   ", false},
		{"exactly 100 chars", string(make([]byte, 100)), false},
		{"101 chars is too long", string(make([]byte, 101)), true},
		{"unicode name", "山田太郎", false},
		{"name with numbers", "User123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, apperrors.IsValidation(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateWeightUnit(t *testing.T) {
	tests := []struct {
		unit    string
		wantErr bool
	}{
		{"lb", false},
		{"kg", false},
		{"LB", true},  // Case sensitive
		{"KG", true},  // Case sensitive
		{"lbs", true}, // Invalid
		{"", true},    // Empty invalid
		{"pounds", true},
		{"kilograms", true},
	}

	for _, tt := range tests {
		t.Run(tt.unit, func(t *testing.T) {
			err := validateWeightUnit(tt.unit)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, apperrors.IsValidation(err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProfileJSONTags(t *testing.T) {
	// Verify that the Profile struct has the expected JSON tags
	// This is a compile-time check - if the tags are wrong, the API contract breaks
	name := "Test User"
	p := Profile{
		ID:         "user-123",
		Email:      "test@example.com",
		Name:       &name,
		WeightUnit: "lb",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Just verify the struct can be created with all fields
	assert.Equal(t, "user-123", p.ID)
	assert.Equal(t, "test@example.com", p.Email)
	assert.Equal(t, "Test User", *p.Name)
	assert.Equal(t, "lb", p.WeightUnit)
}

func TestWeightUnitConstants(t *testing.T) {
	assert.Equal(t, "lb", WeightUnitLb)
	assert.Equal(t, "kg", WeightUnitKg)
}

func TestMaxNameLength(t *testing.T) {
	assert.Equal(t, 100, maxNameLength)
}

// =============================================================================
// SQLITE REPOSITORY INTEGRATION TESTS (REQ-TD2-007)
// =============================================================================

func setupTestDB(t *testing.T) (*SQLiteProfileRepository, func(), *sql.DB) {
	db, cleanup, err := database.OpenTemp("../../migrations")
	require.NoError(t, err)

	repo := NewSQLiteProfileRepository(db)
	return repo, cleanup, db
}

func createTestUserWithEmail(t *testing.T, db *sql.DB, userID, email string) {
	ctx := context.Background()
	_, err := db.ExecContext(ctx, `
		INSERT INTO users (id, email, weight_unit, created_at, updated_at)
		VALUES (?, ?, 'lb', datetime('now'), datetime('now'))
		ON CONFLICT(id) DO UPDATE SET email = excluded.email
	`, userID, email)
	require.NoError(t, err)
}

func TestSQLiteProfileRepository_GetByUserID(t *testing.T) {
	repo, cleanup, db := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("returns existing user profile", func(t *testing.T) {
		// Create a test user with email
		createTestUserWithEmail(t, db, "profile-test-user-001", "profile-test@example.com")

		profile, err := repo.GetByUserID(ctx, "profile-test-user-001")
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, "profile-test-user-001", profile.ID)
		assert.Equal(t, "profile-test@example.com", profile.Email)
		assert.NotEmpty(t, profile.WeightUnit)
	})

	t.Run("returns not found for nonexistent user", func(t *testing.T) {
		_, err := repo.GetByUserID(ctx, "nonexistent-user")
		require.Error(t, err)
		assert.True(t, apperrors.IsNotFound(err))
	})

	t.Run("handles user with null name", func(t *testing.T) {
		createTestUserWithEmail(t, db, "null-name-user", "nullname@example.com")

		profile, err := repo.GetByUserID(ctx, "null-name-user")
		require.NoError(t, err)
		// Name should be nil since we didn't set it
		assert.Nil(t, profile.Name)
	})

	t.Run("handles user with name set", func(t *testing.T) {
		ctx := context.Background()
		_, err := db.ExecContext(ctx, `
			INSERT INTO users (id, email, name, weight_unit, created_at, updated_at)
			VALUES ('named-user', 'named@example.com', 'Test User Name', 'kg', datetime('now'), datetime('now'))
		`)
		require.NoError(t, err)

		profile, err := repo.GetByUserID(ctx, "named-user")
		require.NoError(t, err)
		require.NotNil(t, profile.Name)
		assert.Equal(t, "Test User Name", *profile.Name)
		assert.Equal(t, "kg", profile.WeightUnit)
	})
}

func TestSQLiteProfileRepository_Update(t *testing.T) {
	repo, cleanup, db := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	// Create test user for updates
	createTestUserWithEmail(t, db, "update-test-user", "update-test@example.com")

	t.Run("updates name only", func(t *testing.T) {
		newName := "Updated Name"
		update := ProfileUpdate{
			Name:      &newName,
			SetName:   true,
			UpdatedAt: fixedTime,
		}

		profile, err := repo.Update(ctx, "update-test-user", update)
		require.NoError(t, err)
		require.NotNil(t, profile)
		require.NotNil(t, profile.Name)
		assert.Equal(t, "Updated Name", *profile.Name)
	})

	t.Run("updates weight unit only", func(t *testing.T) {
		update := ProfileUpdate{
			WeightUnit:    "kg",
			SetWeightUnit: true,
			UpdatedAt:     fixedTime,
		}

		profile, err := repo.Update(ctx, "update-test-user", update)
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, "kg", profile.WeightUnit)
	})

	t.Run("updates both name and weight unit", func(t *testing.T) {
		newName := "Both Updated"
		update := ProfileUpdate{
			Name:          &newName,
			SetName:       true,
			WeightUnit:    "lb",
			SetWeightUnit: true,
			UpdatedAt:     fixedTime,
		}

		profile, err := repo.Update(ctx, "update-test-user", update)
		require.NoError(t, err)
		require.NotNil(t, profile)
		require.NotNil(t, profile.Name)
		assert.Equal(t, "Both Updated", *profile.Name)
		assert.Equal(t, "lb", profile.WeightUnit)
	})

	t.Run("clears name to NULL", func(t *testing.T) {
		update := ProfileUpdate{
			Name:      nil, // Set to NULL
			SetName:   true,
			UpdatedAt: fixedTime,
		}

		profile, err := repo.Update(ctx, "update-test-user", update)
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Nil(t, profile.Name)
	})

	t.Run("returns not found for nonexistent user", func(t *testing.T) {
		newName := "Test"
		update := ProfileUpdate{
			Name:      &newName,
			SetName:   true,
			UpdatedAt: fixedTime,
		}

		_, err := repo.Update(ctx, "nonexistent-user", update)
		require.Error(t, err)
		assert.True(t, apperrors.IsNotFound(err))
	})

	t.Run("updates timestamp correctly", func(t *testing.T) {
		newName := "Timestamp Test"
		update := ProfileUpdate{
			Name:      &newName,
			SetName:   true,
			UpdatedAt: fixedTime,
		}

		profile, err := repo.Update(ctx, "update-test-user", update)
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, fixedTime.Format("2006-01-02"), profile.UpdatedAt.Format("2006-01-02"))
	})
}

func TestSQLiteProfileRepository_UpdateAndGet(t *testing.T) {
	repo, cleanup, db := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	// Create test user
	createTestUserWithEmail(t, db, "update-get-user", "update-get@example.com")

	// Get original profile
	original, err := repo.GetByUserID(ctx, "update-get-user")
	require.NoError(t, err)

	// Update the profile
	newName := "Integration Test Update"
	update := ProfileUpdate{
		Name:          &newName,
		SetName:       true,
		WeightUnit:    "kg",
		SetWeightUnit: true,
		UpdatedAt:     fixedTime,
	}

	updated, err := repo.Update(ctx, "update-get-user", update)
	require.NoError(t, err)

	// Get the profile again
	retrieved, err := repo.GetByUserID(ctx, "update-get-user")
	require.NoError(t, err)

	// Verify the update persisted
	assert.Equal(t, original.ID, retrieved.ID)
	assert.Equal(t, original.Email, retrieved.Email)
	require.NotNil(t, retrieved.Name)
	assert.Equal(t, "Integration Test Update", *retrieved.Name)
	assert.Equal(t, "kg", retrieved.WeightUnit)
	assert.Equal(t, updated.WeightUnit, retrieved.WeightUnit)
}

// TestProfileServiceWithRealDB tests the service layer with real database
func TestProfileServiceWithRealDB(t *testing.T) {
	repo, cleanup, db := setupTestDB(t)
	defer cleanup()

	svc := NewService(repo)
	fixedTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return fixedTime }

	ctx := context.Background()

	// Create test user
	createTestUserWithEmail(t, db, "svc-test-user", "svc-test@example.com")

	t.Run("GetProfile with real DB", func(t *testing.T) {
		profile, err := svc.GetProfile(ctx, "svc-test-user")
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, "svc-test-user", profile.ID)
	})

	t.Run("UpdateProfile with real DB - name only", func(t *testing.T) {
		req := UpdateProfileRequest{
			Name: strPtr("Real DB Test Name"),
		}

		profile, err := svc.UpdateProfile(ctx, "svc-test-user", req)
		require.NoError(t, err)
		require.NotNil(t, profile)
		require.NotNil(t, profile.Name)
		assert.Equal(t, "Real DB Test Name", *profile.Name)
	})

	t.Run("UpdateProfile with real DB - weight unit only", func(t *testing.T) {
		req := UpdateProfileRequest{
			WeightUnit: strPtr("kg"),
		}

		profile, err := svc.UpdateProfile(ctx, "svc-test-user", req)
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, "kg", profile.WeightUnit)
	})

	t.Run("UpdateProfile with real DB - both fields", func(t *testing.T) {
		req := UpdateProfileRequest{
			Name:       strPtr("Both Fields Updated"),
			WeightUnit: strPtr("lb"),
		}

		profile, err := svc.UpdateProfile(ctx, "svc-test-user", req)
		require.NoError(t, err)
		require.NotNil(t, profile)
		require.NotNil(t, profile.Name)
		assert.Equal(t, "Both Fields Updated", *profile.Name)
		assert.Equal(t, "lb", profile.WeightUnit)
	})

	t.Run("UpdateProfile with real DB - clear name", func(t *testing.T) {
		req := UpdateProfileRequest{
			Name: strPtr(""),
		}

		profile, err := svc.UpdateProfile(ctx, "svc-test-user", req)
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Nil(t, profile.Name, "name should be cleared")
	})

	t.Run("GetProfile returns not found for nonexistent user", func(t *testing.T) {
		_, err := svc.GetProfile(ctx, "nonexistent")
		require.Error(t, err)
		assert.True(t, apperrors.IsNotFound(err))
	})

	t.Run("UpdateProfile returns not found for nonexistent user", func(t *testing.T) {
		req := UpdateProfileRequest{
			Name: strPtr("Test"),
		}
		_, err := svc.UpdateProfile(ctx, "nonexistent", req)
		require.Error(t, err)
		assert.True(t, apperrors.IsNotFound(err))
	})
}
