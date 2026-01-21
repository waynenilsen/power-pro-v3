// Package repository provides database repository implementations.
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/dailylookup"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/prescription"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
	"github.com/waynenilsen/power-pro-v3/internal/domain/weeklylookup"
	"github.com/waynenilsen/power-pro-v3/internal/domain/workout"
)

// WorkoutRepository implements workout generation data access using sqlc-generated queries.
type WorkoutRepository struct {
	queries         *db.Queries
	strategyFactory *loadstrategy.StrategyFactory
	schemeFactory   *setscheme.SchemeFactory
}

// NewWorkoutRepository creates a new WorkoutRepository.
func NewWorkoutRepository(sqlDB *sql.DB, strategyFactory *loadstrategy.StrategyFactory, schemeFactory *setscheme.SchemeFactory) *WorkoutRepository {
	return &WorkoutRepository{
		queries:         db.New(sqlDB),
		strategyFactory: strategyFactory,
		schemeFactory:   schemeFactory,
	}
}

// EnrollmentData contains all data needed for workout generation.
type EnrollmentData struct {
	UserID                string
	ProgramID             string
	ProgramName           string
	ProgramSlug           string
	CycleID               string
	CycleLengthWeeks      int
	WeeklyLookupID        *string
	DailyLookupID         *string
	DefaultRounding       *float64
	CurrentWeek           int
	CurrentCycleIteration int
	CurrentDayIndex       *int
}

// GetEnrollmentForWorkout retrieves the user's enrollment with all program context.
func (r *WorkoutRepository) GetEnrollmentForWorkout(userID string) (*EnrollmentData, error) {
	ctx := context.Background()
	row, err := r.queries.GetEnrollmentForWorkout(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get enrollment for workout: %w", err)
	}

	var weeklyLookupID *string
	if row.WeeklyLookupID.Valid {
		weeklyLookupID = &row.WeeklyLookupID.String
	}

	var dailyLookupID *string
	if row.DailyLookupID.Valid {
		dailyLookupID = &row.DailyLookupID.String
	}

	var defaultRounding *float64
	if row.DefaultRounding.Valid {
		defaultRounding = &row.DefaultRounding.Float64
	}

	var currentDayIndex *int
	if row.CurrentDayIndex.Valid {
		idx := int(row.CurrentDayIndex.Int64)
		currentDayIndex = &idx
	}

	return &EnrollmentData{
		UserID:                row.UserID,
		ProgramID:             row.ProgramID,
		ProgramName:           row.ProgramName,
		ProgramSlug:           row.ProgramSlug,
		CycleID:               row.CycleID,
		CycleLengthWeeks:      int(row.CycleLengthWeeks),
		WeeklyLookupID:        weeklyLookupID,
		DailyLookupID:         dailyLookupID,
		DefaultRounding:       defaultRounding,
		CurrentWeek:           int(row.CurrentWeek),
		CurrentCycleIteration: int(row.CurrentCycleIteration),
		CurrentDayIndex:       currentDayIndex,
	}, nil
}

// GetWeekByNumberAndCycle retrieves a week by its number within a cycle.
func (r *WorkoutRepository) GetWeekByNumberAndCycle(cycleID string, weekNumber int) (string, error) {
	ctx := context.Background()
	week, err := r.queries.GetWeekByNumberAndCycle(ctx, db.GetWeekByNumberAndCycleParams{
		CycleID:    cycleID,
		WeekNumber: int64(weekNumber),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get week: %w", err)
	}
	return week.ID, nil
}

// DayData contains day information for workout generation.
type DayData struct {
	ID   string
	Name string
	Slug string
}

// GetDayBySlugAndWeek retrieves a day by its slug within a week.
func (r *WorkoutRepository) GetDayBySlugAndWeek(slug, weekID string) (*DayData, error) {
	ctx := context.Background()
	day, err := r.queries.GetDayBySlugAndWeek(ctx, db.GetDayBySlugAndWeekParams{
		Slug:   slug,
		WeekID: weekID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get day by slug: %w", err)
	}
	return &DayData{
		ID:   day.ID,
		Name: day.Name,
		Slug: day.Slug,
	}, nil
}

// GetDayByIndexInWeek retrieves a day by its position index within a week.
func (r *WorkoutRepository) GetDayByIndexInWeek(weekID string, index int) (*DayData, error) {
	ctx := context.Background()
	day, err := r.queries.GetDayByIndexInWeek(ctx, db.GetDayByIndexInWeekParams{
		WeekID: weekID,
		Offset: int64(index),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get day by index: %w", err)
	}
	return &DayData{
		ID:   day.ID,
		Name: day.Name,
		Slug: day.Slug,
	}, nil
}

// GetDaysForWeek retrieves all days in a week in order.
func (r *WorkoutRepository) GetDaysForWeek(weekID string) ([]DayData, error) {
	ctx := context.Background()
	rows, err := r.queries.GetDaysForWeek(ctx, weekID)
	if err != nil {
		return nil, fmt.Errorf("failed to get days for week: %w", err)
	}

	days := make([]DayData, len(rows))
	for i, row := range rows {
		days[i] = DayData{
			ID:   row.ID,
			Name: row.Name,
			Slug: row.Slug,
		}
	}
	return days, nil
}

// CountDaysInWeek returns the number of days in a week.
func (r *WorkoutRepository) CountDaysInWeek(weekID string) (int, error) {
	ctx := context.Background()
	count, err := r.queries.CountDaysInWeek(ctx, weekID)
	if err != nil {
		return 0, fmt.Errorf("failed to count days in week: %w", err)
	}
	return int(count), nil
}

// GetPrescriptionsForDay retrieves all prescriptions for a day in order.
func (r *WorkoutRepository) GetPrescriptionsForDay(dayID string) ([]*prescription.Prescription, error) {
	ctx := context.Background()
	dbPrescriptions, err := r.queries.GetPrescriptionsForDay(ctx, dayID)
	if err != nil {
		return nil, fmt.Errorf("failed to get prescriptions for day: %w", err)
	}

	prescriptions := make([]*prescription.Prescription, 0, len(dbPrescriptions))
	for _, dbPrescription := range dbPrescriptions {
		p, err := r.dbPrescriptionToDomain(dbPrescription)
		if err != nil {
			return nil, fmt.Errorf("failed to convert prescription: %w", err)
		}
		prescriptions = append(prescriptions, p)
	}

	return prescriptions, nil
}

// GetWeeklyLookup retrieves a weekly lookup by ID.
func (r *WorkoutRepository) GetWeeklyLookup(id string) (*weeklylookup.WeeklyLookup, error) {
	ctx := context.Background()
	dbLookup, err := r.queries.GetWeeklyLookup(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get weekly lookup: %w", err)
	}
	return dbWeeklyLookupToDomain(dbLookup)
}

// GetDailyLookup retrieves a daily lookup by ID.
func (r *WorkoutRepository) GetDailyLookup(id string) (*dailylookup.DailyLookup, error) {
	ctx := context.Background()
	dbLookup, err := r.queries.GetDailyLookup(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get daily lookup: %w", err)
	}
	return dbDailyLookupToDomain(dbLookup)
}

// LiftLookupAdapter provides lift lookup functionality for prescription resolution.
type LiftLookupAdapter struct {
	queries *db.Queries
}

// NewLiftLookupAdapter creates a new LiftLookupAdapter.
func NewLiftLookupAdapter(sqlDB *sql.DB) *LiftLookupAdapter {
	return &LiftLookupAdapter{
		queries: db.New(sqlDB),
	}
}

// GetLiftByID retrieves lift information by ID.
func (a *LiftLookupAdapter) GetLiftByID(ctx context.Context, liftID string) (*prescription.LiftInfo, error) {
	lift, err := a.queries.GetLift(ctx, liftID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get lift: %w", err)
	}
	return &prescription.LiftInfo{
		ID:   lift.ID,
		Name: lift.Name,
		Slug: lift.Slug,
	}, nil
}

// MaxLookupAdapter provides max lookup functionality for load strategy resolution.
type MaxLookupAdapter struct {
	queries *db.Queries
}

// NewMaxLookupAdapter creates a new MaxLookupAdapter.
func NewMaxLookupAdapter(sqlDB *sql.DB) *MaxLookupAdapter {
	return &MaxLookupAdapter{
		queries: db.New(sqlDB),
	}
}

// GetCurrentMax retrieves the current max for a user, lift, and max type.
func (a *MaxLookupAdapter) GetCurrentMax(ctx context.Context, userID, liftID, maxType string) (*loadstrategy.MaxValue, error) {
	max, err := a.queries.GetCurrentMax(ctx, db.GetCurrentMaxParams{
		UserID: userID,
		LiftID: liftID,
		Type:   maxType,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get current max: %w", err)
	}
	return &loadstrategy.MaxValue{
		Value:         max.Value,
		EffectiveDate: max.EffectiveDate,
	}, nil
}

// InjectMaxLookup injects a MaxLookup into prescriptions that have load strategies supporting it.
func InjectMaxLookup(prescriptions []*prescription.Prescription, maxLookup loadstrategy.MaxLookup) {
	for _, p := range prescriptions {
		// Check if the load strategy has a SetMaxLookup method (like PercentOfLoadStrategy)
		if setter, ok := p.LoadStrategy.(interface{ SetMaxLookup(loadstrategy.MaxLookup) }); ok {
			setter.SetMaxLookup(maxLookup)
		}
	}
}

// Helper function to convert DB prescription to domain
func (r *WorkoutRepository) dbPrescriptionToDomain(dbPrescription db.Prescription) (*prescription.Prescription, error) {
	// Unmarshal load strategy
	loadStrategy, err := r.strategyFactory.CreateFromJSON(json.RawMessage(dbPrescription.LoadStrategy))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal load strategy: %w", err)
	}

	// Unmarshal set scheme
	setScheme, err := r.schemeFactory.CreateFromJSON(json.RawMessage(dbPrescription.SetScheme))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal set scheme: %w", err)
	}

	var restSeconds *int
	if dbPrescription.RestSeconds.Valid {
		val := int(dbPrescription.RestSeconds.Int64)
		restSeconds = &val
	}

	var notes string
	if dbPrescription.Notes.Valid {
		notes = dbPrescription.Notes.String
	}

	return &prescription.Prescription{
		ID:           dbPrescription.ID,
		LiftID:       dbPrescription.LiftID,
		LoadStrategy: loadStrategy,
		SetScheme:    setScheme,
		Order:        int(dbPrescription.Order),
		Notes:        notes,
		RestSeconds:  restSeconds,
	}, nil
}

// WorkoutGenerationData bundles all data needed for workout generation.
type WorkoutGenerationData struct {
	Enrollment    *EnrollmentData
	WeekID        string
	Day           *DayData
	Prescriptions []*prescription.Prescription
	WeeklyLookup  *weeklylookup.WeeklyLookup
	DailyLookup   *dailylookup.DailyLookup
}

// GetWorkoutGenerationData retrieves all data needed for workout generation.
func (r *WorkoutRepository) GetWorkoutGenerationData(userID string, weekNumber *int, daySlug *string) (*WorkoutGenerationData, error) {
	// Get enrollment data
	enrollment, err := r.GetEnrollmentForWorkout(userID)
	if err != nil {
		return nil, err
	}
	if enrollment == nil {
		return nil, workout.ErrUserNotEnrolled
	}

	// Determine week number to use
	targetWeek := enrollment.CurrentWeek
	if weekNumber != nil {
		targetWeek = *weekNumber
	}

	// Get the week by number
	weekID, err := r.GetWeekByNumberAndCycle(enrollment.CycleID, targetWeek)
	if err != nil {
		return nil, err
	}
	if weekID == "" {
		return nil, workout.ErrWeekNotFound
	}

	// Determine day to use
	var day *DayData
	if daySlug != nil {
		// Use specified day slug
		day, err = r.GetDayBySlugAndWeek(*daySlug, weekID)
		if err != nil {
			return nil, err
		}
		if day == nil {
			return nil, workout.ErrDayNotFound
		}
	} else if enrollment.CurrentDayIndex != nil {
		// Use day index from state
		day, err = r.GetDayByIndexInWeek(weekID, *enrollment.CurrentDayIndex)
		if err != nil {
			return nil, err
		}
		if day == nil {
			return nil, workout.ErrDayNotFound
		}
	} else {
		// Default to first day in week
		days, err := r.GetDaysForWeek(weekID)
		if err != nil {
			return nil, err
		}
		if len(days) == 0 {
			return nil, workout.ErrDayNotFound
		}
		day = &days[0]
	}

	// Get prescriptions for the day
	prescriptions, err := r.GetPrescriptionsForDay(day.ID)
	if err != nil {
		return nil, err
	}

	// Get lookups if configured
	var weeklyLookup *weeklylookup.WeeklyLookup
	if enrollment.WeeklyLookupID != nil {
		weeklyLookup, err = r.GetWeeklyLookup(*enrollment.WeeklyLookupID)
		if err != nil {
			return nil, err
		}
	}

	var dailyLookup *dailylookup.DailyLookup
	if enrollment.DailyLookupID != nil {
		dailyLookup, err = r.GetDailyLookup(*enrollment.DailyLookupID)
		if err != nil {
			return nil, err
		}
	}

	// Override week number in enrollment for response
	enrollment.CurrentWeek = targetWeek

	return &WorkoutGenerationData{
		Enrollment:    enrollment,
		WeekID:        weekID,
		Day:           day,
		Prescriptions: prescriptions,
		WeeklyLookup:  weeklyLookup,
		DailyLookup:   dailyLookup,
	}, nil
}
