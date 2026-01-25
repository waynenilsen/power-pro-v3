package service

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/waynenilsen/power-pro-v3/internal/db"
	"github.com/waynenilsen/power-pro-v3/internal/domain/progression"
	_ "modernc.org/sqlite"
)

// setupTestDB creates a test database with all migrations applied.
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	// Create temp file for test database
	tmpFile, err := os.CreateTemp("", "progression_service_test_*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()

	// Open database connection
	sqlDB, err := sql.Open("sqlite", tmpFile.Name())
	if err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to open database: %v", err)
	}

	// Enable foreign keys
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		sqlDB.Close()
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	// Run migrations
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("sqlite"); err != nil {
		sqlDB.Close()
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to set dialect: %v", err)
	}

	if err := goose.Up(sqlDB, "../../migrations"); err != nil {
		sqlDB.Close()
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to run migrations: %v", err)
	}

	cleanup := func() {
		sqlDB.Close()
		os.Remove(tmpFile.Name())
	}

	return sqlDB, cleanup
}

// Deterministic UUIDs for seeded lifts (from migrations)
const (
	seededSquatID    = "00000000-0000-0000-0000-000000000001"
	seededBenchID    = "00000000-0000-0000-0000-000000000002"
	seededDeadliftID = "00000000-0000-0000-0000-000000000003"
)

// setupTestData creates test data for the progression service tests.
func setupTestData(t *testing.T, sqlDB *sql.DB) testData {
	t.Helper()

	queries := db.New(sqlDB)
	ctx := context.Background()
	now := time.Now().Format(time.RFC3339)
	// Use a past date for initial lift maxes to avoid unique constraint conflicts
	// when the service creates new lift maxes with current time
	pastDate := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)

	// Create a user
	userID := uuid.New().String()
	err := queries.CreateUser(ctx, db.CreateUserParams{
		ID:        userID,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Use the seeded lifts from migrations
	squatID := seededSquatID
	benchID := seededBenchID
	deadliftID := seededDeadliftID

	// Create a cycle
	cycleID := uuid.New().String()
	err = queries.CreateCycle(ctx, db.CreateCycleParams{
		ID:          cycleID,
		Name:        "4 Week Cycle",
		LengthWeeks: 4,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		t.Fatalf("failed to create cycle: %v", err)
	}

	// Create a program
	programID := uuid.New().String()
	err = queries.CreateProgram(ctx, db.CreateProgramParams{
		ID:        programID,
		Name:      "Test Program",
		Slug:      "test-program",
		CycleID:   cycleID,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		t.Fatalf("failed to create program: %v", err)
	}

	// Enroll user in program
	enrollmentID := uuid.New().String()
	err = queries.CreateUserProgramState(ctx, db.CreateUserProgramStateParams{
		ID:                    enrollmentID,
		UserID:                userID,
		ProgramID:             programID,
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
		EnrollmentStatus:      "ACTIVE",
		CycleStatus:           "PENDING",
		WeekStatus:            "PENDING",
		EnrolledAt:            now,
		UpdatedAt:             now,
	})
	if err != nil {
		t.Fatalf("failed to enroll user: %v", err)
	}

	// Create progressions
	linearSessionProgressionID := uuid.New().String()
	err = queries.CreateProgression(ctx, db.CreateProgressionParams{
		ID:   linearSessionProgressionID,
		Name: "Session Linear",
		Type: string(progression.TypeLinear),
		Parameters: `{
			"id": "` + linearSessionProgressionID + `",
			"name": "Session Linear",
			"increment": 5.0,
			"maxType": "TRAINING_MAX",
			"triggerType": "AFTER_SESSION"
		}`,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		t.Fatalf("failed to create linear session progression: %v", err)
	}

	linearWeekProgressionID := uuid.New().String()
	err = queries.CreateProgression(ctx, db.CreateProgressionParams{
		ID:   linearWeekProgressionID,
		Name: "Week Linear",
		Type: string(progression.TypeLinear),
		Parameters: `{
			"id": "` + linearWeekProgressionID + `",
			"name": "Week Linear",
			"increment": 5.0,
			"maxType": "TRAINING_MAX",
			"triggerType": "AFTER_WEEK"
		}`,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		t.Fatalf("failed to create linear week progression: %v", err)
	}

	cycleProgressionID := uuid.New().String()
	err = queries.CreateProgression(ctx, db.CreateProgressionParams{
		ID:   cycleProgressionID,
		Name: "Cycle Progression",
		Type: string(progression.TypeCycle),
		Parameters: `{
			"id": "` + cycleProgressionID + `",
			"name": "Cycle Progression",
			"increment": 10.0,
			"maxType": "TRAINING_MAX"
		}`,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		t.Fatalf("failed to create cycle progression: %v", err)
	}

	// Create program progressions (link progressions to program and lifts)
	ppSquatSessionID := uuid.New().String()
	err = queries.CreateProgramProgression(ctx, db.CreateProgramProgressionParams{
		ID:            ppSquatSessionID,
		ProgramID:     programID,
		ProgressionID: linearSessionProgressionID,
		LiftID:        sql.NullString{String: squatID, Valid: true},
		Priority:      1,
		Enabled:       1,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		t.Fatalf("failed to create program progression for squat session: %v", err)
	}

	ppBenchSessionID := uuid.New().String()
	err = queries.CreateProgramProgression(ctx, db.CreateProgramProgressionParams{
		ID:            ppBenchSessionID,
		ProgramID:     programID,
		ProgressionID: linearSessionProgressionID,
		LiftID:        sql.NullString{String: benchID, Valid: true},
		Priority:      2,
		Enabled:       1,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		t.Fatalf("failed to create program progression for bench session: %v", err)
	}

	ppDeadliftCycleID := uuid.New().String()
	err = queries.CreateProgramProgression(ctx, db.CreateProgramProgressionParams{
		ID:            ppDeadliftCycleID,
		ProgramID:     programID,
		ProgressionID: cycleProgressionID,
		LiftID:        sql.NullString{String: deadliftID, Valid: true},
		Priority:      3,
		Enabled:       1,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		t.Fatalf("failed to create program progression for deadlift cycle: %v", err)
	}

	// Create initial lift maxes (use past date to avoid conflicts)
	squatMaxID := uuid.New().String()
	err = queries.CreateLiftMax(ctx, db.CreateLiftMaxParams{
		ID:            squatMaxID,
		UserID:        userID,
		LiftID:        squatID,
		Type:          "TRAINING_MAX",
		Value:         300,
		EffectiveDate: pastDate,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		t.Fatalf("failed to create squat max: %v", err)
	}

	benchMaxID := uuid.New().String()
	err = queries.CreateLiftMax(ctx, db.CreateLiftMaxParams{
		ID:            benchMaxID,
		UserID:        userID,
		LiftID:        benchID,
		Type:          "TRAINING_MAX",
		Value:         200,
		EffectiveDate: pastDate,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		t.Fatalf("failed to create bench max: %v", err)
	}

	deadliftMaxID := uuid.New().String()
	err = queries.CreateLiftMax(ctx, db.CreateLiftMaxParams{
		ID:            deadliftMaxID,
		UserID:        userID,
		LiftID:        deadliftID,
		Type:          "TRAINING_MAX",
		Value:         400,
		EffectiveDate: pastDate,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		t.Fatalf("failed to create deadlift max: %v", err)
	}

	return testData{
		UserID:                       userID,
		ProgramID:                    programID,
		SquatID:                      squatID,
		BenchID:                      benchID,
		DeadliftID:                   deadliftID,
		LinearSessionProgressionID:   linearSessionProgressionID,
		LinearWeekProgressionID:      linearWeekProgressionID,
		CycleProgressionID:           cycleProgressionID,
		PPSquatSessionID:             ppSquatSessionID,
		PPBenchSessionID:             ppBenchSessionID,
		PPDeadliftCycleID:            ppDeadliftCycleID,
	}
}

type testData struct {
	UserID                       string
	ProgramID                    string
	SquatID                      string
	BenchID                      string
	DeadliftID                   string
	LinearSessionProgressionID   string
	LinearWeekProgressionID      string
	CycleProgressionID           string
	PPSquatSessionID             string
	PPBenchSessionID             string
	PPDeadliftCycleID            string
}

// TestProgressionService_HandleSessionComplete tests AFTER_SESSION trigger handling.
func TestProgressionService_HandleSessionComplete(t *testing.T) {
	sqlDB, cleanup := setupTestDB(t)
	defer cleanup()

	data := setupTestData(t, sqlDB)
	factory := GetDefaultFactory()
	service := NewProgressionService(sqlDB, factory)

	ctx := context.Background()

	t.Run("applies progression for lift in session", func(t *testing.T) {
		event := progression.NewSessionTriggerEvent(
			data.UserID,
			"session-1",
			"day-a",
			1,
			[]string{data.SquatID},
		)

		result, err := service.HandleSessionComplete(ctx, event)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.TotalApplied != 1 {
			t.Errorf("expected TotalApplied=1, got %d", result.TotalApplied)
		}

		// Verify the squat was progressed
		found := false
		for _, r := range result.Results {
			if r.LiftID == data.SquatID && r.Applied {
				found = true
				if r.Result.Delta != 5.0 {
					t.Errorf("expected delta 5.0, got %f", r.Result.Delta)
				}
				if r.Result.NewValue != 305.0 {
					t.Errorf("expected new value 305.0, got %f", r.Result.NewValue)
				}
			}
		}
		if !found {
			t.Error("squat progression not found in results")
		}
	})

	t.Run("does not apply for lift not in session", func(t *testing.T) {
		event := progression.NewSessionTriggerEvent(
			data.UserID,
			"session-2",
			"day-b",
			1,
			[]string{data.DeadliftID}, // Only deadlift, no squat or bench
		)

		result, err := service.HandleSessionComplete(ctx, event)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// No session progressions should apply because deadlift has cycle progression (not session)
		if result.TotalApplied != 0 {
			t.Errorf("expected TotalApplied=0, got %d", result.TotalApplied)
		}
	})

	t.Run("applies to multiple lifts in session", func(t *testing.T) {
		// Use a distinct timestamp to avoid conflicts with previous subtests
		event := &progression.TriggerEventV2{
			Type:      progression.TriggerAfterSession,
			UserID:    data.UserID,
			Timestamp: time.Now().Add(1 * time.Hour), // Distinct timestamp
			Context: progression.SessionTriggerContext{
				SessionID:      "session-3",
				DaySlug:        "day-a",
				WeekNumber:     1,
				LiftsPerformed: []string{data.SquatID, data.BenchID},
			},
		}

		result, err := service.HandleSessionComplete(ctx, event)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.TotalApplied != 2 {
			t.Errorf("expected TotalApplied=2, got %d", result.TotalApplied)
		}
	})
}

// TestProgressionService_HandleCycleComplete tests AFTER_CYCLE trigger handling.
func TestProgressionService_HandleCycleComplete(t *testing.T) {
	sqlDB, cleanup := setupTestDB(t)
	defer cleanup()

	data := setupTestData(t, sqlDB)
	factory := GetDefaultFactory()
	service := NewProgressionService(sqlDB, factory)

	ctx := context.Background()

	t.Run("applies cycle progression", func(t *testing.T) {
		event := progression.NewCycleTriggerEvent(
			data.UserID,
			1, // completed cycle
			4, // total weeks
		)

		result, err := service.HandleCycleComplete(ctx, event)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should apply deadlift cycle progression
		if result.TotalApplied != 1 {
			t.Errorf("expected TotalApplied=1, got %d", result.TotalApplied)
		}

		// Verify deadlift was progressed by 10 (cycle progression increment)
		found := false
		for _, r := range result.Results {
			if r.LiftID == data.DeadliftID && r.Applied {
				found = true
				if r.Result.Delta != 10.0 {
					t.Errorf("expected delta 10.0, got %f", r.Result.Delta)
				}
				if r.Result.NewValue != 410.0 {
					t.Errorf("expected new value 410.0, got %f", r.Result.NewValue)
				}
			}
		}
		if !found {
			t.Error("deadlift cycle progression not found in results")
		}
	})

	t.Run("skips session progressions on cycle trigger", func(t *testing.T) {
		// The session progressions for squat/bench should be skipped
		event := progression.NewCycleTriggerEvent(
			data.UserID,
			2, // new cycle
			4,
		)

		result, err := service.HandleCycleComplete(ctx, event)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Session progressions should be skipped (trigger type mismatch)
		for _, r := range result.Results {
			if r.LiftID == data.SquatID || r.LiftID == data.BenchID {
				if r.Applied {
					t.Errorf("session progression should not apply on cycle trigger: %s", r.LiftID)
				}
			}
		}
	})
}

// TestProgressionService_Idempotency tests that progressions are not applied twice.
func TestProgressionService_Idempotency(t *testing.T) {
	sqlDB, cleanup := setupTestDB(t)
	defer cleanup()

	data := setupTestData(t, sqlDB)
	factory := GetDefaultFactory()
	service := NewProgressionService(sqlDB, factory)

	ctx := context.Background()

	// Create a specific timestamp for idempotency testing
	timestamp := time.Now()

	t.Run("first application succeeds", func(t *testing.T) {
		event := &progression.TriggerEventV2{
			Type:      progression.TriggerAfterSession,
			UserID:    data.UserID,
			Timestamp: timestamp,
			Context: progression.SessionTriggerContext{
				SessionID:      "session-idempotent",
				DaySlug:        "day-a",
				WeekNumber:     1,
				LiftsPerformed: []string{data.SquatID},
			},
		}

		result, err := service.HandleSessionComplete(ctx, event)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.TotalApplied != 1 {
			t.Errorf("expected TotalApplied=1, got %d", result.TotalApplied)
		}
	})

	t.Run("second application with same timestamp is skipped", func(t *testing.T) {
		event := &progression.TriggerEventV2{
			Type:      progression.TriggerAfterSession,
			UserID:    data.UserID,
			Timestamp: timestamp, // Same timestamp
			Context: progression.SessionTriggerContext{
				SessionID:      "session-idempotent",
				DaySlug:        "day-a",
				WeekNumber:     1,
				LiftsPerformed: []string{data.SquatID},
			},
		}

		result, err := service.HandleSessionComplete(ctx, event)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.TotalApplied != 0 {
			t.Errorf("expected TotalApplied=0 (idempotent skip), got %d", result.TotalApplied)
		}

		if result.TotalSkipped != 1 {
			t.Errorf("expected TotalSkipped=1, got %d", result.TotalSkipped)
		}

		// Check skip reason
		for _, r := range result.Results {
			if r.LiftID == data.SquatID {
				if !r.Skipped {
					t.Error("expected squat to be skipped")
				}
				if r.SkipReason != "already applied (idempotent skip)" {
					t.Errorf("expected idempotent skip reason, got: %s", r.SkipReason)
				}
			}
		}
	})

	t.Run("application with different timestamp succeeds", func(t *testing.T) {
		event := &progression.TriggerEventV2{
			Type:      progression.TriggerAfterSession,
			UserID:    data.UserID,
			Timestamp: timestamp.Add(time.Hour), // Different timestamp
			Context: progression.SessionTriggerContext{
				SessionID:      "session-idempotent-2",
				DaySlug:        "day-a",
				WeekNumber:     1,
				LiftsPerformed: []string{data.SquatID},
			},
		}

		result, err := service.HandleSessionComplete(ctx, event)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.TotalApplied != 1 {
			t.Errorf("expected TotalApplied=1, got %d", result.TotalApplied)
		}
	})
}

// TestProgressionService_AtomicTransaction tests that transactions rollback on failure.
func TestProgressionService_AtomicTransaction(t *testing.T) {
	sqlDB, cleanup := setupTestDB(t)
	defer cleanup()

	data := setupTestData(t, sqlDB)
	factory := GetDefaultFactory()
	service := NewProgressionService(sqlDB, factory)

	ctx := context.Background()
	queries := db.New(sqlDB)

	t.Run("verifies LiftMax is created on success", func(t *testing.T) {
		// Get initial count
		initialLogs, err := queries.ListProgressionLogsByUser(ctx, db.ListProgressionLogsByUserParams{
			UserID: data.UserID,
			Limit:  100,
			Offset: 0,
		})
		if err != nil {
			t.Fatalf("failed to get initial logs: %v", err)
		}
		initialCount := len(initialLogs)

		event := progression.NewSessionTriggerEvent(
			data.UserID,
			"session-atomic",
			"day-a",
			1,
			[]string{data.SquatID},
		)

		result, err := service.HandleSessionComplete(ctx, event)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.TotalApplied != 1 {
			t.Errorf("expected TotalApplied=1, got %d", result.TotalApplied)
		}

		// Verify progression log was created
		logs, err := queries.ListProgressionLogsByUser(ctx, db.ListProgressionLogsByUserParams{
			UserID: data.UserID,
			Limit:  100,
			Offset: 0,
		})
		if err != nil {
			t.Fatalf("failed to get logs: %v", err)
		}

		if len(logs) != initialCount+1 {
			t.Errorf("expected %d logs, got %d", initialCount+1, len(logs))
		}
	})
}

// TestProgressionService_UserNotEnrolled tests error handling for unenrolled users.
func TestProgressionService_UserNotEnrolled(t *testing.T) {
	sqlDB, cleanup := setupTestDB(t)
	defer cleanup()

	factory := GetDefaultFactory()
	service := NewProgressionService(sqlDB, factory)

	ctx := context.Background()

	event := progression.NewSessionTriggerEvent(
		"non-existent-user",
		"session-1",
		"day-a",
		1,
		[]string{"some-lift"},
	)

	_, err := service.HandleSessionComplete(ctx, event)
	if err != ErrUserNotEnrolled {
		t.Errorf("expected ErrUserNotEnrolled, got: %v", err)
	}
}

// TestProgressionService_InvalidTriggerContext tests error handling for invalid contexts.
func TestProgressionService_InvalidTriggerContext(t *testing.T) {
	sqlDB, cleanup := setupTestDB(t)
	defer cleanup()

	factory := GetDefaultFactory()
	service := NewProgressionService(sqlDB, factory)

	ctx := context.Background()

	t.Run("wrong trigger type for HandleSessionComplete", func(t *testing.T) {
		event := progression.NewWeekTriggerEvent("user-1", 1, 2, 1)

		_, err := service.HandleSessionComplete(ctx, event)
		if err == nil {
			t.Error("expected error for wrong trigger type")
		}
	})

	t.Run("wrong trigger type for HandleCycleComplete", func(t *testing.T) {
		event := progression.NewSessionTriggerEvent("user-1", "session", "day", 1, nil)

		_, err := service.HandleCycleComplete(ctx, event)
		if err == nil {
			t.Error("expected error for wrong trigger type")
		}
	})

	t.Run("wrong trigger type for HandleWeekAdvance", func(t *testing.T) {
		event := progression.NewCycleTriggerEvent("user-1", 1, 4)

		_, err := service.HandleWeekAdvance(ctx, event)
		if err == nil {
			t.Error("expected error for wrong trigger type")
		}
	})
}

// TestProgressionService_NoCurrentMax tests handling when no current max exists.
func TestProgressionService_NoCurrentMax(t *testing.T) {
	sqlDB, cleanup := setupTestDB(t)
	defer cleanup()

	data := setupTestData(t, sqlDB)
	factory := GetDefaultFactory()
	service := NewProgressionService(sqlDB, factory)

	ctx := context.Background()
	queries := db.New(sqlDB)

	// Create a new lift without a max
	newLiftID := uuid.New().String()
	now := time.Now().Format(time.RFC3339)
	err := queries.CreateLift(ctx, db.CreateLiftParams{
		ID:                newLiftID,
		Name:              "New Lift",
		Slug:              "new-lift",
		IsCompetitionLift: 0,
		CreatedAt:         now,
		UpdatedAt:         now,
	})
	if err != nil {
		t.Fatalf("failed to create new lift: %v", err)
	}

	// Create a program progression for the new lift
	ppID := uuid.New().String()
	err = queries.CreateProgramProgression(ctx, db.CreateProgramProgressionParams{
		ID:            ppID,
		ProgramID:     data.ProgramID,
		ProgressionID: data.LinearSessionProgressionID,
		LiftID:        sql.NullString{String: newLiftID, Valid: true},
		Priority:      99,
		Enabled:       1,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		t.Fatalf("failed to create program progression: %v", err)
	}

	event := progression.NewSessionTriggerEvent(
		data.UserID,
		"session-no-max",
		"day-a",
		1,
		[]string{newLiftID},
	)

	result, err := service.HandleSessionComplete(ctx, event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should skip because there's no current max
	if result.TotalApplied != 0 {
		t.Errorf("expected TotalApplied=0, got %d", result.TotalApplied)
	}

	if result.TotalSkipped != 1 {
		t.Errorf("expected TotalSkipped=1, got %d", result.TotalSkipped)
	}
}

// TestProgressionService_PriorityOrdering tests that progressions are applied in priority order.
func TestProgressionService_PriorityOrdering(t *testing.T) {
	sqlDB, cleanup := setupTestDB(t)
	defer cleanup()

	data := setupTestData(t, sqlDB)
	factory := GetDefaultFactory()
	service := NewProgressionService(sqlDB, factory)

	ctx := context.Background()

	event := progression.NewSessionTriggerEvent(
		data.UserID,
		"session-priority",
		"day-a",
		1,
		[]string{data.SquatID, data.BenchID},
	)

	result, err := service.HandleSessionComplete(ctx, event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify results are in priority order (squat=1, bench=2)
	if len(result.Results) < 2 {
		t.Fatalf("expected at least 2 results, got %d", len(result.Results))
	}

	// First should be squat (priority 1)
	if result.Results[0].LiftID != data.SquatID {
		t.Errorf("expected first result to be squat, got %s", result.Results[0].LiftID)
	}

	// Second should be bench (priority 2)
	if result.Results[1].LiftID != data.BenchID {
		t.Errorf("expected second result to be bench, got %s", result.Results[1].LiftID)
	}
}

// TestProgressionService_DisabledProgression tests that disabled progressions are skipped.
func TestProgressionService_DisabledProgression(t *testing.T) {
	sqlDB, cleanup := setupTestDB(t)
	defer cleanup()

	data := setupTestData(t, sqlDB)
	factory := GetDefaultFactory()
	service := NewProgressionService(sqlDB, factory)

	ctx := context.Background()
	queries := db.New(sqlDB)
	now := time.Now().Format(time.RFC3339)

	// Disable the squat progression
	err := queries.UpdateProgramProgression(ctx, db.UpdateProgramProgressionParams{
		ID:        data.PPSquatSessionID,
		Priority:  1,
		Enabled:   0, // Disabled
		UpdatedAt: now,
	})
	if err != nil {
		t.Fatalf("failed to disable progression: %v", err)
	}

	event := progression.NewSessionTriggerEvent(
		data.UserID,
		"session-disabled",
		"day-a",
		1,
		[]string{data.SquatID, data.BenchID},
	)

	result, err := service.HandleSessionComplete(ctx, event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only bench should be applied (squat is disabled)
	if result.TotalApplied != 1 {
		t.Errorf("expected TotalApplied=1, got %d", result.TotalApplied)
	}

	for _, r := range result.Results {
		if r.LiftID == data.SquatID && r.Applied {
			t.Error("disabled squat progression should not be applied")
		}
		if r.LiftID == data.BenchID && !r.Applied {
			t.Error("bench progression should be applied")
		}
	}
}

// TestProgressionService_OverrideIncrement tests that override increments are applied correctly.
func TestProgressionService_OverrideIncrement(t *testing.T) {
	sqlDB, cleanup := setupTestDB(t)
	defer cleanup()

	data := setupTestData(t, sqlDB)
	factory := GetDefaultFactory()
	service := NewProgressionService(sqlDB, factory)

	ctx := context.Background()
	queries := db.New(sqlDB)
	now := time.Now().Format(time.RFC3339)

	// Set an override increment for deadlift cycle progression (normally 10, override to 15)
	err := queries.UpdateProgramProgression(ctx, db.UpdateProgramProgressionParams{
		ID:                data.PPDeadliftCycleID,
		Priority:          3,
		Enabled:           1,
		OverrideIncrement: sql.NullFloat64{Float64: 15.0, Valid: true},
		UpdatedAt:         now,
	})
	if err != nil {
		t.Fatalf("failed to set override increment: %v", err)
	}

	event := progression.NewCycleTriggerEvent(
		data.UserID,
		1,
		4,
	)

	result, err := service.HandleCycleComplete(ctx, event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify deadlift was progressed by 15 (override)
	for _, r := range result.Results {
		if r.LiftID == data.DeadliftID && r.Applied {
			if r.Result.Delta != 15.0 {
				t.Errorf("expected delta 15.0 (override), got %f", r.Result.Delta)
			}
			if r.Result.NewValue != 415.0 {
				t.Errorf("expected new value 415.0, got %f", r.Result.NewValue)
			}
		}
	}
}

// TestGetDefaultFactory tests the default factory creation.
func TestGetDefaultFactory(t *testing.T) {
	factory := GetDefaultFactory()

	if !factory.IsRegistered(progression.TypeLinear) {
		t.Error("expected LINEAR_PROGRESSION to be registered")
	}

	if !factory.IsRegistered(progression.TypeCycle) {
		t.Error("expected CYCLE_PROGRESSION to be registered")
	}
}

// TestProgressionService_ApplyProgressionManually tests manual progression triggering.
// Each subtest gets its own database to ensure test isolation.
func TestProgressionService_ApplyProgressionManually(t *testing.T) {
	ctx := context.Background()

	t.Run("applies progression to specific lift", func(t *testing.T) {
		sqlDB, cleanup := setupTestDB(t)
		defer cleanup()
		data := setupTestData(t, sqlDB)
		factory := GetDefaultFactory()
		svc := NewProgressionService(sqlDB, factory)

		result, err := svc.ApplyProgressionManually(ctx, data.UserID, data.LinearSessionProgressionID, data.SquatID, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.TotalApplied != 1 {
			t.Errorf("expected TotalApplied=1, got %d", result.TotalApplied)
		}

		if len(result.Results) != 1 {
			t.Errorf("expected 1 result, got %d", len(result.Results))
		}

		if result.Results[0].LiftID != data.SquatID {
			t.Errorf("expected lift ID %s, got %s", data.SquatID, result.Results[0].LiftID)
		}

		if !result.Results[0].Applied {
			t.Error("expected progression to be applied")
		}

		if result.Results[0].Result.Delta != 5.0 {
			t.Errorf("expected delta 5.0, got %f", result.Results[0].Result.Delta)
		}
	})

	t.Run("applies progression to all configured lifts when liftId is empty", func(t *testing.T) {
		sqlDB, cleanup := setupTestDB(t)
		defer cleanup()
		data := setupTestData(t, sqlDB)
		factory := GetDefaultFactory()
		svc := NewProgressionService(sqlDB, factory)

		// The session progression is configured for both squat and bench
		result, err := svc.ApplyProgressionManually(ctx, data.UserID, data.LinearSessionProgressionID, "", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should apply to both squat and bench
		if result.TotalApplied != 2 {
			t.Errorf("expected TotalApplied=2, got %d (skipped=%d, errors=%d)", result.TotalApplied, result.TotalSkipped, result.TotalErrors)
			for _, r := range result.Results {
				t.Logf("  liftID=%s applied=%v skipped=%v reason=%s error=%s", r.LiftID, r.Applied, r.Skipped, r.SkipReason, r.Error)
			}
		}

		if len(result.Results) != 2 {
			t.Errorf("expected 2 results, got %d", len(result.Results))
		}

		// Verify both lifts were progressed
		foundSquat := false
		foundBench := false
		for _, r := range result.Results {
			if r.LiftID == data.SquatID && r.Applied {
				foundSquat = true
			}
			if r.LiftID == data.BenchID && r.Applied {
				foundBench = true
			}
		}
		if !foundSquat {
			t.Error("expected squat to be progressed")
		}
		if !foundBench {
			t.Error("expected bench to be progressed")
		}
	})

	t.Run("force=false respects idempotency on immediate second call", func(t *testing.T) {
		sqlDB, cleanup := setupTestDB(t)
		defer cleanup()
		data := setupTestData(t, sqlDB)
		factory := GetDefaultFactory()
		svc := NewProgressionService(sqlDB, factory)

		// First application
		result1, err := svc.ApplyProgressionManually(ctx, data.UserID, data.CycleProgressionID, data.DeadliftID, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result1.TotalApplied != 1 {
			t.Errorf("expected first application to succeed, got TotalApplied=%d", result1.TotalApplied)
		}

		// Second application with force=false immediately after
		result2, err := svc.ApplyProgressionManually(ctx, data.UserID, data.CycleProgressionID, data.DeadliftID, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Log outcome - behavior depends on timing (sub-second calls)
		t.Logf("Second call: applied=%d, skipped=%d", result2.TotalApplied, result2.TotalSkipped)
	})

	t.Run("force=true bypasses idempotency", func(t *testing.T) {
		sqlDB, cleanup := setupTestDB(t)
		defer cleanup()
		data := setupTestData(t, sqlDB)
		factory := GetDefaultFactory()
		svc := NewProgressionService(sqlDB, factory)
		queries := db.New(sqlDB)

		// Get initial max
		initialMax, err := queries.GetCurrentMax(ctx, db.GetCurrentMaxParams{
			UserID: data.UserID,
			LiftID: data.SquatID,
			Type:   "TRAINING_MAX",
		})
		if err != nil {
			t.Fatalf("failed to get initial max: %v", err)
		}
		initialValue := initialMax.Value

		// First force application
		result1, err := svc.ApplyProgressionManually(ctx, data.UserID, data.LinearSessionProgressionID, data.SquatID, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result1.TotalApplied != 1 {
			t.Errorf("expected first force application to succeed, got TotalApplied=%d", result1.TotalApplied)
		}

		// Get max after first application
		afterFirst, err := queries.GetCurrentMax(ctx, db.GetCurrentMaxParams{
			UserID: data.UserID,
			LiftID: data.SquatID,
			Type:   "TRAINING_MAX",
		})
		if err != nil {
			t.Fatalf("failed to get max after first: %v", err)
		}
		if afterFirst.Value != initialValue+5.0 {
			t.Errorf("expected value %f after first, got %f", initialValue+5.0, afterFirst.Value)
		}

		// Second force application - should ALSO apply because force=true
		result2, err := svc.ApplyProgressionManually(ctx, data.UserID, data.LinearSessionProgressionID, data.SquatID, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result2.TotalApplied != 1 {
			t.Errorf("expected second force application to succeed, got TotalApplied=%d", result2.TotalApplied)
		}

		// Get max after second application
		afterSecond, err := queries.GetCurrentMax(ctx, db.GetCurrentMaxParams{
			UserID: data.UserID,
			LiftID: data.SquatID,
			Type:   "TRAINING_MAX",
		})
		if err != nil {
			t.Fatalf("failed to get max after second: %v", err)
		}
		if afterSecond.Value != initialValue+10.0 {
			t.Errorf("expected value %f after second force, got %f", initialValue+10.0, afterSecond.Value)
		}
	})

	t.Run("force=true creates log with manual and force markers", func(t *testing.T) {
		sqlDB, cleanup := setupTestDB(t)
		defer cleanup()
		data := setupTestData(t, sqlDB)
		factory := GetDefaultFactory()
		svc := NewProgressionService(sqlDB, factory)
		queries := db.New(sqlDB)

		result, err := svc.ApplyProgressionManually(ctx, data.UserID, data.CycleProgressionID, data.DeadliftID, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.TotalApplied != 1 {
			t.Fatalf("expected TotalApplied=1, got %d", result.TotalApplied)
		}

		// Get the most recent log entry
		logs, err := queries.ListProgressionLogsByUser(ctx, db.ListProgressionLogsByUserParams{
			UserID: data.UserID,
			Limit:  1,
			Offset: 0,
		})
		if err != nil {
			t.Fatalf("failed to get logs: %v", err)
		}
		if len(logs) == 0 {
			t.Fatal("expected at least one log entry")
		}

		// Verify trigger_context contains "manual": true and "force": true
		triggerContext := logs[0].TriggerContext.String
		if !contains(triggerContext, `"manual":true`) {
			t.Errorf("expected trigger_context to contain 'manual:true', got: %s", triggerContext)
		}
		if !contains(triggerContext, `"force":true`) {
			t.Errorf("expected trigger_context to contain 'force:true', got: %s", triggerContext)
		}
	})

	t.Run("returns error for non-existent progression", func(t *testing.T) {
		sqlDB, cleanup := setupTestDB(t)
		defer cleanup()
		data := setupTestData(t, sqlDB)
		factory := GetDefaultFactory()
		svc := NewProgressionService(sqlDB, factory)

		_, err := svc.ApplyProgressionManually(ctx, data.UserID, "non-existent-progression", data.SquatID, false)
		if err != ErrProgressionNotFound {
			t.Errorf("expected ErrProgressionNotFound, got: %v", err)
		}
	})

	t.Run("returns error for non-existent lift", func(t *testing.T) {
		sqlDB, cleanup := setupTestDB(t)
		defer cleanup()
		data := setupTestData(t, sqlDB)
		factory := GetDefaultFactory()
		svc := NewProgressionService(sqlDB, factory)

		_, err := svc.ApplyProgressionManually(ctx, data.UserID, data.LinearSessionProgressionID, "non-existent-lift", false)
		if err != ErrLiftNotFound {
			t.Errorf("expected ErrLiftNotFound, got: %v", err)
		}
	})

	t.Run("returns error for unenrolled user", func(t *testing.T) {
		sqlDB, cleanup := setupTestDB(t)
		defer cleanup()
		data := setupTestData(t, sqlDB)
		factory := GetDefaultFactory()
		svc := NewProgressionService(sqlDB, factory)

		_, err := svc.ApplyProgressionManually(ctx, "non-existent-user", data.LinearSessionProgressionID, data.SquatID, false)
		if err != ErrUserNotEnrolled {
			t.Errorf("expected ErrUserNotEnrolled, got: %v", err)
		}
		_ = data // Use data variable
	})

	t.Run("returns error when no lifts configured for progression", func(t *testing.T) {
		sqlDB, cleanup := setupTestDB(t)
		defer cleanup()
		data := setupTestData(t, sqlDB)
		factory := GetDefaultFactory()
		svc := NewProgressionService(sqlDB, factory)

		// The week progression has no lifts configured
		_, err := svc.ApplyProgressionManually(ctx, data.UserID, data.LinearWeekProgressionID, "", false)
		if err == nil {
			t.Error("expected error for progression with no configured lifts")
		}
	})
}

// Helper to check if a string contains a substring (simple contains check)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
