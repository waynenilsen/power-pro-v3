// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/userprogramstate"
)

// UserProgramStateRepository implements user program state persistence using sqlc-generated queries.
type UserProgramStateRepository struct {
	queries *db.Queries
}

// NewUserProgramStateRepository creates a new UserProgramStateRepository.
func NewUserProgramStateRepository(sqlDB *sql.DB) *UserProgramStateRepository {
	return &UserProgramStateRepository{
		queries: db.New(sqlDB),
	}
}

// GetByUserID retrieves a user's program state by their user ID.
func (r *UserProgramStateRepository) GetByUserID(userID string) (*userprogramstate.UserProgramState, error) {
	ctx := context.Background()
	dbState, err := r.queries.GetUserProgramStateByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user program state: %w", err)
	}
	return dbUserProgramStateToDomain(dbState), nil
}

// GetByID retrieves a user's program state by its ID.
func (r *UserProgramStateRepository) GetByID(id string) (*userprogramstate.UserProgramState, error) {
	ctx := context.Background()
	dbState, err := r.queries.GetUserProgramStateByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user program state by ID: %w", err)
	}
	return dbUserProgramStateToDomain(dbState), nil
}

// GetEnrollmentWithProgram retrieves a user's enrollment along with program details.
func (r *UserProgramStateRepository) GetEnrollmentWithProgram(userID string) (*userprogramstate.EnrollmentWithProgram, error) {
	ctx := context.Background()
	row, err := r.queries.GetEnrollmentWithProgram(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get enrollment with program: %w", err)
	}

	enrolledAt, _ := time.Parse(time.RFC3339, row.EnrolledAt)
	updatedAt, _ := time.Parse(time.RFC3339, row.UpdatedAt)

	state := &userprogramstate.UserProgramState{
		ID:                    row.ID,
		UserID:                row.UserID,
		ProgramID:             row.ProgramID,
		CurrentWeek:           int(row.CurrentWeek),
		CurrentCycleIteration: int(row.CurrentCycleIteration),
		CurrentDayIndex:       nullInt64ToIntPtr(row.CurrentDayIndex),
		EnrolledAt:            enrolledAt,
		UpdatedAt:             updatedAt,
	}

	return &userprogramstate.EnrollmentWithProgram{
		State:              state,
		ProgramName:        row.ProgramName,
		ProgramSlug:        row.ProgramSlug,
		ProgramDescription: nullStringToStringPtr(row.ProgramDescription),
		CycleLengthWeeks:   int(row.CycleLengthWeeks),
	}, nil
}

// Create persists a new user program state to the database.
func (r *UserProgramStateRepository) Create(state *userprogramstate.UserProgramState) error {
	ctx := context.Background()

	err := r.queries.CreateUserProgramState(ctx, db.CreateUserProgramStateParams{
		ID:                    state.ID,
		UserID:                state.UserID,
		ProgramID:             state.ProgramID,
		CurrentWeek:           int64(state.CurrentWeek),
		CurrentCycleIteration: int64(state.CurrentCycleIteration),
		CurrentDayIndex:       intPtrToNullInt64(state.CurrentDayIndex),
		EnrolledAt:            state.EnrolledAt.Format(time.RFC3339),
		UpdatedAt:             state.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to create user program state: %w", err)
	}
	return nil
}

// Update persists changes to an existing user program state.
func (r *UserProgramStateRepository) Update(state *userprogramstate.UserProgramState) error {
	ctx := context.Background()

	err := r.queries.UpdateUserProgramState(ctx, db.UpdateUserProgramStateParams{
		UserID:                state.UserID,
		ProgramID:             state.ProgramID,
		CurrentWeek:           int64(state.CurrentWeek),
		CurrentCycleIteration: int64(state.CurrentCycleIteration),
		CurrentDayIndex:       intPtrToNullInt64(state.CurrentDayIndex),
		UpdatedAt:             state.UpdatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("failed to update user program state: %w", err)
	}
	return nil
}

// DeleteByUserID removes a user's program state from the database.
func (r *UserProgramStateRepository) DeleteByUserID(userID string) error {
	ctx := context.Background()

	err := r.queries.DeleteUserProgramStateByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user program state: %w", err)
	}
	return nil
}

// UserIsEnrolled checks if a user is enrolled in any program.
func (r *UserProgramStateRepository) UserIsEnrolled(userID string) (bool, error) {
	ctx := context.Background()

	isEnrolled, err := r.queries.UserIsEnrolled(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check if user is enrolled: %w", err)
	}
	return isEnrolled == 1, nil
}

// Helper functions

func dbUserProgramStateToDomain(dbState db.UserProgramState) *userprogramstate.UserProgramState {
	enrolledAt, _ := time.Parse(time.RFC3339, dbState.EnrolledAt)
	updatedAt, _ := time.Parse(time.RFC3339, dbState.UpdatedAt)

	return &userprogramstate.UserProgramState{
		ID:                    dbState.ID,
		UserID:                dbState.UserID,
		ProgramID:             dbState.ProgramID,
		CurrentWeek:           int(dbState.CurrentWeek),
		CurrentCycleIteration: int(dbState.CurrentCycleIteration),
		CurrentDayIndex:       nullInt64ToIntPtr(dbState.CurrentDayIndex),
		EnrolledAt:            enrolledAt,
		UpdatedAt:             updatedAt,
	}
}

// Note: intPtrToNullInt64, nullInt64ToIntPtr, stringPtrToNullString, and nullStringToStringPtr
// are shared across repositories and defined in prescription_repository.go and lift_repository.go.
