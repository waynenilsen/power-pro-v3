// Package integration provides integration tests for cross-component behavior.
// This file tests E1RM, FindRM, and RelativeTo LoadStrategies working together
// in real program scenarios like GZCL Jacked & Tan 2.0 and Calgary Barbell 8-Week.
package integration

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/e1rm"
	"github.com/waynenilsen/power-pro-v3/internal/domain/liftmax"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loggedset"
	"github.com/waynenilsen/power-pro-v3/internal/domain/rpechart"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
)

// =============================================================================
// MOCK SESSION LOOKUP FOR RELATIVETO TESTS
// =============================================================================

// mockSessionLookup implements loadstrategy.SessionLookup for testing.
type mockSessionLookup struct {
	// sets maps sessionID -> liftID -> setIndex -> LoggedSetResult
	sets map[string]map[string]map[int]*loadstrategy.LoggedSetResult
	err  error
}

func newMockSessionLookup() *mockSessionLookup {
	return &mockSessionLookup{
		sets: make(map[string]map[string]map[int]*loadstrategy.LoggedSetResult),
	}
}

func (m *mockSessionLookup) addSet(sessionID, liftID string, setIndex int, result *loadstrategy.LoggedSetResult) {
	if m.sets[sessionID] == nil {
		m.sets[sessionID] = make(map[string]map[int]*loadstrategy.LoggedSetResult)
	}
	if m.sets[sessionID][liftID] == nil {
		m.sets[sessionID][liftID] = make(map[int]*loadstrategy.LoggedSetResult)
	}
	m.sets[sessionID][liftID][setIndex] = result
}

func (m *mockSessionLookup) GetLoggedSetByIndex(_ context.Context, sessionID string, liftID string, setIndex int) (*loadstrategy.LoggedSetResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.sets[sessionID] == nil {
		return nil, nil
	}
	if m.sets[sessionID][liftID] == nil {
		return nil, nil
	}
	return m.sets[sessionID][liftID][setIndex], nil
}

// =============================================================================
// MOCK MAX LOOKUP WITH E1RM SUPPORT
// =============================================================================

// e1rmMaxLookup implements loadstrategy.MaxLookup with E1RM storage support.
type e1rmMaxLookup struct {
	maxes map[string]*loadstrategy.MaxValue // key: "userID:liftID:maxType"
	err   error
}

func newE1RMMaxLookup() *e1rmMaxLookup {
	return &e1rmMaxLookup{
		maxes: make(map[string]*loadstrategy.MaxValue),
	}
}

func (m *e1rmMaxLookup) SetMax(userID, liftID, maxType string, value float64, date string) {
	key := userID + ":" + liftID + ":" + maxType
	m.maxes[key] = &loadstrategy.MaxValue{
		Value:         value,
		EffectiveDate: date,
	}
}

func (m *e1rmMaxLookup) GetCurrentMax(_ context.Context, userID, liftID, maxType string) (*loadstrategy.MaxValue, error) {
	if m.err != nil {
		return nil, m.err
	}
	key := userID + ":" + liftID + ":" + maxType
	return m.maxes[key], nil
}

// =============================================================================
// E1RM INTEGRATION TESTS
// =============================================================================

// TestE1RMIntegration runs all E1RM integration tests.
func TestE1RMIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("GZCL FindRM with Backoffs", testGZCL_FindRM_WithBackoffs)
	t.Run("Calgary Barbell RPE to E1RM to Backoffs", testCalgaryBarbell_RPE_E1RM_Backoffs)
	t.Run("E1RM Stored as LiftMax", testE1RM_StoredAsLiftMax)
	t.Run("E1RM Used in PercentOf Prescription", testE1RM_UsedInPercentOfPrescription)
	t.Run("FindRM then RelativeTo Multiple Backoff Sets", testFindRM_MultipleBackoffSets)
	t.Run("RPETarget to E1RM Calculation Accuracy", testRPETarget_E1RM_Accuracy)
	t.Run("RelativeTo with Missing Reference Set", testRelativeTo_MissingReferenceSet)
	t.Run("E1RM Weekly Progression", testE1RM_WeeklyProgression)
}

// testGZCL_FindRM_WithBackoffs tests the GZCL Jacked & Tan 2.0 style workflow.
// Week 1: Find 10RM, then 3×10 @ 85% of found weight.
//
// Setup:
//   - User performs FindRM set: 315 × 10
//   - Back-off sets calculated at 85% of 315
//
// Expected:
//   - FindRM strategy returns 0 (user decides weight)
//   - After logging 315 × 10, RelativeTo calculates: 315 × 0.85 = 267.75 → 270 lbs
func testGZCL_FindRM_WithBackoffs(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	squatLiftID := uuid.New().String()

	// Step 1: Create FindRM prescription
	findRMStrategy := loadstrategy.NewFindRMLoadStrategy(10) // Find 10RM

	// Verify FindRM returns 0 (no prescribed weight)
	params := loadstrategy.LoadCalculationParams{
		UserID: userID,
		LiftID: squatLiftID,
	}
	weight, err := findRMStrategy.CalculateLoad(ctx, params)
	if err != nil {
		t.Fatalf("FindRM CalculateLoad failed: %v", err)
	}
	if weight != 0 {
		t.Errorf("FindRM should return 0, got %.1f", weight)
	}

	// Step 2: Simulate user logging their FindRM set: 315 × 10
	sessionLookup := newMockSessionLookup()
	sessionLookup.addSet(sessionID, squatLiftID, 0, &loadstrategy.LoggedSetResult{
		Weight: 315.0,
		Reps:   10,
		RPE:    ptrFloat64(8.5),
	})

	// Step 3: Create RelativeTo strategy for back-offs at 85%
	relativeToStrategy := loadstrategy.NewRelativeToLoadStrategy(
		0,                          // reference set index (first set)
		85.0,                       // percentage
		5.0,                        // rounding increment
		loadstrategy.RoundNearest,  // rounding direction
		sessionLookup,
	)

	// Calculate back-off weight
	backoffParams := loadstrategy.LoadCalculationParams{
		UserID: userID,
		LiftID: squatLiftID,
		Context: map[string]interface{}{
			"sessionID": sessionID,
		},
	}

	backoffWeight, err := relativeToStrategy.CalculateLoad(ctx, backoffParams)
	if err != nil {
		t.Fatalf("RelativeTo CalculateLoad failed: %v", err)
	}

	// Expected: 315 × 0.85 = 267.75 → 270 lbs (rounded to nearest 5)
	expectedBackoff := 270.0
	if math.Abs(backoffWeight-expectedBackoff) > 0.0001 {
		t.Errorf("expected back-off weight %.1f, got %.1f", expectedBackoff, backoffWeight)
	}

	// Step 4: Generate back-off sets
	setScheme, err := setscheme.NewFixedSetScheme(3, 10) // 3x10
	if err != nil {
		t.Fatalf("failed to create set scheme: %v", err)
	}

	sets, err := setScheme.GenerateSets(backoffWeight, setscheme.DefaultSetGenerationContext())
	if err != nil {
		t.Fatalf("failed to generate sets: %v", err)
	}

	// Verify back-off sets
	if len(sets) != 3 {
		t.Errorf("expected 3 back-off sets, got %d", len(sets))
	}
	for i, set := range sets {
		if set.Weight != expectedBackoff {
			t.Errorf("set %d: expected weight %.1f, got %.1f", i+1, expectedBackoff, set.Weight)
		}
		if set.TargetReps != 10 {
			t.Errorf("set %d: expected 10 reps, got %d", i+1, set.TargetReps)
		}
	}
}

// testCalgaryBarbell_RPE_E1RM_Backoffs tests the Calgary Barbell 8-Week style workflow.
// RPE-based top set → calculate E1RM → back-off sets at % of top set.
//
// Setup:
//   - User has squat 1RM = 400 lbs
//   - Prescription: top single @ RPE 8
//   - User logs: 365 × 1 @ RPE 8
//   - Calculate E1RM from logged set
//   - Back-offs at 80% of top set weight
//
// Expected:
//   - RPETarget prescription: 400 × 0.91 (1 rep @ RPE 8) = 364 → 365 lbs
//   - Actual logged: 365 × 1 @ RPE 8
//   - E1RM: 365 / 0.91 = 401.1 → 400 lbs (rounded to 2.5)
//   - Back-offs: 365 × 0.80 = 292 lbs
func testCalgaryBarbell_RPE_E1RM_Backoffs(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	squatLiftID := uuid.New().String()

	// Setup max lookup with 1RM = 400 lbs
	maxLookup := newE1RMMaxLookup()
	maxLookup.SetMax(userID, squatLiftID, "ONE_RM", 400.0, "2024-01-15")

	chart := rpechart.NewDefaultRPEChart()

	// Step 1: Create RPETarget prescription for top single @ RPE 8
	rpeStrategy := loadstrategy.NewRPETargetLoadStrategy(
		1,                          // target reps (single)
		8.0,                        // target RPE
		5.0,                        // rounding increment
		loadstrategy.RoundNearest,
		maxLookup,
		chart,
	)

	params := loadstrategy.LoadCalculationParams{
		UserID: userID,
		LiftID: squatLiftID,
	}

	prescribedWeight, err := rpeStrategy.CalculateLoad(ctx, params)
	if err != nil {
		t.Fatalf("RPETarget CalculateLoad failed: %v", err)
	}

	// Verify prescribed weight: 400 × 0.91 = 364 → 365
	expectedPrescribed := 365.0
	if math.Abs(prescribedWeight-expectedPrescribed) > 0.0001 {
		t.Errorf("expected prescribed weight %.1f, got %.1f", expectedPrescribed, prescribedWeight)
	}

	// Step 2: Simulate user logging the top set: 365 × 1 @ RPE 8
	actualWeight := 365.0
	actualReps := 1
	actualRPE := 8.0

	sessionLookup := newMockSessionLookup()
	sessionLookup.addSet(sessionID, squatLiftID, 0, &loadstrategy.LoggedSetResult{
		Weight: actualWeight,
		Reps:   actualReps,
		RPE:    &actualRPE,
	})

	// Step 3: Calculate E1RM from the logged set
	calculator := e1rm.NewCalculator(chart)
	calculatedE1RM, err := calculator.Calculate(actualWeight, actualReps, actualRPE)
	if err != nil {
		t.Fatalf("E1RM calculation failed: %v", err)
	}

	// Expected E1RM: 365 / 0.91 = 401.1 → 400.0 (rounded to 2.5)
	expectedE1RM := 400.0
	if math.Abs(calculatedE1RM-expectedE1RM) > 0.0001 {
		t.Errorf("expected E1RM %.1f, got %.1f", expectedE1RM, calculatedE1RM)
	}

	// Step 4: Calculate back-offs at 80% of top set weight
	relativeToStrategy := loadstrategy.NewRelativeToLoadStrategy(
		0,                          // reference set index (top set)
		80.0,                       // percentage
		5.0,                        // rounding increment
		loadstrategy.RoundNearest,
		sessionLookup,
	)

	backoffParams := loadstrategy.LoadCalculationParams{
		UserID: userID,
		LiftID: squatLiftID,
		Context: map[string]interface{}{
			"sessionID": sessionID,
		},
	}

	backoffWeight, err := relativeToStrategy.CalculateLoad(ctx, backoffParams)
	if err != nil {
		t.Fatalf("RelativeTo CalculateLoad failed: %v", err)
	}

	// Expected back-off: 365 × 0.80 = 292 → 290 (rounded to nearest 5)
	expectedBackoff := 290.0
	if math.Abs(backoffWeight-expectedBackoff) > 0.0001 {
		t.Errorf("expected back-off weight %.1f, got %.1f", expectedBackoff, backoffWeight)
	}
}

// testE1RM_StoredAsLiftMax tests storing calculated E1RM as a LiftMax entity.
//
// Setup:
//   - Calculate E1RM from logged set: 290 × 5 @ RPE 8
//   - Store as LiftMax with type E1RM
//
// Expected:
//   - E1RM: 290 / 0.77 = 376.6 → 377.5 lbs
//   - LiftMax entity created with type E1RM
func testE1RM_StoredAsLiftMax(t *testing.T) {
	// Setup
	chart := rpechart.NewDefaultRPEChart()
	calculator := e1rm.NewCalculator(chart)

	// Calculate E1RM from a logged set
	weight := 290.0
	reps := 5
	rpe := 8.0

	calculatedE1RM, err := calculator.Calculate(weight, reps, rpe)
	if err != nil {
		t.Fatalf("E1RM calculation failed: %v", err)
	}

	// Expected: 290 / 0.77 = 376.6 → 377.5
	expectedE1RM := 377.5
	if math.Abs(calculatedE1RM-expectedE1RM) > 0.0001 {
		t.Errorf("expected E1RM %.1f, got %.1f", expectedE1RM, calculatedE1RM)
	}

	// Store as LiftMax with type E1RM
	effectiveDate := time.Now()
	input := liftmax.CreateLiftMaxInput{
		UserID:        uuid.New().String(),
		LiftID:        uuid.New().String(),
		Type:          liftmax.E1RM,
		Value:         calculatedE1RM,
		EffectiveDate: &effectiveDate,
	}

	storedE1RM, result := liftmax.CreateLiftMax(input, uuid.New().String(), nil)
	if !result.Valid {
		t.Fatalf("failed to create E1RM LiftMax: %v", result.Errors)
	}

	// Verify stored values
	if storedE1RM.Type != liftmax.E1RM {
		t.Errorf("expected type %s, got %s", liftmax.E1RM, storedE1RM.Type)
	}
	if storedE1RM.Value != calculatedE1RM {
		t.Errorf("expected value %.1f, got %.1f", calculatedE1RM, storedE1RM.Value)
	}
}

// testE1RM_UsedInPercentOfPrescription tests using stored E1RM in a PercentOf prescription.
//
// Setup:
//   - E1RM stored as 400 lbs
//   - Prescription: 75% of E1RM
//
// Expected:
//   - Calculated weight: 400 × 0.75 = 300 lbs
func testE1RM_UsedInPercentOfPrescription(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New().String()
	liftID := uuid.New().String()

	// Setup max lookup with E1RM = 400 lbs
	maxLookup := newE1RMMaxLookup()
	maxLookup.SetMax(userID, liftID, "E1RM", 400.0, "2024-01-15")

	// Create PercentOf strategy referencing E1RM
	percentOfStrategy := loadstrategy.NewPercentOfLoadStrategy(
		loadstrategy.ReferenceE1RM, // reference E1RM instead of ONE_RM
		75.0,                        // percentage
		5.0,                         // rounding increment
		loadstrategy.RoundNearest,
		maxLookup,
	)

	params := loadstrategy.LoadCalculationParams{
		UserID: userID,
		LiftID: liftID,
	}

	weight, err := percentOfStrategy.CalculateLoad(ctx, params)
	if err != nil {
		t.Fatalf("PercentOf CalculateLoad failed: %v", err)
	}

	// Expected: 400 × 0.75 = 300 lbs
	expectedWeight := 300.0
	if math.Abs(weight-expectedWeight) > 0.0001 {
		t.Errorf("expected weight %.1f, got %.1f", expectedWeight, weight)
	}
}

// testFindRM_MultipleBackoffSets tests FindRM followed by multiple back-off set variations.
//
// Setup:
//   - Find 8RM: User achieves 340 × 8
//   - Back-off 1: 3×8 @ 85% = 289 → 290 lbs
//   - Back-off 2: 3×8 @ 80% = 272 → 270 lbs
func testFindRM_MultipleBackoffSets(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	liftID := uuid.New().String()

	// Setup session lookup with the FindRM result
	sessionLookup := newMockSessionLookup()
	sessionLookup.addSet(sessionID, liftID, 0, &loadstrategy.LoggedSetResult{
		Weight: 340.0,
		Reps:   8,
		RPE:    ptrFloat64(9.0),
	})

	testCases := []struct {
		name           string
		percentage     float64
		expectedWeight float64
	}{
		{
			name:           "85% back-off",
			percentage:     85.0,
			expectedWeight: 290.0, // 340 × 0.85 = 289 → 290
		},
		{
			name:           "80% back-off",
			percentage:     80.0,
			expectedWeight: 270.0, // 340 × 0.80 = 272 → 270
		},
		{
			name:           "75% back-off",
			percentage:     75.0,
			expectedWeight: 255.0, // 340 × 0.75 = 255
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			strategy := loadstrategy.NewRelativeToLoadStrategy(
				0,
				tc.percentage,
				5.0,
				loadstrategy.RoundNearest,
				sessionLookup,
			)

			params := loadstrategy.LoadCalculationParams{
				UserID: userID,
				LiftID: liftID,
				Context: map[string]interface{}{
					"sessionID": sessionID,
				},
			}

			weight, err := strategy.CalculateLoad(ctx, params)
			if err != nil {
				t.Fatalf("CalculateLoad failed: %v", err)
			}

			if math.Abs(weight-tc.expectedWeight) > 0.0001 {
				t.Errorf("expected %.1f, got %.1f", tc.expectedWeight, weight)
			}
		})
	}
}

// testRPETarget_E1RM_Accuracy tests E1RM calculation accuracy across different rep/RPE combinations.
func testRPETarget_E1RM_Accuracy(t *testing.T) {
	chart := rpechart.NewDefaultRPEChart()
	calculator := e1rm.NewCalculator(chart)

	testCases := []struct {
		name       string
		weight     float64
		reps       int
		rpe        float64
		expectedE1RM float64
	}{
		{
			name:       "Heavy single @ RPE 10 (true 1RM)",
			weight:     400.0,
			reps:       1,
			rpe:        10.0,
			expectedE1RM: 400.0, // 400 / 1.00 = 400
		},
		{
			name:       "Single @ RPE 8",
			weight:     364.0,
			reps:       1,
			rpe:        8.0,
			expectedE1RM: 400.0, // 364 / 0.91 = 400
		},
		{
			name:       "5 reps @ RPE 9",
			weight:     320.0,
			reps:       5,
			rpe:        9.0,
			expectedE1RM: 400.0, // 320 / 0.80 = 400
		},
		{
			name:       "8 reps @ RPE 8",
			weight:     264.0,
			reps:       8,
			rpe:        8.0,
			expectedE1RM: 400.0, // 264 / 0.66 = 400
		},
		{
			name:       "3 reps @ RPE 9.5",
			weight:     362.0,
			reps:       3,
			rpe:        9.5,
			expectedE1RM: 400.0, // 362 / 0.905 = 400
		},
		{
			name:       "10 reps @ RPE 7",
			weight:     240.0,
			reps:       10,
			rpe:        7.0,
			expectedE1RM: 400.0, // 240 / 0.60 = 400
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e1rm, err := calculator.Calculate(tc.weight, tc.reps, tc.rpe)
			if err != nil {
				t.Fatalf("E1RM calculation failed: %v", err)
			}

			if math.Abs(e1rm-tc.expectedE1RM) > 0.0001 {
				t.Errorf("expected E1RM %.1f, got %.1f", tc.expectedE1RM, e1rm)
			}
		})
	}
}

// testRelativeTo_MissingReferenceSet tests error handling when the reference set isn't logged yet.
func testRelativeTo_MissingReferenceSet(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New().String()
	sessionID := uuid.New().String()
	liftID := uuid.New().String()

	// Empty session lookup - no sets logged
	sessionLookup := newMockSessionLookup()

	strategy := loadstrategy.NewRelativeToLoadStrategy(
		0,
		85.0,
		5.0,
		loadstrategy.RoundNearest,
		sessionLookup,
	)

	params := loadstrategy.LoadCalculationParams{
		UserID: userID,
		LiftID: liftID,
		Context: map[string]interface{}{
			"sessionID": sessionID,
		},
	}

	_, err := strategy.CalculateLoad(ctx, params)
	if err == nil {
		t.Error("expected error when reference set not found, got nil")
	}
}

// testE1RM_WeeklyProgression tests E1RM calculation and storage over multiple weeks.
// This simulates a Calgary Barbell style progression where E1RM is tracked weekly.
func testE1RM_WeeklyProgression(t *testing.T) {
	chart := rpechart.NewDefaultRPEChart()
	calculator := e1rm.NewCalculator(chart)

	userID := uuid.New().String()
	liftID := uuid.New().String()

	// Simulate 4 weeks of progression
	// 1 rep @ RPE 8 = 0.91, 1 rep @ RPE 8.5 = 0.93
	weeklyData := []struct {
		week         int
		weight       float64
		reps         int
		rpe          float64
		expectedE1RM float64
	}{
		{week: 1, weight: 355.0, reps: 1, rpe: 8.0, expectedE1RM: 390.0},   // 355 / 0.91 = 390.1 → 390
		{week: 2, weight: 360.0, reps: 1, rpe: 8.0, expectedE1RM: 395.0},   // 360 / 0.91 = 395.6 → 395
		{week: 3, weight: 365.0, reps: 1, rpe: 8.0, expectedE1RM: 400.0},   // 365 / 0.91 = 401.1 → 400
		{week: 4, weight: 370.0, reps: 1, rpe: 8.5, expectedE1RM: 397.5},   // 370 / 0.93 = 397.8 → 397.5
	}

	var storedE1RMs []*liftmax.LiftMax

	for _, week := range weeklyData {
		// Calculate E1RM for this week
		e1rmValue, err := calculator.Calculate(week.weight, week.reps, week.rpe)
		if err != nil {
			t.Fatalf("week %d: E1RM calculation failed: %v", week.week, err)
		}

		if math.Abs(e1rmValue-week.expectedE1RM) > 0.0001 {
			t.Errorf("week %d: expected E1RM %.2f, got %.2f", week.week, week.expectedE1RM, e1rmValue)
		}

		// Store E1RM
		effectiveDate := time.Now().AddDate(0, 0, (week.week-1)*7)
		input := liftmax.CreateLiftMaxInput{
			UserID:        userID,
			LiftID:        liftID,
			Type:          liftmax.E1RM,
			Value:         e1rmValue,
			EffectiveDate: &effectiveDate,
		}

		storedE1RM, result := liftmax.CreateLiftMax(input, uuid.New().String(), nil)
		if !result.Valid {
			t.Fatalf("week %d: failed to create E1RM LiftMax: %v", week.week, result.Errors)
		}

		storedE1RMs = append(storedE1RMs, storedE1RM)
	}

	// Verify we tracked 4 weeks of progression
	if len(storedE1RMs) != 4 {
		t.Errorf("expected 4 stored E1RMs, got %d", len(storedE1RMs))
	}

	// Verify progression trend (E1RM should generally increase)
	for i := 1; i < len(storedE1RMs); i++ {
		if storedE1RMs[i].Value < storedE1RMs[i-1].Value-5 { // Allow small fluctuations
			t.Logf("Note: E1RM decreased from week %d (%.1f) to week %d (%.1f)",
				i, storedE1RMs[i-1].Value, i+1, storedE1RMs[i].Value)
		}
	}
}

// =============================================================================
// E1RM ACCEPTANCE TESTS
// =============================================================================

// TestE1RMAcceptanceCriteria validates E1RM business requirements.
func TestE1RMAcceptanceCriteria(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}

	// AC1: GZCL FindRM→RelativeTo flow works end-to-end
	t.Run("AC1: GZCL FindRM to RelativeTo flow", func(t *testing.T) {
		ctx := context.Background()
		sessionID := uuid.New().String()
		liftID := uuid.New().String()

		// FindRM returns 0
		findRM := loadstrategy.NewFindRMLoadStrategy(10)
		weight, _ := findRM.CalculateLoad(ctx, loadstrategy.LoadCalculationParams{
			UserID: uuid.New().String(),
			LiftID: liftID,
		})
		if weight != 0 {
			t.Errorf("FindRM should return 0, got %.1f", weight)
		}

		// After logging, RelativeTo calculates correctly
		sessionLookup := newMockSessionLookup()
		sessionLookup.addSet(sessionID, liftID, 0, &loadstrategy.LoggedSetResult{
			Weight: 300.0,
			Reps:   10,
		})

		relativeTo := loadstrategy.NewRelativeToLoadStrategy(0, 85.0, 5.0, loadstrategy.RoundNearest, sessionLookup)
		backoff, err := relativeTo.CalculateLoad(ctx, loadstrategy.LoadCalculationParams{
			UserID: uuid.New().String(),
			LiftID: liftID,
			Context: map[string]interface{}{
				"sessionID": sessionID,
			},
		})
		if err != nil {
			t.Fatalf("RelativeTo failed: %v", err)
		}

		// 300 × 0.85 = 255
		if backoff != 255.0 {
			t.Errorf("expected 255.0, got %.1f", backoff)
		}
	})

	// AC2: Calgary Barbell RPE→E1RM→RelativeTo flow works
	t.Run("AC2: Calgary Barbell RPE to E1RM to RelativeTo flow", func(t *testing.T) {
		ctx := context.Background()
		sessionID := uuid.New().String()
		liftID := uuid.New().String()
		userID := uuid.New().String()

		// Setup 1RM for RPETarget
		maxLookup := newE1RMMaxLookup()
		maxLookup.SetMax(userID, liftID, "ONE_RM", 400.0, "2024-01-15")

		chart := rpechart.NewDefaultRPEChart()

		// RPETarget prescribes weight
		// 1 rep @ RPE 8 = 0.91, so 400 × 0.91 = 364 → 365
		rpeTarget := loadstrategy.NewRPETargetLoadStrategy(1, 8.0, 5.0, loadstrategy.RoundNearest, maxLookup, chart)
		prescribed, _ := rpeTarget.CalculateLoad(ctx, loadstrategy.LoadCalculationParams{
			UserID: userID,
			LiftID: liftID,
		})
		if prescribed != 365.0 { // 400 × 0.91 = 364 → 365
			t.Errorf("expected prescribed 365.0, got %.1f", prescribed)
		}

		// User logs actual set
		sessionLookup := newMockSessionLookup()
		sessionLookup.addSet(sessionID, liftID, 0, &loadstrategy.LoggedSetResult{
			Weight: 365.0,
			Reps:   1,
			RPE:    ptrFloat64(8.0),
		})

		// Calculate E1RM
		// 365 / 0.91 = 401.1 → 400 (rounded to 2.5)
		calculator := e1rm.NewCalculator(chart)
		e1rmValue, _ := calculator.Calculate(365.0, 1, 8.0)
		if math.Abs(e1rmValue-400.0) > 0.0001 {
			t.Errorf("expected E1RM 400.0, got %.1f", e1rmValue)
		}

		// Back-offs at 80%
		relativeTo := loadstrategy.NewRelativeToLoadStrategy(0, 80.0, 5.0, loadstrategy.RoundNearest, sessionLookup)
		backoff, _ := relativeTo.CalculateLoad(ctx, loadstrategy.LoadCalculationParams{
			UserID: userID,
			LiftID: liftID,
			Context: map[string]interface{}{
				"sessionID": sessionID,
			},
		})
		if backoff != 290.0 { // 365 × 0.80 = 292 → 290
			t.Errorf("expected backoff 290.0, got %.1f", backoff)
		}
	})

	// AC3: E1RM can be stored and used in future prescriptions
	t.Run("AC3: E1RM stored and used in PercentOf", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New().String()
		liftID := uuid.New().String()

		// Calculate and store E1RM
		chart := rpechart.NewDefaultRPEChart()
		calculator := e1rm.NewCalculator(chart)
		e1rmValue, _ := calculator.Calculate(320.0, 5, 9.0) // 320 / 0.80 = 400

		effectiveDate := time.Now()
		input := liftmax.CreateLiftMaxInput{
			UserID:        userID,
			LiftID:        liftID,
			Type:          liftmax.E1RM,
			Value:         e1rmValue,
			EffectiveDate: &effectiveDate,
		}
		stored, result := liftmax.CreateLiftMax(input, uuid.New().String(), nil)
		if !result.Valid {
			t.Fatalf("failed to store E1RM: %v", result.Errors)
		}

		// Use stored E1RM in PercentOf
		maxLookup := newE1RMMaxLookup()
		maxLookup.SetMax(userID, liftID, "E1RM", stored.Value, "2024-01-15")

		percentOf := loadstrategy.NewPercentOfLoadStrategy(
			loadstrategy.ReferenceE1RM,
			80.0,
			5.0,
			loadstrategy.RoundNearest,
			maxLookup,
		)

		weight, err := percentOf.CalculateLoad(ctx, loadstrategy.LoadCalculationParams{
			UserID: userID,
			LiftID: liftID,
		})
		if err != nil {
			t.Fatalf("PercentOf failed: %v", err)
		}

		// 400 × 0.80 = 320
		if weight != 320.0 {
			t.Errorf("expected 320.0, got %.1f", weight)
		}
	})

	// AC4: All edge cases handled
	t.Run("AC4: Edge cases handled", func(t *testing.T) {
		ctx := context.Background()

		// Missing session ID
		relativeTo := loadstrategy.NewRelativeToLoadStrategy(0, 85.0, 5.0, loadstrategy.RoundNearest, newMockSessionLookup())
		_, err := relativeTo.CalculateLoad(ctx, loadstrategy.LoadCalculationParams{
			UserID: uuid.New().String(),
			LiftID: uuid.New().String(),
			// Missing Context with sessionID
		})
		if err == nil {
			t.Error("expected error for missing sessionID")
		}

		// Invalid RPE for E1RM
		chart := rpechart.NewDefaultRPEChart()
		calculator := e1rm.NewCalculator(chart)
		_, err = calculator.Calculate(300.0, 5, 6.0) // RPE 6 is invalid (must be 7-10)
		if err == nil {
			t.Error("expected error for invalid RPE")
		}

		// Invalid reps for E1RM
		_, err = calculator.Calculate(300.0, 15, 8.0) // 15 reps is invalid (must be 1-12)
		if err == nil {
			t.Error("expected error for invalid reps")
		}

		// Zero weight for E1RM
		_, err = calculator.Calculate(0, 5, 8.0)
		if err == nil {
			t.Error("expected error for zero weight")
		}
	})

	// AC5: Tests are integration-level (touch multiple domains)
	t.Run("AC5: Integration touches multiple domains", func(t *testing.T) {
		ctx := context.Background()

		// This test exercises:
		// - loadstrategy (FindRM, RelativeTo)
		// - e1rm (Calculator)
		// - liftmax (LiftMax entity)
		// - rpechart (RPEChart)
		// - loggedset (LoggedSet)
		// - setscheme (SetScheme)

		userID := uuid.New().String()
		sessionID := uuid.New().String()
		prescriptionID := uuid.New().String()
		liftID := uuid.New().String()

		// 1. FindRM strategy
		findRM := loadstrategy.NewFindRMLoadStrategy(10)
		if findRM.Type() != loadstrategy.TypeFindRM {
			t.Error("incorrect strategy type")
		}

		// 2. Log a set with RPE
		rpe := 8.5
		input := loggedset.CreateLoggedSetInput{
			UserID:         userID,
			SessionID:      sessionID,
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      1,
			Weight:         290.0,
			TargetReps:     10,
			RepsPerformed:  10,
			IsAMRAP:        false,
			RPE:            &rpe,
		}
		logged, valResult := loggedset.NewLoggedSet(input, uuid.New().String())
		if !valResult.Valid {
			t.Fatalf("failed to create logged set: %v", valResult.Errors)
		}

		// 3. Calculate E1RM from logged set
		chart := rpechart.NewDefaultRPEChart()
		calculator := e1rm.NewCalculator(chart)
		e1rmValue, err := calculator.Calculate(logged.Weight, logged.RepsPerformed, *logged.RPE)
		if err != nil {
			t.Fatalf("E1RM calculation failed: %v", err)
		}

		// 4. Store E1RM as LiftMax
		effectiveDate := time.Now()
		liftMaxInput := liftmax.CreateLiftMaxInput{
			UserID:        userID,
			LiftID:        liftID,
			Type:          liftmax.E1RM,
			Value:         e1rmValue,
			EffectiveDate: &effectiveDate,
		}
		storedMax, maxResult := liftmax.CreateLiftMax(liftMaxInput, uuid.New().String(), nil)
		if !maxResult.Valid {
			t.Fatalf("failed to create LiftMax: %v", maxResult.Errors)
		}

		// 5. RelativeTo uses logged set
		sessionLookup := newMockSessionLookup()
		sessionLookup.addSet(sessionID, liftID, 0, &loadstrategy.LoggedSetResult{
			Weight: logged.Weight,
			Reps:   logged.RepsPerformed,
			RPE:    logged.RPE,
		})

		relativeTo := loadstrategy.NewRelativeToLoadStrategy(0, 85.0, 5.0, loadstrategy.RoundNearest, sessionLookup)
		backoff, err := relativeTo.CalculateLoad(ctx, loadstrategy.LoadCalculationParams{
			UserID: userID,
			LiftID: liftID,
			Context: map[string]interface{}{
				"sessionID": sessionID,
			},
		})
		if err != nil {
			t.Fatalf("RelativeTo failed: %v", err)
		}

		// 6. Generate sets with SetScheme
		setScheme, err := setscheme.NewFixedSetScheme(3, 10)
		if err != nil {
			t.Fatalf("failed to create set scheme: %v", err)
		}
		sets, err := setScheme.GenerateSets(backoff, setscheme.DefaultSetGenerationContext())
		if err != nil {
			t.Fatalf("failed to generate sets: %v", err)
		}

		// Verify full integration
		if len(sets) != 3 {
			t.Errorf("expected 3 sets, got %d", len(sets))
		}
		if storedMax.Type != liftmax.E1RM {
			t.Error("stored max should be E1RM type")
		}
		t.Logf("Integration complete: E1RM=%.1f, backoff=%.1f, sets=%d", storedMax.Value, backoff, len(sets))
	})
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// ptrFloat64 returns a pointer to the given float64 value.
func ptrFloat64(v float64) *float64 {
	return &v
}
