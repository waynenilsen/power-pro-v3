package workout

import (
	"context"
	"errors"
	"testing"

	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/prescription"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
)

// ==================== Mock Types ====================

// mockLoadStrategy implements loadstrategy.LoadStrategy for testing
type mockLoadStrategy struct {
	strategyType    loadstrategy.LoadStrategyType
	calculateResult float64
	calculateErr    error
	validateErr     error
}

func (m *mockLoadStrategy) Type() loadstrategy.LoadStrategyType {
	return m.strategyType
}

func (m *mockLoadStrategy) CalculateLoad(ctx context.Context, params loadstrategy.LoadCalculationParams) (float64, error) {
	if m.calculateErr != nil {
		return 0, m.calculateErr
	}
	return m.calculateResult, nil
}

func (m *mockLoadStrategy) Validate() error {
	return m.validateErr
}

// mockSetScheme implements setscheme.SetScheme for testing
type mockSetScheme struct {
	schemeType     setscheme.SetSchemeType
	generateResult []setscheme.GeneratedSet
	generateErr    error
	validateErr    error
}

func (m *mockSetScheme) Type() setscheme.SetSchemeType {
	return m.schemeType
}

func (m *mockSetScheme) GenerateSets(baseWeight float64, ctx setscheme.SetGenerationContext) ([]setscheme.GeneratedSet, error) {
	if m.generateErr != nil {
		return nil, m.generateErr
	}
	return m.generateResult, nil
}

func (m *mockSetScheme) Validate() error {
	return m.validateErr
}

// mockLiftLookup implements prescription.LiftLookup for testing
type mockLiftLookup struct {
	lifts map[string]*prescription.LiftInfo
	err   error
}

func newMockLiftLookup() *mockLiftLookup {
	return &mockLiftLookup{
		lifts: make(map[string]*prescription.LiftInfo),
	}
}

func (m *mockLiftLookup) GetLiftByID(ctx context.Context, liftID string) (*prescription.LiftInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.lifts[liftID], nil
}

func (m *mockLiftLookup) SetLift(liftID string, lift *prescription.LiftInfo) {
	m.lifts[liftID] = lift
}

// ==================== Helper Functions ====================

func validUUID() string {
	return "550e8400-e29b-41d4-a716-446655440000"
}

func anotherValidUUID() string {
	return "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
}

func createValidPrescription(liftID string) *prescription.Prescription {
	return &prescription.Prescription{
		ID:     anotherValidUUID(),
		LiftID: liftID,
		LoadStrategy: &mockLoadStrategy{
			strategyType:    loadstrategy.TypePercentOf,
			calculateResult: 225.0,
		},
		SetScheme: &mockSetScheme{
			schemeType: setscheme.TypeFixed,
			generateResult: []setscheme.GeneratedSet{
				{SetNumber: 1, Weight: 225.0, TargetReps: 5, IsWorkSet: true},
				{SetNumber: 2, Weight: 225.0, TargetReps: 5, IsWorkSet: true},
				{SetNumber: 3, Weight: 225.0, TargetReps: 5, IsWorkSet: true},
			},
		},
		Order:       1,
		Notes:       "Test notes",
		RestSeconds: intPtr(90),
	}
}

func intPtr(i int) *int {
	return &i
}

// ==================== GenerateWorkout Tests ====================

func TestGenerateWorkout_Success(t *testing.T) {
	userID := "user-123"
	liftID := validUUID()

	liftLookup := newMockLiftLookup()
	liftLookup.SetLift(liftID, &prescription.LiftInfo{
		ID:   liftID,
		Name: "Back Squat",
		Slug: "back-squat",
	})

	prescriptions := []*prescription.Prescription{
		createValidPrescription(liftID),
	}

	programCtx := ProgramContext{
		ProgramID:        "program-1",
		ProgramName:      "Test Program",
		CycleID:          "cycle-1",
		CycleLengthWeeks: 4,
	}

	userState := UserState{
		CurrentWeek:           2,
		CurrentCycleIteration: 1,
	}

	dayCtx := DayContext{
		DayID:   "day-1",
		DaySlug: "heavy-day",
		DayName: "Heavy Day",
	}

	genCtx := GenerationContext{
		LiftLookup:    liftLookup,
		SetGenContext: setscheme.DefaultSetGenerationContext(),
	}

	workout, err := GenerateWorkout(
		context.Background(),
		userID,
		programCtx,
		userState,
		dayCtx,
		prescriptions,
		genCtx,
		"2024-01-15",
	)

	if err != nil {
		t.Fatalf("GenerateWorkout returned error: %v", err)
	}
	if workout == nil {
		t.Fatal("GenerateWorkout returned nil workout")
	}
	if workout.UserID != userID {
		t.Errorf("workout.UserID = %q, want %q", workout.UserID, userID)
	}
	if workout.ProgramID != programCtx.ProgramID {
		t.Errorf("workout.ProgramID = %q, want %q", workout.ProgramID, programCtx.ProgramID)
	}
	if workout.CycleIteration != userState.CurrentCycleIteration {
		t.Errorf("workout.CycleIteration = %d, want %d", workout.CycleIteration, userState.CurrentCycleIteration)
	}
	if workout.WeekNumber != userState.CurrentWeek {
		t.Errorf("workout.WeekNumber = %d, want %d", workout.WeekNumber, userState.CurrentWeek)
	}
	if workout.DaySlug != dayCtx.DaySlug {
		t.Errorf("workout.DaySlug = %q, want %q", workout.DaySlug, dayCtx.DaySlug)
	}
	if workout.Date != "2024-01-15" {
		t.Errorf("workout.Date = %q, want %q", workout.Date, "2024-01-15")
	}
	if len(workout.Exercises) != 1 {
		t.Errorf("len(workout.Exercises) = %d, want 1", len(workout.Exercises))
	}
}

func TestGenerateWorkout_MultipleExercises(t *testing.T) {
	userID := "user-123"
	lift1ID := validUUID()
	lift2ID := anotherValidUUID()

	liftLookup := newMockLiftLookup()
	liftLookup.SetLift(lift1ID, &prescription.LiftInfo{
		ID:   lift1ID,
		Name: "Back Squat",
		Slug: "back-squat",
	})
	liftLookup.SetLift(lift2ID, &prescription.LiftInfo{
		ID:   lift2ID,
		Name: "Bench Press",
		Slug: "bench-press",
	})

	prescriptions := []*prescription.Prescription{
		createValidPrescription(lift1ID),
		createValidPrescription(lift2ID),
	}
	prescriptions[1].ID = "prescription-2"

	programCtx := ProgramContext{
		ProgramID:        "program-1",
		ProgramName:      "Test Program",
		CycleID:          "cycle-1",
		CycleLengthWeeks: 4,
	}

	userState := UserState{
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}

	dayCtx := DayContext{
		DayID:   "day-1",
		DaySlug: "day-a",
		DayName: "Day A",
	}

	genCtx := GenerationContext{
		LiftLookup:    liftLookup,
		SetGenContext: setscheme.DefaultSetGenerationContext(),
	}

	workout, err := GenerateWorkout(
		context.Background(),
		userID,
		programCtx,
		userState,
		dayCtx,
		prescriptions,
		genCtx,
		"2024-01-15",
	)

	if err != nil {
		t.Fatalf("GenerateWorkout returned error: %v", err)
	}
	if len(workout.Exercises) != 2 {
		t.Errorf("len(workout.Exercises) = %d, want 2", len(workout.Exercises))
	}
}

func TestGenerateWorkout_NoPrescriptions(t *testing.T) {
	userID := "user-123"

	programCtx := ProgramContext{
		ProgramID:        "program-1",
		ProgramName:      "Test Program",
		CycleID:          "cycle-1",
		CycleLengthWeeks: 4,
	}

	userState := UserState{
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}

	dayCtx := DayContext{
		DayID:   "day-1",
		DaySlug: "day-a",
		DayName: "Day A",
	}

	genCtx := GenerationContext{
		LiftLookup:    newMockLiftLookup(),
		SetGenContext: setscheme.DefaultSetGenerationContext(),
	}

	_, err := GenerateWorkout(
		context.Background(),
		userID,
		programCtx,
		userState,
		dayCtx,
		[]*prescription.Prescription{}, // Empty prescriptions
		genCtx,
		"2024-01-15",
	)

	if !errors.Is(err, ErrNoPrescriptions) {
		t.Errorf("GenerateWorkout with no prescriptions should return ErrNoPrescriptions, got: %v", err)
	}
}

func TestGenerateWorkout_PrescriptionResolutionError(t *testing.T) {
	userID := "user-123"
	liftID := validUUID()

	// Create a prescription with a strategy that returns an error
	prescriptions := []*prescription.Prescription{
		{
			ID:     anotherValidUUID(),
			LiftID: liftID,
			LoadStrategy: &mockLoadStrategy{
				strategyType: loadstrategy.TypePercentOf,
				calculateErr: errors.New("calculation failed"),
			},
			SetScheme: &mockSetScheme{
				schemeType: setscheme.TypeFixed,
				generateResult: []setscheme.GeneratedSet{
					{SetNumber: 1, Weight: 225.0, TargetReps: 5, IsWorkSet: true},
				},
			},
		},
	}

	// Still need to provide a lift lookup that finds the lift
	liftLookup := newMockLiftLookup()
	liftLookup.SetLift(liftID, &prescription.LiftInfo{
		ID:   liftID,
		Name: "Back Squat",
		Slug: "back-squat",
	})

	programCtx := ProgramContext{
		ProgramID:        "program-1",
		ProgramName:      "Test Program",
		CycleID:          "cycle-1",
		CycleLengthWeeks: 4,
	}

	userState := UserState{
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}

	dayCtx := DayContext{
		DayID:   "day-1",
		DaySlug: "day-a",
		DayName: "Day A",
	}

	genCtx := GenerationContext{
		LiftLookup:    liftLookup,
		SetGenContext: setscheme.DefaultSetGenerationContext(),
	}

	_, err := GenerateWorkout(
		context.Background(),
		userID,
		programCtx,
		userState,
		dayCtx,
		prescriptions,
		genCtx,
		"2024-01-15",
	)

	if err == nil {
		t.Error("GenerateWorkout with failing prescription should return error")
	}
}

func TestGenerateWorkout_ExerciseContainsCorrectInfo(t *testing.T) {
	userID := "user-123"
	liftID := validUUID()
	prescriptionID := anotherValidUUID()

	liftLookup := newMockLiftLookup()
	liftLookup.SetLift(liftID, &prescription.LiftInfo{
		ID:   liftID,
		Name: "Back Squat",
		Slug: "back-squat",
	})

	prescriptions := []*prescription.Prescription{
		{
			ID:     prescriptionID,
			LiftID: liftID,
			LoadStrategy: &mockLoadStrategy{
				strategyType:    loadstrategy.TypePercentOf,
				calculateResult: 200.0,
			},
			SetScheme: &mockSetScheme{
				schemeType: setscheme.TypeFixed,
				generateResult: []setscheme.GeneratedSet{
					{SetNumber: 1, Weight: 200.0, TargetReps: 5, IsWorkSet: true},
					{SetNumber: 2, Weight: 200.0, TargetReps: 5, IsWorkSet: true},
				},
			},
			Notes:       "Form cues",
			RestSeconds: intPtr(120),
		},
	}

	programCtx := ProgramContext{
		ProgramID:        "program-1",
		ProgramName:      "Test Program",
		CycleID:          "cycle-1",
		CycleLengthWeeks: 4,
	}

	userState := UserState{
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}

	dayCtx := DayContext{
		DayID:   "day-1",
		DaySlug: "day-a",
		DayName: "Day A",
	}

	genCtx := GenerationContext{
		LiftLookup:    liftLookup,
		SetGenContext: setscheme.DefaultSetGenerationContext(),
	}

	workout, err := GenerateWorkout(
		context.Background(),
		userID,
		programCtx,
		userState,
		dayCtx,
		prescriptions,
		genCtx,
		"2024-01-15",
	)

	if err != nil {
		t.Fatalf("GenerateWorkout returned error: %v", err)
	}

	if len(workout.Exercises) != 1 {
		t.Fatalf("len(workout.Exercises) = %d, want 1", len(workout.Exercises))
	}

	exercise := workout.Exercises[0]
	if exercise.PrescriptionID != prescriptionID {
		t.Errorf("exercise.PrescriptionID = %q, want %q", exercise.PrescriptionID, prescriptionID)
	}
	if exercise.Lift.ID != liftID {
		t.Errorf("exercise.Lift.ID = %q, want %q", exercise.Lift.ID, liftID)
	}
	if exercise.Lift.Name != "Back Squat" {
		t.Errorf("exercise.Lift.Name = %q, want %q", exercise.Lift.Name, "Back Squat")
	}
	if exercise.Lift.Slug != "back-squat" {
		t.Errorf("exercise.Lift.Slug = %q, want %q", exercise.Lift.Slug, "back-squat")
	}
	if len(exercise.Sets) != 2 {
		t.Errorf("len(exercise.Sets) = %d, want 2", len(exercise.Sets))
	}
	if exercise.Notes != "Form cues" {
		t.Errorf("exercise.Notes = %q, want %q", exercise.Notes, "Form cues")
	}
	if exercise.RestSeconds == nil || *exercise.RestSeconds != 120 {
		t.Errorf("exercise.RestSeconds = %v, want 120", exercise.RestSeconds)
	}
}

func TestGenerateWorkout_SetContainsCorrectInfo(t *testing.T) {
	userID := "user-123"
	liftID := validUUID()

	liftLookup := newMockLiftLookup()
	liftLookup.SetLift(liftID, &prescription.LiftInfo{
		ID:   liftID,
		Name: "Back Squat",
		Slug: "back-squat",
	})

	prescriptions := []*prescription.Prescription{
		{
			ID:     anotherValidUUID(),
			LiftID: liftID,
			LoadStrategy: &mockLoadStrategy{
				strategyType:    loadstrategy.TypePercentOf,
				calculateResult: 185.0,
			},
			SetScheme: &mockSetScheme{
				schemeType: setscheme.TypeRamp,
				generateResult: []setscheme.GeneratedSet{
					{SetNumber: 1, Weight: 135.0, TargetReps: 5, IsWorkSet: false},
					{SetNumber: 2, Weight: 155.0, TargetReps: 3, IsWorkSet: false},
					{SetNumber: 3, Weight: 185.0, TargetReps: 1, IsWorkSet: true},
				},
			},
		},
	}

	programCtx := ProgramContext{
		ProgramID:        "program-1",
		ProgramName:      "Test Program",
		CycleID:          "cycle-1",
		CycleLengthWeeks: 4,
	}

	userState := UserState{
		CurrentWeek:           1,
		CurrentCycleIteration: 1,
	}

	dayCtx := DayContext{
		DayID:   "day-1",
		DaySlug: "day-a",
		DayName: "Day A",
	}

	genCtx := GenerationContext{
		LiftLookup:    liftLookup,
		SetGenContext: setscheme.DefaultSetGenerationContext(),
	}

	workout, err := GenerateWorkout(
		context.Background(),
		userID,
		programCtx,
		userState,
		dayCtx,
		prescriptions,
		genCtx,
		"2024-01-15",
	)

	if err != nil {
		t.Fatalf("GenerateWorkout returned error: %v", err)
	}

	exercise := workout.Exercises[0]
	if len(exercise.Sets) != 3 {
		t.Fatalf("len(exercise.Sets) = %d, want 3", len(exercise.Sets))
	}

	// Verify first set (warmup)
	set1 := exercise.Sets[0]
	if set1.SetNumber != 1 {
		t.Errorf("set1.SetNumber = %d, want 1", set1.SetNumber)
	}
	if set1.Weight != 135.0 {
		t.Errorf("set1.Weight = %f, want 135.0", set1.Weight)
	}
	if set1.TargetReps != 5 {
		t.Errorf("set1.TargetReps = %d, want 5", set1.TargetReps)
	}
	if set1.IsWorkSet {
		t.Error("set1.IsWorkSet should be false")
	}

	// Verify third set (work set)
	set3 := exercise.Sets[2]
	if set3.SetNumber != 3 {
		t.Errorf("set3.SetNumber = %d, want 3", set3.SetNumber)
	}
	if set3.Weight != 185.0 {
		t.Errorf("set3.Weight = %f, want 185.0", set3.Weight)
	}
	if set3.TargetReps != 1 {
		t.Errorf("set3.TargetReps = %d, want 1", set3.TargetReps)
	}
	if !set3.IsWorkSet {
		t.Error("set3.IsWorkSet should be true")
	}
}

func TestGenerateWorkout_WithLookupContext(t *testing.T) {
	userID := "user-123"
	liftID := validUUID()

	liftLookup := newMockLiftLookup()
	liftLookup.SetLift(liftID, &prescription.LiftInfo{
		ID:   liftID,
		Name: "Back Squat",
		Slug: "back-squat",
	})

	prescriptions := []*prescription.Prescription{
		createValidPrescription(liftID),
	}

	programCtx := ProgramContext{
		ProgramID:        "program-1",
		ProgramName:      "Test Program",
		CycleID:          "cycle-1",
		CycleLengthWeeks: 4,
	}

	userState := UserState{
		CurrentWeek:           2,
		CurrentCycleIteration: 1,
	}

	dayCtx := DayContext{
		DayID:   "day-1",
		DaySlug: "heavy-day",
		DayName: "Heavy Day",
	}

	genCtx := GenerationContext{
		LiftLookup:    liftLookup,
		SetGenContext: setscheme.DefaultSetGenerationContext(),
		LookupContext: &loadstrategy.LookupContext{
			WeekNumber: 2,
			DaySlug:    "heavy-day",
		},
	}

	workout, err := GenerateWorkout(
		context.Background(),
		userID,
		programCtx,
		userState,
		dayCtx,
		prescriptions,
		genCtx,
		"2024-01-15",
	)

	if err != nil {
		t.Fatalf("GenerateWorkout returned error: %v", err)
	}
	if workout == nil {
		t.Fatal("GenerateWorkout returned nil workout")
	}
	// The workout should be generated successfully with the lookup context
	if len(workout.Exercises) != 1 {
		t.Errorf("len(workout.Exercises) = %d, want 1", len(workout.Exercises))
	}
}

// ==================== ConvertSets Tests ====================

func TestConvertSets(t *testing.T) {
	sets := []setscheme.GeneratedSet{
		{SetNumber: 1, Weight: 135.0, TargetReps: 8, IsWorkSet: false},
		{SetNumber: 2, Weight: 185.0, TargetReps: 5, IsWorkSet: true},
		{SetNumber: 3, Weight: 185.0, TargetReps: 5, IsWorkSet: true},
	}

	result := convertSets(sets)

	if len(result) != 3 {
		t.Fatalf("len(result) = %d, want 3", len(result))
	}

	// Check first set
	if result[0].SetNumber != 1 {
		t.Errorf("result[0].SetNumber = %d, want 1", result[0].SetNumber)
	}
	if result[0].Weight != 135.0 {
		t.Errorf("result[0].Weight = %f, want 135.0", result[0].Weight)
	}
	if result[0].TargetReps != 8 {
		t.Errorf("result[0].TargetReps = %d, want 8", result[0].TargetReps)
	}
	if result[0].IsWorkSet {
		t.Error("result[0].IsWorkSet should be false")
	}

	// Check second set
	if result[1].SetNumber != 2 {
		t.Errorf("result[1].SetNumber = %d, want 2", result[1].SetNumber)
	}
	if !result[1].IsWorkSet {
		t.Error("result[1].IsWorkSet should be true")
	}
}

func TestConvertSets_Empty(t *testing.T) {
	result := convertSets([]setscheme.GeneratedSet{})
	if len(result) != 0 {
		t.Errorf("len(result) = %d, want 0", len(result))
	}
}

// ==================== GetDateString Tests ====================

func TestGetDateString(t *testing.T) {
	date := GetDateString()

	// Should be in YYYY-MM-DD format
	if len(date) != 10 {
		t.Errorf("GetDateString() returned %q, expected 10 characters", date)
	}
	if date[4] != '-' || date[7] != '-' {
		t.Errorf("GetDateString() returned %q, expected YYYY-MM-DD format", date)
	}
}

// ==================== ValidateGenerationParams Tests ====================

func TestValidateGenerationParams_Valid(t *testing.T) {
	tests := []struct {
		name   string
		params GenerationParams
	}{
		{
			name: "no overrides",
			params: GenerationParams{
				UserID: "user-123",
			},
		},
		{
			name: "with week number",
			params: GenerationParams{
				UserID:     "user-123",
				WeekNumber: intPtr(2),
			},
		},
		{
			name: "with day slug",
			params: GenerationParams{
				UserID:  "user-123",
				DaySlug: strPtr("day-a"),
			},
		},
		{
			name: "with all params",
			params: GenerationParams{
				UserID:     "user-123",
				WeekNumber: intPtr(3),
				DaySlug:    strPtr("heavy"),
				Date:       strPtr("2024-01-15"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGenerationParams(tt.params)
			if err != nil {
				t.Errorf("ValidateGenerationParams() returned error: %v", err)
			}
		})
	}
}

func TestValidateGenerationParams_InvalidWeekNumber(t *testing.T) {
	params := GenerationParams{
		UserID:     "user-123",
		WeekNumber: intPtr(0),
	}

	err := ValidateGenerationParams(params)
	if !errors.Is(err, ErrInvalidWeekNumber) {
		t.Errorf("ValidateGenerationParams() with week 0 should return ErrInvalidWeekNumber, got: %v", err)
	}

	params.WeekNumber = intPtr(-1)
	err = ValidateGenerationParams(params)
	if !errors.Is(err, ErrInvalidWeekNumber) {
		t.Errorf("ValidateGenerationParams() with week -1 should return ErrInvalidWeekNumber, got: %v", err)
	}
}

// ==================== ValidatePreviewParams Tests ====================

func TestValidatePreviewParams_Valid(t *testing.T) {
	params := PreviewParams{
		UserID:     "user-123",
		WeekNumber: 2,
		DaySlug:    "heavy",
	}

	err := ValidatePreviewParams(params)
	if err != nil {
		t.Errorf("ValidatePreviewParams() returned error: %v", err)
	}
}

func TestValidatePreviewParams_InvalidWeekNumber(t *testing.T) {
	tests := []struct {
		name   string
		week   int
	}{
		{"zero", 0},
		{"negative", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := PreviewParams{
				UserID:     "user-123",
				WeekNumber: tt.week,
				DaySlug:    "heavy",
			}

			err := ValidatePreviewParams(params)
			if !errors.Is(err, ErrInvalidWeekNumber) {
				t.Errorf("ValidatePreviewParams() with week %d should return ErrInvalidWeekNumber, got: %v", tt.week, err)
			}
		})
	}
}

func TestValidatePreviewParams_EmptyDaySlug(t *testing.T) {
	params := PreviewParams{
		UserID:     "user-123",
		WeekNumber: 2,
		DaySlug:    "",
	}

	err := ValidatePreviewParams(params)
	if !errors.Is(err, ErrInvalidDaySlug) {
		t.Errorf("ValidatePreviewParams() with empty day slug should return ErrInvalidDaySlug, got: %v", err)
	}
}

// ==================== DefaultGenerationContext Tests ====================

func TestDefaultGenerationContext(t *testing.T) {
	liftLookup := newMockLiftLookup()
	ctx := DefaultGenerationContext(liftLookup)

	if ctx.LiftLookup != liftLookup {
		t.Error("LiftLookup should be set")
	}
	if ctx.SetGenContext.WorkSetThreshold != 80.0 {
		t.Errorf("WorkSetThreshold = %v, want 80.0", ctx.SetGenContext.WorkSetThreshold)
	}
	if ctx.LookupContext != nil {
		t.Error("LookupContext should be nil by default")
	}
}

// ==================== Error Constants Tests ====================

func TestErrorConstants(t *testing.T) {
	if ErrUserNotEnrolled.Error() != "user is not enrolled in any program" {
		t.Errorf("ErrUserNotEnrolled message incorrect: %s", ErrUserNotEnrolled.Error())
	}
	if ErrDayNotFound.Error() != "day not found for the specified position" {
		t.Errorf("ErrDayNotFound message incorrect: %s", ErrDayNotFound.Error())
	}
	if ErrInvalidWeekNumber.Error() != "week number must be >= 1" {
		t.Errorf("ErrInvalidWeekNumber message incorrect: %s", ErrInvalidWeekNumber.Error())
	}
	if ErrInvalidDaySlug.Error() != "day slug is required" {
		t.Errorf("ErrInvalidDaySlug message incorrect: %s", ErrInvalidDaySlug.Error())
	}
	if ErrNoPrescriptions.Error() != "day has no prescriptions" {
		t.Errorf("ErrNoPrescriptions message incorrect: %s", ErrNoPrescriptions.Error())
	}
	if ErrProgramNotFound.Error() != "program not found" {
		t.Errorf("ErrProgramNotFound message incorrect: %s", ErrProgramNotFound.Error())
	}
	if ErrCycleNotFound.Error() != "cycle not found" {
		t.Errorf("ErrCycleNotFound message incorrect: %s", ErrCycleNotFound.Error())
	}
	if ErrWeekNotFound.Error() != "week not found in cycle" {
		t.Errorf("ErrWeekNotFound message incorrect: %s", ErrWeekNotFound.Error())
	}
}

// ==================== Helper Functions ====================

func strPtr(s string) *string {
	return &s
}
