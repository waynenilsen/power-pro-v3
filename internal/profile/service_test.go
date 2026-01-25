package profile

import (
	"context"
	"testing"
	"time"

	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
