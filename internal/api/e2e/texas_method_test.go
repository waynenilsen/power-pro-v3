// Package e2e provides end-to-end tests for complete program workflows.
// These tests validate entire program configurations from setup through execution.
package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// =============================================================================
// TEXAS METHOD E2E TEST
// =============================================================================

// TestTexasMethodProgram validates the complete Texas Method program
// configuration and execution through the API.
//
// Texas Method characteristics:
// - 3-day weekly cycle: Monday (Volume), Wednesday (Recovery), Friday (Intensity)
// - Daily intensity variation using DailyLookup:
//   - Monday (Volume): 90% of Friday intensity
//   - Wednesday (Recovery): 80% of Monday (~72% of Friday)
//   - Friday (Intensity): 100% - PR attempt day
// - Different set schemes by day:
//   - Monday: 5x5 (volume accumulation)
//   - Wednesday: 2x5 (active recovery)
//   - Friday: 1x5 (intensity/PR attempt)
// - Weekly LinearProgression: +5lb lower body, +2.5lb upper body after successful Friday
func TestTexasMethodProgram(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Test-unique identifiers
	testID := uuid.New().String()[:8]
	// Use a seeded test user (required for foreign key constraints)
	userID := "workout-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"

	// Create Press lift (not seeded)
	pressSlug := "press-tm-" + testID
	pressBody := fmt.Sprintf(`{"name": "Overhead Press", "slug": "%s", "isCompetitionLift": false}`, pressSlug)
	pressResp, err := adminPost(ts.URL("/lifts"), pressBody)
	if err != nil {
		t.Fatalf("Failed to create press lift: %v", err)
	}
	if pressResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(pressResp.Body)
		pressResp.Body.Close()
		t.Fatalf("Failed to create press lift, status %d: %s", pressResp.StatusCode, body)
	}
	var pressEnvelope LiftResponse
	json.NewDecoder(pressResp.Body).Decode(&pressEnvelope)
	pressResp.Body.Close()
	pressID := pressEnvelope.Data.ID

	// Texas Method training maxes (based on Friday Intensity Day weights)
	squatMax := 315.0 // Squat 5RM for Friday
	benchMax := 225.0 // Bench 5RM for Friday
	pressMax := 135.0 // Press 5RM for Friday

	// Create training maxes for the user
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatMax)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchMax)
	createLiftMax(t, ts, userID, pressID, "TRAINING_MAX", pressMax)

	// =============================================================================
	// Create Daily Lookup for Volume/Recovery/Intensity intensities
	// - volume: 90% (Monday - high volume at moderate intensity)
	// - recovery: 72% (Wednesday - 80% of Monday = 80% × 90% = 72%)
	// - intensity: 100% (Friday - PR attempts)
	// =============================================================================
	dailyLookupBody := `{
		"name": "Texas Method V/R/I Intensity",
		"entries": [
			{"dayIdentifier": "volume", "percentageModifier": 90.0, "intensityLevel": "MEDIUM"},
			{"dayIdentifier": "recovery", "percentageModifier": 72.0, "intensityLevel": "LIGHT"},
			{"dayIdentifier": "intensity", "percentageModifier": 100.0, "intensityLevel": "HEAVY"}
		]
	}`
	dailyLookupResp, err := adminPost(ts.URL("/daily-lookups"), dailyLookupBody)
	if err != nil {
		t.Fatalf("Failed to create daily lookup: %v", err)
	}
	if dailyLookupResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(dailyLookupResp.Body)
		dailyLookupResp.Body.Close()
		t.Fatalf("Failed to create daily lookup, status %d: %s", dailyLookupResp.StatusCode, body)
	}
	var dailyLookupEnvelope struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	json.NewDecoder(dailyLookupResp.Body).Decode(&dailyLookupEnvelope)
	dailyLookupResp.Body.Close()
	dailyLookupID := dailyLookupEnvelope.Data.ID

	// =============================================================================
	// Create prescriptions for each day
	// Volume Day (Monday): 5x5 FIXED sets
	// Recovery Day (Wednesday): 2x5 FIXED sets
	// Intensity Day (Friday): 1x5 FIXED sets (PR attempt)
	// =============================================================================

	// Volume Day prescriptions (5x5)
	squatVolumePrescID := createFixedPrescriptionWithLookup(t, ts, squatID, 5, 5, 100.0, 0, "day")
	benchVolumePrescID := createFixedPrescriptionWithLookup(t, ts, benchID, 5, 5, 100.0, 1, "day")

	// Recovery Day prescriptions (2x5)
	squatRecoveryPrescID := createFixedPrescriptionWithLookup(t, ts, squatID, 2, 5, 100.0, 0, "day")
	pressRecoveryPrescID := createFixedPrescriptionWithLookup(t, ts, pressID, 2, 5, 100.0, 1, "day")

	// Intensity Day prescriptions (1x5 - PR attempt)
	squatIntensityPrescID := createFixedPrescriptionWithLookup(t, ts, squatID, 1, 5, 100.0, 0, "day")
	benchIntensityPrescID := createFixedPrescriptionWithLookup(t, ts, benchID, 1, 5, 100.0, 1, "day")

	// =============================================================================
	// Create Days: Volume (Monday), Recovery (Wednesday), Intensity (Friday)
	// The day slug must match the dailyLookup dayIdentifier for intensity lookup
	// =============================================================================

	// Volume Day (Monday)
	volumeDayBody := `{"name": "Volume Day", "slug": "volume"}`
	volumeDayResp, _ := adminPost(ts.URL("/days"), volumeDayBody)
	var volumeDayEnvelope DayResponse
	json.NewDecoder(volumeDayResp.Body).Decode(&volumeDayEnvelope)
	volumeDayResp.Body.Close()
	volumeDayID := volumeDayEnvelope.Data.ID

	addPrescToDay(t, ts, volumeDayID, squatVolumePrescID)
	addPrescToDay(t, ts, volumeDayID, benchVolumePrescID)

	// Recovery Day (Wednesday)
	recoveryDayBody := `{"name": "Recovery Day", "slug": "recovery"}`
	recoveryDayResp, _ := adminPost(ts.URL("/days"), recoveryDayBody)
	var recoveryDayEnvelope DayResponse
	json.NewDecoder(recoveryDayResp.Body).Decode(&recoveryDayEnvelope)
	recoveryDayResp.Body.Close()
	recoveryDayID := recoveryDayEnvelope.Data.ID

	addPrescToDay(t, ts, recoveryDayID, squatRecoveryPrescID)
	addPrescToDay(t, ts, recoveryDayID, pressRecoveryPrescID)

	// Intensity Day (Friday)
	intensityDayBody := `{"name": "Intensity Day", "slug": "intensity"}`
	intensityDayResp, _ := adminPost(ts.URL("/days"), intensityDayBody)
	var intensityDayEnvelope DayResponse
	json.NewDecoder(intensityDayResp.Body).Decode(&intensityDayEnvelope)
	intensityDayResp.Body.Close()
	intensityDayID := intensityDayEnvelope.Data.ID

	addPrescToDay(t, ts, intensityDayID, squatIntensityPrescID)
	addPrescToDay(t, ts, intensityDayID, benchIntensityPrescID)

	// =============================================================================
	// Create 1-week cycle with Volume/Recovery/Intensity pattern (Mon/Wed/Fri)
	// =============================================================================
	cycleName := "Texas Method Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 1}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	// Create week 1 in the cycle
	weekBody := fmt.Sprintf(`{"weekNumber": 1, "cycleId": "%s"}`, cycleID)
	weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
	var weekEnvelope WeekResponse
	json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
	weekResp.Body.Close()
	weekID := weekEnvelope.Data.ID

	// Add days to week: Volume/Recovery/Intensity pattern
	addDayToWeek(t, ts, weekID, volumeDayID, "MONDAY")
	addDayToWeek(t, ts, weekID, recoveryDayID, "WEDNESDAY")
	addDayToWeek(t, ts, weekID, intensityDayID, "FRIDAY")

	// =============================================================================
	// Create program and link daily lookup
	// =============================================================================
	programSlug := "texas-method-" + testID
	programBody := fmt.Sprintf(`{"name": "Texas Method", "slug": "%s", "cycleId": "%s", "dailyLookupId": "%s"}`,
		programSlug, cycleID, dailyLookupID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// Create Linear Progressions (AFTER_WEEK trigger)
	// Lower body: +5lb, Upper body: +2.5lb
	// =============================================================================
	lowerProgBody := `{"name": "TM Lower Linear", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_WEEK"}}`
	lowerProgResp, _ := adminPost(ts.URL("/progressions"), lowerProgBody)
	var lowerProgEnvelope ProgressionResponse
	json.NewDecoder(lowerProgResp.Body).Decode(&lowerProgEnvelope)
	lowerProgResp.Body.Close()
	lowerProgID := lowerProgEnvelope.Data.ID

	upperProgBody := `{"name": "TM Upper Linear", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 2.5, "maxType": "TRAINING_MAX", "triggerType": "AFTER_WEEK"}}`
	upperProgResp, _ := adminPost(ts.URL("/progressions"), upperProgBody)
	var upperProgEnvelope ProgressionResponse
	json.NewDecoder(upperProgResp.Body).Decode(&upperProgEnvelope)
	upperProgResp.Body.Close()
	upperProgID := upperProgEnvelope.Data.ID

	// Link progressions to program
	linkProgressionToProgram(t, ts, programID, lowerProgID, squatID, 1)
	linkProgressionToProgram(t, ts, programID, upperProgID, benchID, 2)
	linkProgressionToProgram(t, ts, programID, upperProgID, pressID, 3)

	// =============================================================================
	// Enroll user in program
	// =============================================================================
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, err := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	if err != nil {
		t.Fatalf("Failed to enroll user: %v", err)
	}
	if enrollResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(enrollResp.Body)
		enrollResp.Body.Close()
		t.Fatalf("Failed to enroll user, status %d: %s", enrollResp.StatusCode, body)
	}
	enrollResp.Body.Close()

	// =============================================================================
	// EXECUTION PHASE: Volume Day (Monday - Workout 1)
	// =============================================================================
	t.Run("Volume Day generates 5x5 sets at 90% intensity", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		if workoutResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(workoutResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", workoutResp.StatusCode, body)
		}

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Verify Volume Day
		if workout.Data.DaySlug != "volume" {
			t.Errorf("Expected day slug 'volume', got '%s'", workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) != 2 {
			t.Fatalf("Expected 2 exercises on Volume Day, got %d", len(workout.Data.Exercises))
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Squat on Volume Day: 5x5 at 90% daily intensity
		// Expected weight: 315 * 90% = 283.5, rounded to nearest 5 = 285
		if squat, ok := exercisesByLift[squatID]; ok {
			if len(squat.Sets) != 5 {
				t.Errorf("Squat: expected 5 sets, got %d", len(squat.Sets))
			}
			expectedWeight := 285.0 // 315 * 0.90 = 283.5 → rounded to 285
			for i, set := range squat.Sets {
				if !closeEnough(set.Weight, expectedWeight) {
					t.Errorf("Squat set %d: expected weight ~%.1f, got %.1f", i+1, expectedWeight, set.Weight)
				}
				if set.TargetReps != 5 {
					t.Errorf("Squat set %d: expected 5 reps, got %d", i+1, set.TargetReps)
				}
			}
		} else {
			t.Error("Volume Day missing Squat exercise")
		}

		// Bench on Volume Day: 5x5 at 90% daily intensity
		// Expected weight: 225 * 90% = 202.5, rounded to nearest 5 = 205
		if bench, ok := exercisesByLift[benchID]; ok {
			if len(bench.Sets) != 5 {
				t.Errorf("Bench: expected 5 sets, got %d", len(bench.Sets))
			}
			expectedWeight := 205.0 // 225 * 0.90 = 202.5 → rounded to 205
			for i, set := range bench.Sets {
				if !closeEnough(set.Weight, expectedWeight) {
					t.Errorf("Bench set %d: expected weight ~%.1f, got %.1f", i+1, expectedWeight, set.Weight)
				}
				if set.TargetReps != 5 {
					t.Errorf("Bench set %d: expected 5 reps, got %d", i+1, set.TargetReps)
				}
			}
		} else {
			t.Error("Volume Day missing Bench exercise")
		}
	})

	// Complete Volume Day workout using explicit state machine flow
	completeTMWorkoutDay(t, ts, userID)

	// =============================================================================
	// EXECUTION PHASE: Recovery Day (Wednesday - Workout 2)
	// =============================================================================
	t.Run("Recovery Day generates 2x5 sets at 72% intensity", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		if workoutResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(workoutResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", workoutResp.StatusCode, body)
		}

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Verify Recovery Day
		if workout.Data.DaySlug != "recovery" {
			t.Errorf("Expected day slug 'recovery', got '%s'", workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) != 2 {
			t.Fatalf("Expected 2 exercises on Recovery Day, got %d", len(workout.Data.Exercises))
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Squat on Recovery Day: 2x5 at 72% daily intensity
		// Expected weight: 315 * 72% = 226.8, rounded to nearest 5 = 225
		if squat, ok := exercisesByLift[squatID]; ok {
			if len(squat.Sets) != 2 {
				t.Errorf("Recovery Squat: expected 2 sets, got %d", len(squat.Sets))
			}
			expectedWeight := 225.0 // 315 * 0.72 = 226.8 → rounded to 225
			for i, set := range squat.Sets {
				if !closeEnough(set.Weight, expectedWeight) {
					t.Errorf("Recovery Squat set %d: expected weight ~%.1f, got %.1f", i+1, expectedWeight, set.Weight)
				}
				if set.TargetReps != 5 {
					t.Errorf("Recovery Squat set %d: expected 5 reps, got %d", i+1, set.TargetReps)
				}
			}
		} else {
			t.Error("Recovery Day missing Squat exercise")
		}

		// Press on Recovery Day: 2x5 at 72% daily intensity
		// Expected weight: 135 * 72% = 97.2, rounded to nearest 5 = 95
		if press, ok := exercisesByLift[pressID]; ok {
			if len(press.Sets) != 2 {
				t.Errorf("Recovery Press: expected 2 sets, got %d", len(press.Sets))
			}
			expectedWeight := 95.0 // 135 * 0.72 = 97.2 → rounded to 95
			for i, set := range press.Sets {
				if !closeEnough(set.Weight, expectedWeight) {
					t.Errorf("Recovery Press set %d: expected weight ~%.1f, got %.1f", i+1, expectedWeight, set.Weight)
				}
				if set.TargetReps != 5 {
					t.Errorf("Recovery Press set %d: expected 5 reps, got %d", i+1, set.TargetReps)
				}
			}
		} else {
			t.Error("Recovery Day missing Press exercise")
		}
	})

	// Complete Recovery Day workout using explicit state machine flow
	completeTMWorkoutDay(t, ts, userID)

	// =============================================================================
	// EXECUTION PHASE: Intensity Day (Friday - Workout 3)
	// =============================================================================
	t.Run("Intensity Day generates 1x5 sets at 100% intensity (PR attempt)", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		if workoutResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(workoutResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", workoutResp.StatusCode, body)
		}

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Verify Intensity Day
		if workout.Data.DaySlug != "intensity" {
			t.Errorf("Expected day slug 'intensity', got '%s'", workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) != 2 {
			t.Fatalf("Expected 2 exercises on Intensity Day, got %d", len(workout.Data.Exercises))
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Squat on Intensity Day: 1x5 at 100% daily intensity (PR attempt)
		// Expected weight: 315 * 100% = 315
		if squat, ok := exercisesByLift[squatID]; ok {
			if len(squat.Sets) != 1 {
				t.Errorf("Intensity Squat: expected 1 set, got %d", len(squat.Sets))
			}
			expectedWeight := squatMax * 1.00 // 315
			for i, set := range squat.Sets {
				if !closeEnough(set.Weight, expectedWeight) {
					t.Errorf("Intensity Squat set %d: expected weight ~%.1f, got %.1f", i+1, expectedWeight, set.Weight)
				}
				if set.TargetReps != 5 {
					t.Errorf("Intensity Squat set %d: expected 5 reps, got %d", i+1, set.TargetReps)
				}
			}
		} else {
			t.Error("Intensity Day missing Squat exercise")
		}

		// Bench on Intensity Day: 1x5 at 100% daily intensity (PR attempt)
		// Expected weight: 225 * 100% = 225
		if bench, ok := exercisesByLift[benchID]; ok {
			if len(bench.Sets) != 1 {
				t.Errorf("Intensity Bench: expected 1 set, got %d", len(bench.Sets))
			}
			expectedWeight := benchMax * 1.00 // 225
			for i, set := range bench.Sets {
				if !closeEnough(set.Weight, expectedWeight) {
					t.Errorf("Intensity Bench set %d: expected weight ~%.1f, got %.1f", i+1, expectedWeight, set.Weight)
				}
				if set.TargetReps != 5 {
					t.Errorf("Intensity Bench set %d: expected 5 reps, got %d", i+1, set.TargetReps)
				}
			}
		} else {
			t.Error("Intensity Day missing Bench exercise")
		}
	})

	// =============================================================================
	// PROGRESSION PHASE: Trigger weekly progression after completing Intensity Day
	// =============================================================================
	t.Run("Weekly progression applies correct increments (+5lb lower, +2.5lb upper)", func(t *testing.T) {
		// Trigger lower body progression for squat (+5lb)
		squatTrigger := triggerProgressionForLift(t, ts, userID, lowerProgID, squatID)

		if squatTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected squat progression to apply, got TotalApplied=%d", squatTrigger.Data.TotalApplied)
		}
		if len(squatTrigger.Data.Results) > 0 && squatTrigger.Data.Results[0].Result != nil {
			if squatTrigger.Data.Results[0].Result.Delta != 5.0 {
				t.Errorf("Expected squat delta +5, got %f", squatTrigger.Data.Results[0].Result.Delta)
			}
			expectedNewSquat := squatMax + 5.0 // 320
			if squatTrigger.Data.Results[0].Result.NewValue != expectedNewSquat {
				t.Errorf("Expected squat new value %f, got %f", expectedNewSquat, squatTrigger.Data.Results[0].Result.NewValue)
			}
		}

		// Trigger upper body progression for bench (+2.5lb)
		benchTrigger := triggerProgressionForLift(t, ts, userID, upperProgID, benchID)

		if benchTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected bench progression to apply")
		}
		if len(benchTrigger.Data.Results) > 0 && benchTrigger.Data.Results[0].Result != nil {
			if benchTrigger.Data.Results[0].Result.Delta != 2.5 {
				t.Errorf("Expected bench delta +2.5, got %f", benchTrigger.Data.Results[0].Result.Delta)
			}
			expectedNewBench := benchMax + 2.5 // 227.5
			if benchTrigger.Data.Results[0].Result.NewValue != expectedNewBench {
				t.Errorf("Expected bench new value %f, got %f", expectedNewBench, benchTrigger.Data.Results[0].Result.NewValue)
			}
		}

		// Trigger upper body progression for press (+2.5lb)
		pressTrigger := triggerProgressionForLift(t, ts, userID, upperProgID, pressID)

		if pressTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected press progression to apply")
		}
		if len(pressTrigger.Data.Results) > 0 && pressTrigger.Data.Results[0].Result != nil {
			if pressTrigger.Data.Results[0].Result.Delta != 2.5 {
				t.Errorf("Expected press delta +2.5, got %f", pressTrigger.Data.Results[0].Result.Delta)
			}
			expectedNewPress := pressMax + 2.5 // 137.5
			if pressTrigger.Data.Results[0].Result.NewValue != expectedNewPress {
				t.Errorf("Expected press new value %f, got %f", expectedNewPress, pressTrigger.Data.Results[0].Result.NewValue)
			}
		}
	})

	// Complete Intensity Day workout using explicit state machine flow
	completeTMWorkoutDay(t, ts, userID)

	// =============================================================================
	// VALIDATION PHASE: Next week's Volume Day should show increased weights
	// =============================================================================
	t.Run("Next week Volume Day shows progression on all lifts", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		if workoutResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(workoutResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", workoutResp.StatusCode, body)
		}

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Should be Volume Day again
		if workout.Data.DaySlug != "volume" {
			t.Errorf("Expected day slug 'volume', got '%s'", workout.Data.DaySlug)
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Squat should now be at (315 + 5) * 90% = 320 * 0.9 = 288, rounded to 290
		expectedSquatWeight := 290.0 // 320 * 0.90 = 288 → rounded to 290
		if squat, ok := exercisesByLift[squatID]; ok {
			for i, set := range squat.Sets {
				if !closeEnough(set.Weight, expectedSquatWeight) {
					t.Errorf("Week 2 Volume Squat set %d: expected weight ~%.1f, got %.1f", i+1, expectedSquatWeight, set.Weight)
				}
			}
		}

		// Bench should now be at (225 + 2.5) * 90% = 227.5 * 0.9 = 204.75, rounded to 205
		expectedBenchWeight := 205.0 // 227.5 * 0.90 = 204.75 → rounded to 205
		if bench, ok := exercisesByLift[benchID]; ok {
			for i, set := range bench.Sets {
				if !closeEnough(set.Weight, expectedBenchWeight) {
					t.Errorf("Week 2 Volume Bench set %d: expected weight ~%.1f, got %.1f", i+1, expectedBenchWeight, set.Weight)
				}
			}
		}
	})

	// Complete Week 2 Volume Day workout using explicit state machine flow
	completeTMWorkoutDay(t, ts, userID)

	// =============================================================================
	// VALIDATION PHASE: Next week's Recovery Day should show increased weights
	// =============================================================================
	t.Run("Next week Recovery Day shows progression on all lifts", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		if workoutResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(workoutResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", workoutResp.StatusCode, body)
		}

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Should be Recovery Day
		if workout.Data.DaySlug != "recovery" {
			t.Errorf("Expected day slug 'recovery', got '%s'", workout.Data.DaySlug)
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Squat should now be at (315 + 5) * 72% = 320 * 0.72 = 230.4, rounded to 230
		expectedSquatWeight := 230.0 // 320 * 0.72 = 230.4 → rounded to 230
		if squat, ok := exercisesByLift[squatID]; ok {
			for i, set := range squat.Sets {
				if !closeEnough(set.Weight, expectedSquatWeight) {
					t.Errorf("Week 2 Recovery Squat set %d: expected weight ~%.1f, got %.1f", i+1, expectedSquatWeight, set.Weight)
				}
			}
		}

		// Press should now be at (135 + 2.5) * 72% = 137.5 * 0.72 = 99, rounded to 100
		expectedPressWeight := 100.0 // 137.5 * 0.72 = 99 → rounded to 100
		if press, ok := exercisesByLift[pressID]; ok {
			for i, set := range press.Sets {
				if !closeEnough(set.Weight, expectedPressWeight) {
					t.Errorf("Week 2 Recovery Press set %d: expected weight ~%.1f, got %.1f", i+1, expectedPressWeight, set.Weight)
				}
			}
		}
	})

	// Complete Week 2 Recovery Day workout using explicit state machine flow
	completeTMWorkoutDay(t, ts, userID)

	// =============================================================================
	// VALIDATION PHASE: Next week's Intensity Day should show increased weights
	// =============================================================================
	t.Run("Next week Intensity Day shows progression at 100% (new PR attempt)", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		if workoutResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(workoutResp.Body)
			t.Fatalf("Expected status 200, got %d: %s", workoutResp.StatusCode, body)
		}

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Should be Intensity Day
		if workout.Data.DaySlug != "intensity" {
			t.Errorf("Expected day slug 'intensity', got '%s'", workout.Data.DaySlug)
		}

		exercisesByLift := make(map[string]WorkoutExerciseData)
		for _, ex := range workout.Data.Exercises {
			exercisesByLift[ex.Lift.ID] = ex
		}

		// Squat should now be at 320 (315 + 5) at 100%
		expectedSquatWeight := 320.0 // 315 + 5 = 320
		if squat, ok := exercisesByLift[squatID]; ok {
			if len(squat.Sets) != 1 {
				t.Errorf("Intensity Squat: expected 1 set, got %d", len(squat.Sets))
			}
			for i, set := range squat.Sets {
				if !closeEnough(set.Weight, expectedSquatWeight) {
					t.Errorf("Week 2 Intensity Squat set %d: expected weight ~%.1f, got %.1f", i+1, expectedSquatWeight, set.Weight)
				}
			}
		}

		// Bench should now be at 227.5 (225 + 2.5) at 100%, rounded to 230
		expectedBenchWeight := 230.0 // 225 + 2.5 = 227.5 → rounded to 230
		if bench, ok := exercisesByLift[benchID]; ok {
			if len(bench.Sets) != 1 {
				t.Errorf("Intensity Bench: expected 1 set, got %d", len(bench.Sets))
			}
			for i, set := range bench.Sets {
				if !closeEnough(set.Weight, expectedBenchWeight) {
					t.Errorf("Week 2 Intensity Bench set %d: expected weight ~%.1f, got %.1f", i+1, expectedBenchWeight, set.Weight)
				}
			}
		}
	})
}

// =============================================================================
// HELPER FUNCTIONS (specific to this test file)
// =============================================================================

// completeTMWorkoutDay completes a Texas Method workout day using explicit state machine flow.
func completeTMWorkoutDay(t *testing.T, ts *testutil.TestServer, userID string) {
	t.Helper()

	sessionID := startWorkoutSession(t, ts, userID)

	workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
	var workout WorkoutResponse
	json.NewDecoder(workoutResp.Body).Decode(&workout)
	workoutResp.Body.Close()

	for _, ex := range workout.Data.Exercises {
		for _, set := range ex.Sets {
			logTMSet(t, ts, userID, sessionID, ex.PrescriptionID, ex.Lift.ID, set.SetNumber, set.Weight, set.TargetReps, set.TargetReps)
		}
	}

	finishWorkoutSession(t, ts, sessionID, userID)
	advanceUserState(t, ts, userID)
}

// logTMSet logs a single set for Texas Method workout.
func logTMSet(t *testing.T, ts *testutil.TestServer, userID, sessionID, prescriptionID, liftID string, setNumber int, weight float64, targetReps, repsPerformed int) {
	t.Helper()

	loggedSetBody := fmt.Sprintf(`{
		"sets": [{
			"prescriptionId": "%s",
			"liftId": "%s",
			"setNumber": %d,
			"weight": %.1f,
			"targetReps": %d,
			"repsPerformed": %d
		}]
	}`, prescriptionID, liftID, setNumber, weight, targetReps, repsPerformed)

	logResp, err := userPost(ts.URL("/sessions/"+sessionID+"/sets"), loggedSetBody, userID)
	if err != nil {
		t.Fatalf("Failed to log set: %v", err)
	}
	if logResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(logResp.Body)
		logResp.Body.Close()
		t.Fatalf("Failed to log set, status %d: %s", logResp.StatusCode, body)
	}
	logResp.Body.Close()
}

// createFixedPrescriptionWithLookup creates a prescription with FIXED set scheme and daily lookup.
func createFixedPrescriptionWithLookup(t *testing.T, ts *testutil.TestServer, liftID string, sets, reps int, percentage float64, order int, lookupKey string) string {
	t.Helper()

	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": %.1f, "lookupKey": "%s"},
		"setScheme": {"type": "FIXED", "sets": %d, "reps": %d},
		"order": %d
	}`, liftID, percentage, lookupKey, sets, reps, order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create fixed prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create fixed prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}
