// Package dashboard provides dashboard aggregation functionality.
// This package combines data from multiple sources (enrollment, workouts, maxes)
// into a single dashboard response for frontend consumption.
package dashboard

import (
	"context"
	"database/sql"
	"log"
	"sort"
	"time"

	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/profile"
)

// EnrollmentSummary represents the enrollment section of the dashboard.
type EnrollmentSummary struct {
	Status         string `json:"status"`
	ProgramName    string `json:"programName"`
	CycleIteration int    `json:"cycleIteration"`
	CycleStatus    string `json:"cycleStatus"`
	WeekNumber     int    `json:"weekNumber"`
	WeekStatus     string `json:"weekStatus"`
}

// NextWorkoutPreview represents the next workout section of the dashboard.
type NextWorkoutPreview struct {
	DayName       string `json:"dayName"`
	DaySlug       string `json:"daySlug"`
	ExerciseCount int    `json:"exerciseCount"`
	EstimatedSets int    `json:"estimatedSets"`
}

// SessionSummary represents the current session section of the dashboard.
type SessionSummary struct {
	SessionID     string    `json:"sessionId"`
	DayName       string    `json:"dayName"`
	StartedAt     time.Time `json:"startedAt"`
	SetsCompleted int       `json:"setsCompleted"`
	TotalSets     int       `json:"totalSets"`
}

// WorkoutSummary represents a recent workout in the dashboard.
type WorkoutSummary struct {
	Date          string `json:"date"`
	DayName       string `json:"dayName"`
	SetsCompleted int    `json:"setsCompleted"`
}

// MaxSummary represents a current max in the dashboard.
type MaxSummary struct {
	Lift  string  `json:"lift"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

// Dashboard represents the complete dashboard response.
type Dashboard struct {
	Enrollment     *EnrollmentSummary   `json:"enrollment"`
	NextWorkout    *NextWorkoutPreview  `json:"nextWorkout"`
	CurrentSession *SessionSummary      `json:"currentSession"`
	RecentWorkouts []WorkoutSummary     `json:"recentWorkouts"`
	CurrentMaxes   []MaxSummary         `json:"currentMaxes"`
}

// Service provides dashboard aggregation operations.
type Service struct {
	db             *sql.DB
	queries        *db.Queries
	profileService *profile.Service
}

// NewService creates a new dashboard service.
func NewService(sqlDB *sql.DB, profileService *profile.Service) *Service {
	return &Service{
		db:             sqlDB,
		queries:        db.New(sqlDB),
		profileService: profileService,
	}
}

// GetDashboard retrieves the aggregated dashboard for a user.
func (s *Service) GetDashboard(ctx context.Context, userID string) (*Dashboard, error) {
	dashboard := &Dashboard{
		RecentWorkouts: []WorkoutSummary{},
		CurrentMaxes:   []MaxSummary{},
	}

	// Get user's weight unit preference for max conversion
	weightUnit := "lb" // default
	userProfile, err := s.profileService.GetProfile(ctx, userID)
	if err != nil {
		log.Printf("Warning: failed to get user profile for dashboard: %v", err)
	} else if userProfile != nil {
		weightUnit = userProfile.WeightUnit
	}

	// Get enrollment status
	enrollment, err := s.aggregateEnrollment(ctx, userID)
	if err != nil {
		log.Printf("Warning: failed to get enrollment for dashboard: %v", err)
	}
	dashboard.Enrollment = enrollment

	// Get current session (if any)
	currentSession, err := s.getCurrentSession(ctx, userID)
	if err != nil {
		log.Printf("Warning: failed to get current session for dashboard: %v", err)
	}
	dashboard.CurrentSession = currentSession

	// Get next workout (only if enrolled and no active session)
	if enrollment != nil && currentSession == nil {
		nextWorkout, err := s.calculateNextWorkout(ctx, userID)
		if err != nil {
			log.Printf("Warning: failed to calculate next workout for dashboard: %v", err)
		}
		dashboard.NextWorkout = nextWorkout
	}

	// Get recent workouts
	recentWorkouts, err := s.getRecentWorkouts(ctx, userID, 5)
	if err != nil {
		log.Printf("Warning: failed to get recent workouts for dashboard: %v", err)
	} else {
		dashboard.RecentWorkouts = recentWorkouts
	}

	// Get current maxes
	currentMaxes, err := s.getCurrentMaxes(ctx, userID, weightUnit)
	if err != nil {
		log.Printf("Warning: failed to get current maxes for dashboard: %v", err)
	} else {
		dashboard.CurrentMaxes = currentMaxes
	}

	return dashboard, nil
}

// aggregateEnrollment retrieves the enrollment summary for a user.
func (s *Service) aggregateEnrollment(ctx context.Context, userID string) (*EnrollmentSummary, error) {
	row, err := s.queries.GetEnrollmentWithProgram(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &EnrollmentSummary{
		Status:         row.EnrollmentStatus,
		ProgramName:    row.ProgramName,
		CycleIteration: int(row.CurrentCycleIteration),
		CycleStatus:    row.CycleStatus,
		WeekNumber:     int(row.CurrentWeek),
		WeekStatus:     row.WeekStatus,
	}, nil
}

// getCurrentSession retrieves the current active session for a user.
func (s *Service) getCurrentSession(ctx context.Context, userID string) (*SessionSummary, error) {
	row, err := s.queries.GetActiveWorkoutSessionByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Get day name for the session
	dayName, err := s.getDayNameForSession(ctx, userID, int(row.WeekNumber), int(row.DayIndex))
	if err != nil {
		log.Printf("Warning: failed to get day name for session: %v", err)
		dayName = "Unknown Day"
	}

	// Count sets completed in this session
	setsCompleted, err := s.queries.CountLoggedSetsBySession(ctx, row.ID)
	if err != nil {
		log.Printf("Warning: failed to count logged sets for session: %v", err)
		setsCompleted = 0
	}

	// Get total sets expected for this day
	totalSets, err := s.getTotalSetsForDay(ctx, userID, int(row.WeekNumber), int(row.DayIndex))
	if err != nil {
		log.Printf("Warning: failed to get total sets for day: %v", err)
		totalSets = 0
	}

	startedAt, _ := time.Parse(time.RFC3339, row.StartedAt)

	return &SessionSummary{
		SessionID:     row.ID,
		DayName:       dayName,
		StartedAt:     startedAt,
		SetsCompleted: int(setsCompleted),
		TotalSets:     totalSets,
	}, nil
}

// calculateNextWorkout determines the next workout for a user.
func (s *Service) calculateNextWorkout(ctx context.Context, userID string) (*NextWorkoutPreview, error) {
	// Get the user's enrollment state
	state, err := s.queries.GetUserProgramStateByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Determine the next day index
	nextDayIndex := 0
	if state.CurrentDayIndex.Valid {
		nextDayIndex = int(state.CurrentDayIndex.Int64)
	}

	// Get the day for this position
	day, err := s.getDayForWeekPosition(ctx, state.ProgramID, int(state.CurrentWeek), nextDayIndex)
	if err != nil {
		log.Printf("Warning: failed to get day for week position: %v", err)
		return nil, nil
	}
	if day == nil {
		return nil, nil
	}

	// Count exercises and sets for this day
	exerciseCount, estimatedSets, err := s.getExerciseAndSetCounts(ctx, day.ID)
	if err != nil {
		log.Printf("Warning: failed to get exercise/set counts: %v", err)
	}

	return &NextWorkoutPreview{
		DayName:       day.Name,
		DaySlug:       day.Slug,
		ExerciseCount: exerciseCount,
		EstimatedSets: estimatedSets,
	}, nil
}

// getRecentWorkouts retrieves the recent completed workouts for a user.
func (s *Service) getRecentWorkouts(ctx context.Context, userID string, limit int) ([]WorkoutSummary, error) {
	rows, err := s.queries.GetRecentCompletedWorkouts(ctx, db.GetRecentCompletedWorkoutsParams{
		UserID: userID,
		Limit:  int64(limit),
	})
	if err != nil {
		return nil, err
	}

	summaries := make([]WorkoutSummary, 0, len(rows))
	for _, row := range rows {
		// Parse the finished_at to extract date
		finishedAt, err := time.Parse(time.RFC3339, row.FinishedAt.String)
		if err != nil {
			log.Printf("Warning: failed to parse finished_at: %v", err)
			continue
		}

		summaries = append(summaries, WorkoutSummary{
			Date:          finishedAt.Format("2006-01-02"),
			DayName:       row.DayName,
			SetsCompleted: int(row.SetsCompleted),
		})
	}

	return summaries, nil
}

// getCurrentMaxes retrieves the current training maxes for a user.
func (s *Service) getCurrentMaxes(ctx context.Context, userID string, weightUnit string) ([]MaxSummary, error) {
	rows, err := s.queries.GetCurrentMaxesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	summaries := make([]MaxSummary, 0, len(rows))
	for _, row := range rows {
		value := row.Value

		// Convert to user's preferred unit if needed
		// Maxes are stored in whatever unit the user entered them
		// For now, we assume they're stored in the user's preferred unit
		// Future enhancement: store all maxes in a canonical unit and convert

		summaries = append(summaries, MaxSummary{
			Lift:  row.LiftName,
			Value: value,
			Type:  row.Type,
		})
	}

	// Sort by lift name alphabetically
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Lift < summaries[j].Lift
	})

	return summaries, nil
}

// Helper functions

// getDayNameForSession gets the day name for a session.
func (s *Service) getDayNameForSession(ctx context.Context, userID string, weekNumber, dayIndex int) (string, error) {
	// Get the user's program state to find the program ID
	state, err := s.queries.GetUserProgramStateByUserID(ctx, userID)
	if err != nil {
		return "", err
	}

	day, err := s.getDayForWeekPosition(ctx, state.ProgramID, weekNumber, dayIndex)
	if err != nil {
		return "", err
	}
	if day == nil {
		return "Unknown Day", nil
	}
	return day.Name, nil
}

// dayInfo holds basic day information.
type dayInfo struct {
	ID   string
	Name string
	Slug string
}

// getDayForWeekPosition gets the day at a specific position in a week.
func (s *Service) getDayForWeekPosition(ctx context.Context, programID string, weekNumber, dayIndex int) (*dayInfo, error) {
	row, err := s.queries.GetDayForWeekPosition(ctx, db.GetDayForWeekPositionParams{
		ID:         programID,
		WeekNumber: int64(weekNumber),
		Offset:     int64(dayIndex),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &dayInfo{
		ID:   row.ID,
		Name: row.Name,
		Slug: row.Slug,
	}, nil
}

// getTotalSetsForDay calculates total sets for a day based on prescriptions.
func (s *Service) getTotalSetsForDay(ctx context.Context, userID string, weekNumber, dayIndex int) (int, error) {
	// Get the user's program state to find the program ID
	state, err := s.queries.GetUserProgramStateByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	day, err := s.getDayForWeekPosition(ctx, state.ProgramID, weekNumber, dayIndex)
	if err != nil || day == nil {
		return 0, err
	}

	_, totalSets, err := s.getExerciseAndSetCounts(ctx, day.ID)
	return totalSets, err
}

// getExerciseAndSetCounts counts distinct exercises and total sets for a day.
func (s *Service) getExerciseAndSetCounts(ctx context.Context, dayID string) (int, int, error) {
	row, err := s.queries.GetDayExerciseAndSetCounts(ctx, dayID)
	if err != nil {
		return 0, 0, err
	}

	// TotalSets is interface{} due to COALESCE, need to convert
	totalSets := 0
	switch v := row.TotalSets.(type) {
	case int64:
		totalSets = int(v)
	case int:
		totalSets = v
	case float64:
		totalSets = int(v)
	}

	return int(row.ExerciseCount), totalSets, nil
}
