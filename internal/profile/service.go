// Package profile provides user profile management functionality.
// This package handles profile retrieval, updates, and validation.
package profile

import (
	"context"
	"database/sql"
	"strings"
	"time"

	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
)

const (
	// maxNameLength is the maximum allowed length for a user's name.
	maxNameLength = 100
)

// Valid weight units.
const (
	WeightUnitLb = "lb"
	WeightUnitKg = "kg"
)

// Profile represents a user's profile information.
type Profile struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Name       *string   `json:"name"`
	WeightUnit string    `json:"weightUnit"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// UpdateProfileRequest represents a request to update a user's profile.
// Fields are pointers to distinguish between "not provided" (nil) and "set to empty" (empty string).
type UpdateProfileRequest struct {
	// Name is the user's display name. Nil means don't change, empty string means clear.
	Name *string
	// WeightUnit is the user's preferred weight unit ("lb" or "kg"). Nil means don't change.
	WeightUnit *string
}

// ProfileUpdate represents the changes to apply to a profile.
// This is an internal type used by the repository.
type ProfileUpdate struct {
	// Name is the new name value. Only used if SetName is true.
	Name *string
	// SetName indicates whether to update the name field.
	SetName bool
	// WeightUnit is the new weight unit. Only used if SetWeightUnit is true.
	WeightUnit string
	// SetWeightUnit indicates whether to update the weight unit field.
	SetWeightUnit bool
	// UpdatedAt is the timestamp for the update.
	UpdatedAt time.Time
}

// ProfileRepository defines the interface for profile persistence.
type ProfileRepository interface {
	GetByUserID(ctx context.Context, userID string) (*Profile, error)
	Update(ctx context.Context, userID string, update ProfileUpdate) (*Profile, error)
}

// Service provides profile operations.
type Service struct {
	profileRepo ProfileRepository
	now         func() time.Time
}

// NewService creates a new profile service.
func NewService(profileRepo ProfileRepository) *Service {
	return &Service{
		profileRepo: profileRepo,
		now:         time.Now,
	}
}

// GetProfile retrieves a user's profile by their user ID.
func (s *Service) GetProfile(ctx context.Context, userID string) (*Profile, error) {
	if userID == "" {
		return nil, apperrors.NewBadRequest("user ID is required")
	}

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// UpdateProfile updates a user's profile with the provided changes.
// Only non-nil fields in the request will be updated.
func (s *Service) UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*Profile, error) {
	if userID == "" {
		return nil, apperrors.NewBadRequest("user ID is required")
	}

	// Validate name if provided
	if req.Name != nil {
		if err := validateName(*req.Name); err != nil {
			return nil, err
		}
	}

	// Validate weight unit if provided
	if req.WeightUnit != nil {
		if err := validateWeightUnit(*req.WeightUnit); err != nil {
			return nil, err
		}
	}

	// Check if there's anything to update
	if req.Name == nil && req.WeightUnit == nil {
		// Nothing to update, just return the current profile
		return s.profileRepo.GetByUserID(ctx, userID)
	}

	// Build the update
	update := ProfileUpdate{
		UpdatedAt: s.now(),
	}

	// Handle name update - empty string means clear (set to NULL)
	if req.Name != nil {
		update.SetName = true
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			update.Name = nil // Set to NULL
		} else {
			update.Name = &trimmed
		}
	}

	// Handle weight unit update
	if req.WeightUnit != nil {
		update.SetWeightUnit = true
		update.WeightUnit = *req.WeightUnit
	}

	// Update the profile
	profile, err := s.profileRepo.Update(ctx, userID, update)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// validateName validates the user's name.
func validateName(name string) error {
	trimmed := strings.TrimSpace(name)
	// Empty string is valid - it means "clear the name"
	if trimmed == "" {
		return nil
	}
	if len(trimmed) > maxNameLength {
		return apperrors.NewValidation("name", "name must be 100 characters or less")
	}
	return nil
}

// validateWeightUnit validates the weight unit.
func validateWeightUnit(unit string) error {
	if unit != WeightUnitLb && unit != WeightUnitKg {
		return apperrors.NewValidation("weightUnit", "weight unit must be 'lb' or 'kg'")
	}
	return nil
}

// SQLiteProfileRepository implements ProfileRepository using SQLite.
type SQLiteProfileRepository struct {
	db *sql.DB
}

// NewSQLiteProfileRepository creates a new SQLite-backed profile repository.
func NewSQLiteProfileRepository(db *sql.DB) *SQLiteProfileRepository {
	return &SQLiteProfileRepository{db: db}
}

// GetByUserID retrieves a user's profile by their user ID.
func (r *SQLiteProfileRepository) GetByUserID(ctx context.Context, userID string) (*Profile, error) {
	var profile Profile
	var name sql.NullString
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, name, weight_unit, created_at, updated_at
		FROM users WHERE id = ?
	`, userID).Scan(&profile.ID, &profile.Email, &name, &profile.WeightUnit, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, apperrors.NewNotFound("user", userID)
	}
	if err != nil {
		return nil, apperrors.NewInternal("failed to retrieve profile", err)
	}

	if name.Valid {
		profile.Name = &name.String
	}
	profile.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	profile.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &profile, nil
}

// Update updates a user's profile with the provided changes.
func (r *SQLiteProfileRepository) Update(ctx context.Context, userID string, update ProfileUpdate) (*Profile, error) {
	// Build dynamic UPDATE query based on what fields are being updated
	query := "UPDATE users SET updated_at = ?"
	args := []interface{}{update.UpdatedAt.Format(time.RFC3339)}

	if update.SetName {
		if update.Name == nil {
			query += ", name = NULL"
		} else {
			query += ", name = ?"
			args = append(args, *update.Name)
		}
	}

	if update.SetWeightUnit {
		query += ", weight_unit = ?"
		args = append(args, update.WeightUnit)
	}

	query += " WHERE id = ?"
	args = append(args, userID)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, apperrors.NewInternal("failed to update profile", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, apperrors.NewInternal("failed to check update result", err)
	}
	if rowsAffected == 0 {
		return nil, apperrors.NewNotFound("user", userID)
	}

	// Fetch and return the updated profile
	return r.GetByUserID(ctx, userID)
}
