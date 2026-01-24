// Package integration provides integration tests for cross-component behavior.
// This file tests RTS Intermediate program support with RPE-based load calculation.
package integration

import (
	"context"
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loadstrategy"
	"github.com/waynenilsen/power-pro-v3/internal/domain/loggedset"
	"github.com/waynenilsen/power-pro-v3/internal/domain/rpechart"
	"github.com/waynenilsen/power-pro-v3/internal/domain/setscheme"
)

// =============================================================================
// MOCK MAX LOOKUP FOR RTS TESTS
// =============================================================================

// rtsMaxLookup implements loadstrategy.MaxLookup for testing.
type rtsMaxLookup struct {
	maxes map[string]*loadstrategy.MaxValue // key: "userID:liftID:maxType"
	err   error
}

func newRTSMaxLookup() *rtsMaxLookup {
	return &rtsMaxLookup{
		maxes: make(map[string]*loadstrategy.MaxValue),
	}
}

func (m *rtsMaxLookup) SetMax(userID, liftID, maxType string, value float64, date string) {
	key := userID + ":" + liftID + ":" + maxType
	m.maxes[key] = &loadstrategy.MaxValue{
		Value:         value,
		EffectiveDate: date,
	}
}

func (m *rtsMaxLookup) GetCurrentMax(_ context.Context, userID, liftID, maxType string) (*loadstrategy.MaxValue, error) {
	if m.err != nil {
		return nil, m.err
	}
	key := userID + ":" + liftID + ":" + maxType
	return m.maxes[key], nil
}

// =============================================================================
// RTS INTERMEDIATE E2E TESTS
// =============================================================================

// TestRTSIntermediate runs all RTS Intermediate integration tests.
func TestRTSIntermediate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("Basic RPE-Based Workout Generation", testRTSIntermediateBasicWorkout)
	t.Run("Multiple RPE Prescriptions in Day", testRTSIntermediateMultiplePrescriptions)
	t.Run("RPE with Different Rounding", testRTSIntermediateRounding)
	t.Run("Log Set with RPE", testRTSIntermediateLogSetWithRPE)
	t.Run("Custom RPE Chart", testRTSIntermediateCustomRPEChart)
}

// testRTSIntermediateBasicWorkout tests basic RPE-based weight calculation.
// Setup:
//   - User with squat 1RM = 365 lbs
//   - Prescription with RPETarget strategy (5 reps @ RPE 9)
//
// Expected:
//   - RPE chart: 5 reps @ RPE 9 = 80% (0.80)
//   - Calculated weight: 365 × 0.80 = 292 → rounds to 290 lbs
func testRTSIntermediateBasicWorkout(t *testing.T) {
	ctx := context.Background()

	// Setup user and lift
	userID := uuid.New().String()
	squatLiftID := uuid.New().String()

	// Create mock max lookup with squat 1RM = 365 lbs
	maxLookup := newRTSMaxLookup()
	maxLookup.SetMax(userID, squatLiftID, "ONE_RM", 365.0, "2024-01-15")

	// Create default RPE chart
	chart := rpechart.NewDefaultRPEChart()

	// Create RPE Target strategy: 5 reps @ RPE 9
	strategy := loadstrategy.NewRPETargetLoadStrategy(
		5,                           // target reps
		9.0,                         // target RPE
		5.0,                         // rounding increment
		loadstrategy.RoundNearest,   // rounding direction
		maxLookup,
		chart,
	)

	// Verify strategy type
	if strategy.Type() != loadstrategy.TypeRPETarget {
		t.Errorf("expected strategy type %s, got %s", loadstrategy.TypeRPETarget, strategy.Type())
	}

	// Calculate load
	params := loadstrategy.LoadCalculationParams{
		UserID: userID,
		LiftID: squatLiftID,
	}

	weight, err := strategy.CalculateLoad(ctx, params)
	if err != nil {
		t.Fatalf("CalculateLoad failed: %v", err)
	}

	// Verify: 365 × 0.80 = 292 → rounds to 290 lbs
	expectedWeight := 290.0
	if math.Abs(weight-expectedWeight) > 0.0001 {
		t.Errorf("expected weight %.1f, got %.1f", expectedWeight, weight)
	}

	// Generate sets using fixed set scheme with calculated weight
	setScheme, err := setscheme.NewFixedSetScheme(3, 5) // 3x5
	if err != nil {
		t.Fatalf("failed to create set scheme: %v", err)
	}

	sets, err := setScheme.GenerateSets(weight, setscheme.DefaultSetGenerationContext())
	if err != nil {
		t.Fatalf("failed to generate sets: %v", err)
	}

	// Verify generated sets
	if len(sets) != 3 {
		t.Errorf("expected 3 sets, got %d", len(sets))
	}

	for i, set := range sets {
		if set.Weight != expectedWeight {
			t.Errorf("set %d: expected weight %.1f, got %.1f", i+1, expectedWeight, set.Weight)
		}
		if set.TargetReps != 5 {
			t.Errorf("set %d: expected 5 target reps, got %d", i+1, set.TargetReps)
		}
		if !set.IsWorkSet {
			t.Errorf("set %d: expected IsWorkSet=true", i+1)
		}
	}
}

// testRTSIntermediateMultiplePrescriptions tests multiple lifts with different RPE prescriptions.
// Setup:
//   - Squat: 4 reps @ RPE 9 (82% = 300 lbs from 365 1RM)
//   - Bench: 5 reps @ RPE 8 (77% = 230 lbs from 300 1RM)
func testRTSIntermediateMultiplePrescriptions(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New().String()
	squatLiftID := uuid.New().String()
	benchLiftID := uuid.New().String()

	// Setup maxes
	maxLookup := newRTSMaxLookup()
	maxLookup.SetMax(userID, squatLiftID, "ONE_RM", 365.0, "2024-01-15")
	maxLookup.SetMax(userID, benchLiftID, "ONE_RM", 300.0, "2024-01-15")

	chart := rpechart.NewDefaultRPEChart()

	// Test cases for different lift/RPE combinations
	testCases := []struct {
		name           string
		liftID         string
		targetReps     int
		targetRPE      float64
		oneRM          float64
		expectedWeight float64
		description    string
	}{
		{
			name:           "Squat 4x@9",
			liftID:         squatLiftID,
			targetReps:     4,
			targetRPE:      9.0,
			oneRM:          365.0,
			expectedWeight: 300.0, // 365 × 0.82 = 299.3 → 300 rounded
			description:    "4 reps @ RPE 9 = 82%",
		},
		{
			name:           "Bench 5x@8",
			liftID:         benchLiftID,
			targetReps:     5,
			targetRPE:      8.0,
			oneRM:          300.0,
			expectedWeight: 230.0, // 300 × 0.77 = 231 → 230 rounded
			description:    "5 reps @ RPE 8 = 77%",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			strategy := loadstrategy.NewRPETargetLoadStrategy(
				tc.targetReps,
				tc.targetRPE,
				5.0,
				loadstrategy.RoundNearest,
				maxLookup,
				chart,
			)

			params := loadstrategy.LoadCalculationParams{
				UserID: userID,
				LiftID: tc.liftID,
			}

			weight, err := strategy.CalculateLoad(ctx, params)
			if err != nil {
				t.Fatalf("CalculateLoad failed: %v", err)
			}

			if math.Abs(weight-tc.expectedWeight) > 0.0001 {
				t.Errorf("%s: expected %.1f lbs, got %.1f lbs", tc.description, tc.expectedWeight, weight)
			}
		})
	}
}

// testRTSIntermediateRounding tests RPE calculations with different rounding directions.
func testRTSIntermediateRounding(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New().String()
	liftID := uuid.New().String()

	// 1RM = 400 lbs, 5 reps @ RPE 8 = 77% → 400 × 0.77 = 308 lbs
	maxLookup := newRTSMaxLookup()
	maxLookup.SetMax(userID, liftID, "ONE_RM", 400.0, "2024-01-15")

	chart := rpechart.NewDefaultRPEChart()

	testCases := []struct {
		name           string
		direction      loadstrategy.RoundingDirection
		increment      float64
		expectedWeight float64
	}{
		{
			name:           "NEAREST rounds 308 to 310",
			direction:      loadstrategy.RoundNearest,
			increment:      5.0,
			expectedWeight: 310.0, // 308 → 310 (nearest)
		},
		{
			name:           "DOWN rounds 308 to 305",
			direction:      loadstrategy.RoundDown,
			increment:      5.0,
			expectedWeight: 305.0, // 308 → 305 (floor)
		},
		{
			name:           "UP rounds 308 to 310",
			direction:      loadstrategy.RoundUp,
			increment:      5.0,
			expectedWeight: 310.0, // 308 → 310 (ceiling)
		},
		{
			name:           "NEAREST with 2.5 increment rounds 308 to 307.5",
			direction:      loadstrategy.RoundNearest,
			increment:      2.5,
			expectedWeight: 307.5, // 308 → 307.5 (nearest 2.5)
		},
		{
			name:           "DOWN with 2.5 increment rounds 308 to 307.5",
			direction:      loadstrategy.RoundDown,
			increment:      2.5,
			expectedWeight: 307.5, // 308 → 307.5 (floor to 2.5)
		},
		{
			name:           "UP with 2.5 increment rounds 308 to 310",
			direction:      loadstrategy.RoundUp,
			increment:      2.5,
			expectedWeight: 310.0, // 308 → 310 (ceiling to 2.5)
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			strategy := loadstrategy.NewRPETargetLoadStrategy(
				5,           // target reps
				8.0,         // target RPE (77%)
				tc.increment,
				tc.direction,
				maxLookup,
				chart,
			)

			params := loadstrategy.LoadCalculationParams{
				UserID: userID,
				LiftID: liftID,
			}

			weight, err := strategy.CalculateLoad(ctx, params)
			if err != nil {
				t.Fatalf("CalculateLoad failed: %v", err)
			}

			if math.Abs(weight-tc.expectedWeight) > 0.0001 {
				t.Errorf("expected %.1f lbs, got %.1f lbs", tc.expectedWeight, weight)
			}
		})
	}
}

// testRTSIntermediateLogSetWithRPE tests logging sets with actual RPE values.
// This data will be used in Phase 7 for e1RM calculations.
func testRTSIntermediateLogSetWithRPE(t *testing.T) {
	userID := uuid.New().String()
	sessionID := uuid.New().String()
	prescriptionID := uuid.New().String()
	liftID := uuid.New().String()

	testCases := []struct {
		name          string
		targetRPE     float64
		actualRPE     float64
		weight        float64
		targetReps    int
		repsPerformed int
		description   string
	}{
		{
			name:          "Set completed at target RPE",
			targetRPE:     8.0,
			actualRPE:     8.0,
			weight:        290.0,
			targetReps:    5,
			repsPerformed: 5,
			description:   "5 reps @ target RPE 8, felt exactly RPE 8",
		},
		{
			name:          "Set harder than expected",
			targetRPE:     8.0,
			actualRPE:     9.0,
			weight:        290.0,
			targetReps:    5,
			repsPerformed: 5,
			description:   "5 reps @ target RPE 8, felt RPE 9 (harder)",
		},
		{
			name:          "Set easier than expected",
			targetRPE:     8.0,
			actualRPE:     7.0,
			weight:        290.0,
			targetReps:    5,
			repsPerformed: 5,
			description:   "5 reps @ target RPE 8, felt RPE 7 (easier)",
		},
		{
			name:          "Extra reps performed",
			targetRPE:     8.0,
			actualRPE:     9.0,
			weight:        290.0,
			targetReps:    5,
			repsPerformed: 7,
			description:   "7 reps @ target RPE 8, felt RPE 9 after extra reps",
		},
		{
			name:          "Half RPE value",
			targetRPE:     8.5,
			actualRPE:     8.5,
			weight:        295.0,
			targetReps:    5,
			repsPerformed: 5,
			description:   "5 reps @ RPE 8.5 (0.5 increments supported)",
		},
		{
			name:          "Max effort set",
			targetRPE:     10.0,
			actualRPE:     10.0,
			weight:        350.0,
			targetReps:    1,
			repsPerformed: 1,
			description:   "1 rep @ RPE 10 (true max effort)",
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rpe := tc.actualRPE

			input := loggedset.CreateLoggedSetInput{
				UserID:         userID,
				SessionID:      sessionID,
				PrescriptionID: prescriptionID,
				LiftID:         liftID,
				SetNumber:      i + 1,
				Weight:         tc.weight,
				TargetReps:     tc.targetReps,
				RepsPerformed:  tc.repsPerformed,
				IsAMRAP:        false,
				RPE:            &rpe,
			}

			logged, valResult := loggedset.NewLoggedSet(input, uuid.New().String())
			if !valResult.Valid {
				t.Fatalf("failed to create logged set: %v", valResult.Errors)
			}

			// Verify RPE is stored correctly
			if logged.RPE == nil {
				t.Error("expected RPE to be stored, got nil")
			} else if *logged.RPE != tc.actualRPE {
				t.Errorf("expected RPE %.1f, got %.1f", tc.actualRPE, *logged.RPE)
			}

			// Verify other fields
			if logged.Weight != tc.weight {
				t.Errorf("expected weight %.1f, got %.1f", tc.weight, logged.Weight)
			}
			if logged.TargetReps != tc.targetReps {
				t.Errorf("expected target reps %d, got %d", tc.targetReps, logged.TargetReps)
			}
			if logged.RepsPerformed != tc.repsPerformed {
				t.Errorf("expected reps performed %d, got %d", tc.repsPerformed, logged.RepsPerformed)
			}

			// Test performance tracking methods
			if tc.repsPerformed > tc.targetReps {
				if !logged.ExceededTarget() {
					t.Error("expected ExceededTarget() to return true")
				}
				if logged.RepsDifference() <= 0 {
					t.Errorf("expected positive RepsDifference, got %d", logged.RepsDifference())
				}
			}
		})
	}
}

// testRTSIntermediateCustomRPEChart tests using a custom RPE chart.
// Some users may have custom RPE charts based on their training history.
func testRTSIntermediateCustomRPEChart(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New().String()
	liftID := uuid.New().String()

	maxLookup := newRTSMaxLookup()
	maxLookup.SetMax(userID, liftID, "ONE_RM", 400.0, "2024-01-15")

	// Create a custom RPE chart with different percentages
	// This simulates a user who has calibrated their own RPE chart
	customEntries := []rpechart.RPEChartEntry{
		// Custom values: this user can handle heavier weights at RPE 8
		{TargetReps: 5, TargetRPE: 8.0, Percentage: 0.80}, // Custom: 80% instead of 77%
		{TargetReps: 5, TargetRPE: 9.0, Percentage: 0.85}, // Custom: 85% instead of 80%
		// Add required entries for chart validation (1-12 reps, 7-10 RPE)
		{TargetReps: 1, TargetRPE: 10.0, Percentage: 1.00},
		{TargetReps: 1, TargetRPE: 9.5, Percentage: 0.975},
		{TargetReps: 1, TargetRPE: 9.0, Percentage: 0.95},
		{TargetReps: 1, TargetRPE: 8.5, Percentage: 0.93},
		{TargetReps: 1, TargetRPE: 8.0, Percentage: 0.91},
		{TargetReps: 1, TargetRPE: 7.5, Percentage: 0.895},
		{TargetReps: 1, TargetRPE: 7.0, Percentage: 0.88},
	}

	customChart, err := rpechart.NewRPEChart(customEntries)
	if err != nil {
		t.Fatalf("failed to create custom RPE chart: %v", err)
	}

	// Test with custom chart: 5 reps @ RPE 8 = 80% (custom) instead of 77% (default)
	t.Run("Custom chart 5x@8 gives 80%", func(t *testing.T) {
		strategy := loadstrategy.NewRPETargetLoadStrategy(
			5,
			8.0,
			5.0,
			loadstrategy.RoundNearest,
			maxLookup,
			customChart,
		)

		params := loadstrategy.LoadCalculationParams{
			UserID: userID,
			LiftID: liftID,
		}

		weight, err := strategy.CalculateLoad(ctx, params)
		if err != nil {
			t.Fatalf("CalculateLoad failed: %v", err)
		}

		// Custom: 400 × 0.80 = 320 lbs (instead of 400 × 0.77 = 308 → 310 with default)
		expectedWeight := 320.0
		if math.Abs(weight-expectedWeight) > 0.0001 {
			t.Errorf("expected %.1f lbs with custom chart, got %.1f lbs", expectedWeight, weight)
		}
	})

	// Test with custom chart: 5 reps @ RPE 9 = 85% (custom) instead of 80% (default)
	t.Run("Custom chart 5x@9 gives 85%", func(t *testing.T) {
		strategy := loadstrategy.NewRPETargetLoadStrategy(
			5,
			9.0,
			5.0,
			loadstrategy.RoundNearest,
			maxLookup,
			customChart,
		)

		params := loadstrategy.LoadCalculationParams{
			UserID: userID,
			LiftID: liftID,
		}

		weight, err := strategy.CalculateLoad(ctx, params)
		if err != nil {
			t.Fatalf("CalculateLoad failed: %v", err)
		}

		// Custom: 400 × 0.85 = 340 lbs (instead of 400 × 0.80 = 320 with default)
		expectedWeight := 340.0
		if math.Abs(weight-expectedWeight) > 0.0001 {
			t.Errorf("expected %.1f lbs with custom chart, got %.1f lbs", expectedWeight, weight)
		}
	})

	// Test using LookupContext to pass custom chart
	t.Run("Custom chart via LookupContext", func(t *testing.T) {
		// Create strategy without injected chart
		strategy := loadstrategy.NewRPETargetLoadStrategy(
			5,
			8.0,
			5.0,
			loadstrategy.RoundNearest,
			maxLookup,
			nil, // No injected chart
		)

		// Pass custom chart via LookupContext
		params := loadstrategy.LoadCalculationParams{
			UserID: userID,
			LiftID: liftID,
			LookupContext: &loadstrategy.LookupContext{
				RPEChart: customChart,
			},
		}

		weight, err := strategy.CalculateLoad(ctx, params)
		if err != nil {
			t.Fatalf("CalculateLoad failed: %v", err)
		}

		expectedWeight := 320.0 // Custom 80%
		if math.Abs(weight-expectedWeight) > 0.0001 {
			t.Errorf("expected %.1f lbs via LookupContext, got %.1f lbs", expectedWeight, weight)
		}
	})
}

// =============================================================================
// RTS INTERMEDIATE ACCEPTANCE TESTS
// =============================================================================

// TestRTSIntermediateAcceptanceCriteria validates RTS Intermediate business requirements.
func TestRTSIntermediateAcceptanceCriteria(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}

	// AC1: RPE chart lookup returns correct percentages
	t.Run("AC1: RPE chart lookup returns correct percentages", func(t *testing.T) {
		chart := rpechart.NewDefaultRPEChart()

		testCases := []struct {
			reps       int
			rpe        float64
			expected   float64
		}{
			{1, 10.0, 1.00},   // 1 rep @ RPE 10 = 100%
			{5, 9.0, 0.80},    // 5 reps @ RPE 9 = 80%
			{5, 8.0, 0.77},    // 5 reps @ RPE 8 = 77%
			{3, 9.0, 0.89},    // 3 reps @ RPE 9 = 89%
			{4, 9.0, 0.82},    // 4 reps @ RPE 9 = 82%
			{8, 8.0, 0.66},    // 8 reps @ RPE 8 = 66%
			{10, 7.0, 0.60},   // 10 reps @ RPE 7 = 60%
		}

		for _, tc := range testCases {
			percentage, err := chart.GetPercentage(tc.reps, tc.rpe)
			if err != nil {
				t.Errorf("%d reps @ RPE %.1f: lookup failed: %v", tc.reps, tc.rpe, err)
				continue
			}
			if math.Abs(percentage-tc.expected) > 0.0001 {
				t.Errorf("%d reps @ RPE %.1f: expected %.2f, got %.2f", tc.reps, tc.rpe, tc.expected, percentage)
			}
		}
	})

	// AC2: RPETarget LoadStrategy calculates correct weights
	t.Run("AC2: RPETarget LoadStrategy calculates correct weights", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New().String()
		liftID := uuid.New().String()

		maxLookup := newRTSMaxLookup()
		maxLookup.SetMax(userID, liftID, "ONE_RM", 400.0, "2024-01-15")

		chart := rpechart.NewDefaultRPEChart()

		testCases := []struct {
			reps           int
			rpe            float64
			expectedWeight float64
		}{
			{1, 10.0, 400.0},  // 400 × 1.00 = 400
			{5, 9.0, 320.0},   // 400 × 0.80 = 320
			{5, 8.0, 310.0},   // 400 × 0.77 = 308 → 310
			{3, 9.0, 355.0},   // 400 × 0.89 = 356 → 355
		}

		for _, tc := range testCases {
			strategy := loadstrategy.NewRPETargetLoadStrategy(
				tc.reps,
				tc.rpe,
				5.0,
				loadstrategy.RoundNearest,
				maxLookup,
				chart,
			)

			params := loadstrategy.LoadCalculationParams{
				UserID: userID,
				LiftID: liftID,
			}

			weight, err := strategy.CalculateLoad(ctx, params)
			if err != nil {
				t.Fatalf("CalculateLoad failed: %v", err)
			}

			if math.Abs(weight-tc.expectedWeight) > 0.0001 {
				t.Errorf("%d reps @ RPE %.1f: expected %.1f, got %.1f", tc.reps, tc.rpe, tc.expectedWeight, weight)
			}
		}
	})

	// AC3: LoggedSet stores RPE values correctly
	t.Run("AC3: LoggedSet stores RPE values correctly", func(t *testing.T) {
		validRPEValues := []float64{5.0, 5.5, 6.0, 6.5, 7.0, 7.5, 8.0, 8.5, 9.0, 9.5, 10.0}

		for _, rpe := range validRPEValues {
			input := loggedset.CreateLoggedSetInput{
				UserID:         uuid.New().String(),
				SessionID:      uuid.New().String(),
				PrescriptionID: uuid.New().String(),
				LiftID:         uuid.New().String(),
				SetNumber:      1,
				Weight:         200.0,
				TargetReps:     5,
				RepsPerformed:  5,
				IsAMRAP:        false,
				RPE:            &rpe,
			}

			logged, valResult := loggedset.NewLoggedSet(input, uuid.New().String())
			if !valResult.Valid {
				t.Errorf("RPE %.1f: validation failed: %v", rpe, valResult.Errors)
				continue
			}

			if logged.RPE == nil || *logged.RPE != rpe {
				t.Errorf("RPE %.1f: expected RPE to be stored correctly", rpe)
			}
		}
	})

	// AC4: LoggedSet rejects invalid RPE values
	t.Run("AC4: LoggedSet rejects invalid RPE values", func(t *testing.T) {
		invalidRPEValues := []float64{4.5, 4.9, 10.5, 11.0}

		for _, rpe := range invalidRPEValues {
			input := loggedset.CreateLoggedSetInput{
				UserID:         uuid.New().String(),
				SessionID:      uuid.New().String(),
				PrescriptionID: uuid.New().String(),
				LiftID:         uuid.New().String(),
				SetNumber:      1,
				Weight:         200.0,
				TargetReps:     5,
				RepsPerformed:  5,
				IsAMRAP:        false,
				RPE:            &rpe,
			}

			_, valResult := loggedset.NewLoggedSet(input, uuid.New().String())
			if valResult.Valid {
				t.Errorf("RPE %.1f: expected validation to fail for invalid RPE", rpe)
			}
		}
	})

	// AC5: RPE chart validates entry parameters
	t.Run("AC5: RPE chart validates entry parameters", func(t *testing.T) {
		invalidEntries := []struct {
			name  string
			entry rpechart.RPEChartEntry
		}{
			{"reps 0", rpechart.RPEChartEntry{TargetReps: 0, TargetRPE: 8.0, Percentage: 0.77}},
			{"reps 13", rpechart.RPEChartEntry{TargetReps: 13, TargetRPE: 8.0, Percentage: 0.77}},
			{"RPE 6.5", rpechart.RPEChartEntry{TargetReps: 5, TargetRPE: 6.5, Percentage: 0.77}},
			{"RPE 10.5", rpechart.RPEChartEntry{TargetReps: 5, TargetRPE: 10.5, Percentage: 0.77}},
			{"percentage -0.1", rpechart.RPEChartEntry{TargetReps: 5, TargetRPE: 8.0, Percentage: -0.1}},
			{"percentage 1.1", rpechart.RPEChartEntry{TargetReps: 5, TargetRPE: 8.0, Percentage: 1.1}},
		}

		for _, tc := range invalidEntries {
			err := rpechart.ValidateEntry(tc.entry)
			if err == nil {
				t.Errorf("%s: expected validation error, got nil", tc.name)
			}
		}
	})

	// AC6: Different lifts use their own 1RMs
	t.Run("AC6: Different lifts use their own 1RMs", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New().String()
		squatID := uuid.New().String()
		benchID := uuid.New().String()

		maxLookup := newRTSMaxLookup()
		maxLookup.SetMax(userID, squatID, "ONE_RM", 400.0, "2024-01-15")
		maxLookup.SetMax(userID, benchID, "ONE_RM", 300.0, "2024-01-15")

		chart := rpechart.NewDefaultRPEChart()

		// Same RPE prescription, different 1RMs
		strategy := loadstrategy.NewRPETargetLoadStrategy(
			5,
			8.0,
			5.0,
			loadstrategy.RoundNearest,
			maxLookup,
			chart,
		)

		// Squat: 400 × 0.77 = 308 → 310
		squatParams := loadstrategy.LoadCalculationParams{UserID: userID, LiftID: squatID}
		squatWeight, _ := strategy.CalculateLoad(ctx, squatParams)
		if squatWeight != 310.0 {
			t.Errorf("squat: expected 310.0, got %.1f", squatWeight)
		}

		// Bench: 300 × 0.77 = 231 → 230
		benchParams := loadstrategy.LoadCalculationParams{UserID: userID, LiftID: benchID}
		benchWeight, _ := strategy.CalculateLoad(ctx, benchParams)
		if benchWeight != 230.0 {
			t.Errorf("bench: expected 230.0, got %.1f", benchWeight)
		}
	})
}
