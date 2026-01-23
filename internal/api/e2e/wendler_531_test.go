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
// WENDLER 5/3/1 BBB E2E TEST
// =============================================================================

// TestWendler531BBBProgram validates the complete Wendler 5/3/1 Boring But Big
// program configuration and execution through the API.
//
// Wendler 5/3/1 BBB characteristics:
// - 4-Week Cycle: Different rep/percentage schemes each week
//   - Week 1 (5s): 65%/75%/85% at 5/5/5 reps
//   - Week 2 (3s): 70%/80%/90% at 3/3/3 reps
//   - Week 3 (5/3/1): 75%/85%/95% at 5/3/1 reps
//   - Week 4 (Deload): 40%/50%/60% at 5/5/5 reps
// - BBB Accessory: 5x10 at 50% after main work
// - CycleProgression: +10lb lower body, +5lb upper body at cycle end
//
// Implementation: Uses RAMP prescriptions with explicit per-set percentages.
// Each week has dedicated days with week-specific prescriptions to achieve
// the different percentage schemes per week.
func TestWendler531BBBProgram(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Test-unique identifiers
	testID := uuid.New().String()[:8]
	// Use a seeded test user
	userID := "workout-test-user"

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	// Create OHP lift (not seeded)
	ohpSlug := "ohp-" + testID
	ohpBody := fmt.Sprintf(`{"name": "Overhead Press", "slug": "%s", "isCompetitionLift": false}`, ohpSlug)
	ohpResp, err := adminPost(ts.URL("/lifts"), ohpBody)
	if err != nil {
		t.Fatalf("Failed to create OHP lift: %v", err)
	}
	if ohpResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(ohpResp.Body)
		ohpResp.Body.Close()
		t.Fatalf("Failed to create OHP lift, status %d: %s", ohpResp.StatusCode, body)
	}
	var ohpEnvelope LiftResponse
	json.NewDecoder(ohpResp.Body).Decode(&ohpEnvelope)
	ohpResp.Body.Close()
	ohpID := ohpEnvelope.Data.ID

	// Wendler training maxes (typically 90% of true 1RM)
	squatTM := 300.0    // Squat training max
	benchTM := 200.0    // Bench training max
	deadliftTM := 350.0 // Deadlift training max
	ohpTM := 135.0      // OHP training max

	// Create training maxes for the user
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatTM)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchTM)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftTM)
	createLiftMax(t, ts, userID, ohpID, "TRAINING_MAX", ohpTM)

	// =============================================================================
	// Create prescriptions for each week using RAMP scheme
	// Week 1 (5s): 65%/75%/85% x 5 reps each
	// =============================================================================
	week1Steps := []RampStep{
		{Percentage: 65, Reps: 5},
		{Percentage: 75, Reps: 5},
		{Percentage: 85, Reps: 5},
	}
	squatW1PrescID := create531RampPrescription(t, ts, squatID, week1Steps, 0)
	benchW1PrescID := create531RampPrescription(t, ts, benchID, week1Steps, 0)
	deadliftW1PrescID := create531RampPrescription(t, ts, deadliftID, week1Steps, 0)
	ohpW1PrescID := create531RampPrescription(t, ts, ohpID, week1Steps, 0)

	// Week 2 (3s): 70%/80%/90% x 3 reps each
	week2Steps := []RampStep{
		{Percentage: 70, Reps: 3},
		{Percentage: 80, Reps: 3},
		{Percentage: 90, Reps: 3},
	}
	squatW2PrescID := create531RampPrescription(t, ts, squatID, week2Steps, 0)
	benchW2PrescID := create531RampPrescription(t, ts, benchID, week2Steps, 0)
	deadliftW2PrescID := create531RampPrescription(t, ts, deadliftID, week2Steps, 0)
	ohpW2PrescID := create531RampPrescription(t, ts, ohpID, week2Steps, 0)

	// Week 3 (5/3/1): 75%/85%/95% with 5/3/1 reps
	week3Steps := []RampStep{
		{Percentage: 75, Reps: 5},
		{Percentage: 85, Reps: 3},
		{Percentage: 95, Reps: 1},
	}
	squatW3PrescID := create531RampPrescription(t, ts, squatID, week3Steps, 0)
	benchW3PrescID := create531RampPrescription(t, ts, benchID, week3Steps, 0)
	deadliftW3PrescID := create531RampPrescription(t, ts, deadliftID, week3Steps, 0)
	ohpW3PrescID := create531RampPrescription(t, ts, ohpID, week3Steps, 0)

	// Week 4 (Deload): 40%/50%/60% x 5 reps each
	week4Steps := []RampStep{
		{Percentage: 40, Reps: 5},
		{Percentage: 50, Reps: 5},
		{Percentage: 60, Reps: 5},
	}
	squatW4PrescID := create531RampPrescription(t, ts, squatID, week4Steps, 0)
	benchW4PrescID := create531RampPrescription(t, ts, benchID, week4Steps, 0)
	deadliftW4PrescID := create531RampPrescription(t, ts, deadliftID, week4Steps, 0)
	ohpW4PrescID := create531RampPrescription(t, ts, ohpID, week4Steps, 0)

	// =============================================================================
	// Create BBB prescriptions (5x10 at 50% - same for all weeks)
	// =============================================================================
	squatBBBPrescID := createPrescription(t, ts, squatID, 5, 10, 50.0, 1)
	benchBBBPrescID := createPrescription(t, ts, benchID, 5, 10, 50.0, 1)
	deadliftBBBPrescID := createPrescription(t, ts, deadliftID, 5, 10, 50.0, 1)
	ohpBBBPrescID := createPrescription(t, ts, ohpID, 5, 10, 50.0, 1)

	// =============================================================================
	// Create Days for each week (week-specific prescriptions)
	// Week 1 Days
	// =============================================================================
	squatW1DaySlug := "squat-w1-" + testID
	squatW1DayID := createDay(t, ts, "Squat Day W1", squatW1DaySlug)
	addPrescToDay(t, ts, squatW1DayID, squatW1PrescID)
	addPrescToDay(t, ts, squatW1DayID, squatBBBPrescID)

	benchW1DaySlug := "bench-w1-" + testID
	benchW1DayID := createDay(t, ts, "Bench Day W1", benchW1DaySlug)
	addPrescToDay(t, ts, benchW1DayID, benchW1PrescID)
	addPrescToDay(t, ts, benchW1DayID, benchBBBPrescID)

	deadliftW1DayID := createDay(t, ts, "Deadlift Day W1", "deadlift-w1-"+testID)
	addPrescToDay(t, ts, deadliftW1DayID, deadliftW1PrescID)
	addPrescToDay(t, ts, deadliftW1DayID, deadliftBBBPrescID)

	ohpW1DayID := createDay(t, ts, "OHP Day W1", "ohp-w1-"+testID)
	addPrescToDay(t, ts, ohpW1DayID, ohpW1PrescID)
	addPrescToDay(t, ts, ohpW1DayID, ohpBBBPrescID)

	// Week 2 Days
	squatW2DaySlug := "squat-w2-" + testID
	squatW2DayID := createDay(t, ts, "Squat Day W2", squatW2DaySlug)
	addPrescToDay(t, ts, squatW2DayID, squatW2PrescID)
	addPrescToDay(t, ts, squatW2DayID, squatBBBPrescID)

	benchW2DayID := createDay(t, ts, "Bench Day W2", "bench-w2-"+testID)
	addPrescToDay(t, ts, benchW2DayID, benchW2PrescID)
	addPrescToDay(t, ts, benchW2DayID, benchBBBPrescID)

	deadliftW2DayID := createDay(t, ts, "Deadlift Day W2", "deadlift-w2-"+testID)
	addPrescToDay(t, ts, deadliftW2DayID, deadliftW2PrescID)
	addPrescToDay(t, ts, deadliftW2DayID, deadliftBBBPrescID)

	ohpW2DayID := createDay(t, ts, "OHP Day W2", "ohp-w2-"+testID)
	addPrescToDay(t, ts, ohpW2DayID, ohpW2PrescID)
	addPrescToDay(t, ts, ohpW2DayID, ohpBBBPrescID)

	// Week 3 Days
	squatW3DaySlug := "squat-w3-" + testID
	squatW3DayID := createDay(t, ts, "Squat Day W3", squatW3DaySlug)
	addPrescToDay(t, ts, squatW3DayID, squatW3PrescID)
	addPrescToDay(t, ts, squatW3DayID, squatBBBPrescID)

	benchW3DayID := createDay(t, ts, "Bench Day W3", "bench-w3-"+testID)
	addPrescToDay(t, ts, benchW3DayID, benchW3PrescID)
	addPrescToDay(t, ts, benchW3DayID, benchBBBPrescID)

	deadliftW3DayID := createDay(t, ts, "Deadlift Day W3", "deadlift-w3-"+testID)
	addPrescToDay(t, ts, deadliftW3DayID, deadliftW3PrescID)
	addPrescToDay(t, ts, deadliftW3DayID, deadliftBBBPrescID)

	ohpW3DayID := createDay(t, ts, "OHP Day W3", "ohp-w3-"+testID)
	addPrescToDay(t, ts, ohpW3DayID, ohpW3PrescID)
	addPrescToDay(t, ts, ohpW3DayID, ohpBBBPrescID)

	// Week 4 Days (Deload)
	squatW4DaySlug := "squat-w4-" + testID
	squatW4DayID := createDay(t, ts, "Squat Day W4", squatW4DaySlug)
	addPrescToDay(t, ts, squatW4DayID, squatW4PrescID)
	addPrescToDay(t, ts, squatW4DayID, squatBBBPrescID)

	benchW4DayID := createDay(t, ts, "Bench Day W4", "bench-w4-"+testID)
	addPrescToDay(t, ts, benchW4DayID, benchW4PrescID)
	addPrescToDay(t, ts, benchW4DayID, benchBBBPrescID)

	deadliftW4DayID := createDay(t, ts, "Deadlift Day W4", "deadlift-w4-"+testID)
	addPrescToDay(t, ts, deadliftW4DayID, deadliftW4PrescID)
	addPrescToDay(t, ts, deadliftW4DayID, deadliftBBBPrescID)

	ohpW4DayID := createDay(t, ts, "OHP Day W4", "ohp-w4-"+testID)
	addPrescToDay(t, ts, ohpW4DayID, ohpW4PrescID)
	addPrescToDay(t, ts, ohpW4DayID, ohpBBBPrescID)

	// =============================================================================
	// Create 4-week cycle with week-specific days
	// =============================================================================
	cycleName := "Wendler 5/3/1 Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 4}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	// Week 1 with week 1 days
	week1ID := createWeek(t, ts, 1, cycleID)
	addDayToWeek(t, ts, week1ID, squatW1DayID, "MONDAY")
	addDayToWeek(t, ts, week1ID, benchW1DayID, "TUESDAY")
	addDayToWeek(t, ts, week1ID, deadliftW1DayID, "THURSDAY")
	addDayToWeek(t, ts, week1ID, ohpW1DayID, "FRIDAY")

	// Week 2 with week 2 days
	week2ID := createWeek(t, ts, 2, cycleID)
	addDayToWeek(t, ts, week2ID, squatW2DayID, "MONDAY")
	addDayToWeek(t, ts, week2ID, benchW2DayID, "TUESDAY")
	addDayToWeek(t, ts, week2ID, deadliftW2DayID, "THURSDAY")
	addDayToWeek(t, ts, week2ID, ohpW2DayID, "FRIDAY")

	// Week 3 with week 3 days
	week3ID := createWeek(t, ts, 3, cycleID)
	addDayToWeek(t, ts, week3ID, squatW3DayID, "MONDAY")
	addDayToWeek(t, ts, week3ID, benchW3DayID, "TUESDAY")
	addDayToWeek(t, ts, week3ID, deadliftW3DayID, "THURSDAY")
	addDayToWeek(t, ts, week3ID, ohpW3DayID, "FRIDAY")

	// Week 4 (Deload) with week 4 days
	week4ID := createWeek(t, ts, 4, cycleID)
	addDayToWeek(t, ts, week4ID, squatW4DayID, "MONDAY")
	addDayToWeek(t, ts, week4ID, benchW4DayID, "TUESDAY")
	addDayToWeek(t, ts, week4ID, deadliftW4DayID, "THURSDAY")
	addDayToWeek(t, ts, week4ID, ohpW4DayID, "FRIDAY")

	// =============================================================================
	// Create program
	// =============================================================================
	programSlug := "wendler-531-bbb-" + testID
	programBody := fmt.Sprintf(`{"name": "Wendler 5/3/1 BBB", "slug": "%s", "cycleId": "%s"}`,
		programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// Create Cycle Progression with default +5lb (upper body)
	// =============================================================================
	cycleProgBody := `{"name": "Wendler Cycle Progression", "type": "CYCLE_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX"}}`
	cycleProgResp, _ := adminPost(ts.URL("/progressions"), cycleProgBody)
	var cycleProgEnvelope ProgressionResponse
	json.NewDecoder(cycleProgResp.Body).Decode(&cycleProgEnvelope)
	cycleProgResp.Body.Close()
	cycleProgID := cycleProgEnvelope.Data.ID

	// Link progression with lift-specific overrides
	// Lower body: +10lb, Upper body: +5lb (default)
	linkProgressionToProgramWithOverride(t, ts, programID, cycleProgID, squatID, 1, 10.0)
	linkProgressionToProgramWithOverride(t, ts, programID, cycleProgID, deadliftID, 2, 10.0)
	linkProgressionToProgram(t, ts, programID, cycleProgID, benchID, 3)
	linkProgressionToProgram(t, ts, programID, cycleProgID, ohpID, 4)

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
	// EXECUTION PHASE: Week 1 (5s Week)
	// =============================================================================
	t.Run("Week 1 generates main work at 65%/75%/85% and BBB at 50%", func(t *testing.T) {
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

		// Should be squat day week 1
		if workout.Data.DaySlug != squatW1DaySlug {
			t.Errorf("Expected day slug '%s', got '%s'", squatW1DaySlug, workout.Data.DaySlug)
		}

		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got week %d", workout.Data.WeekNumber)
		}

		if len(workout.Data.Exercises) != 2 {
			t.Fatalf("Expected 2 exercises (main + BBB), got %d", len(workout.Data.Exercises))
		}

		// Find main work and BBB exercises
		var mainWork, bbbWork *WorkoutExerciseData
		for i := range workout.Data.Exercises {
			ex := &workout.Data.Exercises[i]
			if len(ex.Sets) == 3 {
				mainWork = ex
			} else if len(ex.Sets) == 5 {
				bbbWork = ex
			}
		}

		if mainWork == nil {
			t.Fatal("Missing main work exercise (3 sets)")
		}
		if bbbWork == nil {
			t.Fatal("Missing BBB exercise (5 sets)")
		}

		// Verify main work: 65%/75%/85% x 5 reps
		week1Percentages := []float64{0.65, 0.75, 0.85}
		for i, set := range mainWork.Sets {
			expectedWeight := squatTM * week1Percentages[i]
			if !closeEnough(set.Weight, expectedWeight) {
				t.Errorf("Main work set %d: expected weight ~%.1f (%.0f%% of %.1f), got %.1f",
					i+1, expectedWeight, week1Percentages[i]*100, squatTM, set.Weight)
			}
			if set.TargetReps != 5 {
				t.Errorf("Main work set %d: expected 5 reps, got %d", i+1, set.TargetReps)
			}
		}

		// Verify BBB: 5x10 at 50%
		expectedBBBWeight := squatTM * 0.50
		for i, set := range bbbWork.Sets {
			if !closeEnough(set.Weight, expectedBBBWeight) {
				t.Errorf("BBB set %d: expected weight ~%.1f (50%% of %.1f), got %.1f",
					i+1, expectedBBBWeight, squatTM, set.Weight)
			}
			if set.TargetReps != 10 {
				t.Errorf("BBB set %d: expected 10 reps, got %d", i+1, set.TargetReps)
			}
		}
	})

	// Advance through week 1 (4 days)
	for range 4 {
		advanceUserState(t, ts, userID)
	}

	// =============================================================================
	// EXECUTION PHASE: Week 2 (3s Week)
	// =============================================================================
	t.Run("Week 2 generates main work at 70%/80%/90% with 3 reps", func(t *testing.T) {
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

		if workout.Data.WeekNumber != 2 {
			t.Errorf("Expected week 2, got week %d", workout.Data.WeekNumber)
		}

		if workout.Data.DaySlug != squatW2DaySlug {
			t.Errorf("Expected day slug '%s', got '%s'", squatW2DaySlug, workout.Data.DaySlug)
		}

		var mainWork *WorkoutExerciseData
		for i := range workout.Data.Exercises {
			ex := &workout.Data.Exercises[i]
			if len(ex.Sets) == 3 {
				mainWork = ex
				break
			}
		}

		if mainWork == nil {
			t.Fatal("Missing main work exercise (3 sets)")
		}

		// Verify main work: 70%/80%/90% x 3 reps
		week2Percentages := []float64{0.70, 0.80, 0.90}
		for i, set := range mainWork.Sets {
			expectedWeight := squatTM * week2Percentages[i]
			if !closeEnough(set.Weight, expectedWeight) {
				t.Errorf("Main work set %d: expected weight ~%.1f (%.0f%% of %.1f), got %.1f",
					i+1, expectedWeight, week2Percentages[i]*100, squatTM, set.Weight)
			}
			if set.TargetReps != 3 {
				t.Errorf("Main work set %d: expected 3 reps, got %d", i+1, set.TargetReps)
			}
		}
	})

	// Advance through week 2
	for range 4 {
		advanceUserState(t, ts, userID)
	}

	// =============================================================================
	// EXECUTION PHASE: Week 3 (5/3/1 Week)
	// =============================================================================
	t.Run("Week 3 generates main work at 75%/85%/95% with 5/3/1 reps", func(t *testing.T) {
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

		if workout.Data.WeekNumber != 3 {
			t.Errorf("Expected week 3, got week %d", workout.Data.WeekNumber)
		}

		if workout.Data.DaySlug != squatW3DaySlug {
			t.Errorf("Expected day slug '%s', got '%s'", squatW3DaySlug, workout.Data.DaySlug)
		}

		var mainWork *WorkoutExerciseData
		for i := range workout.Data.Exercises {
			ex := &workout.Data.Exercises[i]
			if len(ex.Sets) == 3 {
				mainWork = ex
				break
			}
		}

		if mainWork == nil {
			t.Fatal("Missing main work exercise (3 sets)")
		}

		// Verify main work: 75%/85%/95% with 5/3/1 reps
		week3Percentages := []float64{0.75, 0.85, 0.95}
		week3Reps := []int{5, 3, 1}
		for i, set := range mainWork.Sets {
			expectedWeight := squatTM * week3Percentages[i]
			if !closeEnough(set.Weight, expectedWeight) {
				t.Errorf("Main work set %d: expected weight ~%.1f (%.0f%% of %.1f), got %.1f",
					i+1, expectedWeight, week3Percentages[i]*100, squatTM, set.Weight)
			}
			if set.TargetReps != week3Reps[i] {
				t.Errorf("Main work set %d: expected %d reps, got %d", i+1, week3Reps[i], set.TargetReps)
			}
		}
	})

	// Advance through week 3
	for range 4 {
		advanceUserState(t, ts, userID)
	}

	// =============================================================================
	// EXECUTION PHASE: Week 4 (Deload)
	// =============================================================================
	t.Run("Week 4 deload generates main work at 40%/50%/60%", func(t *testing.T) {
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

		if workout.Data.WeekNumber != 4 {
			t.Errorf("Expected week 4, got week %d", workout.Data.WeekNumber)
		}

		if workout.Data.DaySlug != squatW4DaySlug {
			t.Errorf("Expected day slug '%s', got '%s'", squatW4DaySlug, workout.Data.DaySlug)
		}

		var mainWork *WorkoutExerciseData
		for i := range workout.Data.Exercises {
			ex := &workout.Data.Exercises[i]
			if len(ex.Sets) == 3 {
				mainWork = ex
				break
			}
		}

		if mainWork == nil {
			t.Fatal("Missing main work exercise (3 sets)")
		}

		// Verify deload: 40%/50%/60% x 5 reps
		week4Percentages := []float64{0.40, 0.50, 0.60}
		for i, set := range mainWork.Sets {
			expectedWeight := squatTM * week4Percentages[i]
			if !closeEnough(set.Weight, expectedWeight) {
				t.Errorf("Deload set %d: expected weight ~%.1f (%.0f%% of %.1f), got %.1f",
					i+1, expectedWeight, week4Percentages[i]*100, squatTM, set.Weight)
			}
			if set.TargetReps != 5 {
				t.Errorf("Deload set %d: expected 5 reps, got %d", i+1, set.TargetReps)
			}
		}
	})

	// =============================================================================
	// PROGRESSION PHASE: Cycle progression with lift-specific overrides
	// =============================================================================
	t.Run("Cycle progression applies +10lb to lower body, +5lb to upper body", func(t *testing.T) {
		// Trigger progression for squat (lower body, +10lb override)
		triggerBody := ManualTriggerRequest{
			ProgressionID: cycleProgID,
			LiftID:        squatID,
			Force:         true,
		}
		triggerResp, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		if err != nil {
			t.Fatalf("Failed to trigger squat progression: %v", err)
		}
		var squatTrigger TriggerResponse
		json.NewDecoder(triggerResp.Body).Decode(&squatTrigger)
		triggerResp.Body.Close()

		if squatTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected squat progression to apply, got TotalApplied=%d", squatTrigger.Data.TotalApplied)
		}
		if len(squatTrigger.Data.Results) > 0 && squatTrigger.Data.Results[0].Result != nil {
			if squatTrigger.Data.Results[0].Result.Delta != 10.0 {
				t.Errorf("Expected squat delta +10 (override), got %f", squatTrigger.Data.Results[0].Result.Delta)
			}
			expectedNewSquat := squatTM + 10.0
			if squatTrigger.Data.Results[0].Result.NewValue != expectedNewSquat {
				t.Errorf("Expected squat new value %f, got %f", expectedNewSquat, squatTrigger.Data.Results[0].Result.NewValue)
			}
		}

		// Trigger progression for deadlift (lower body, +10lb override)
		triggerBody.LiftID = deadliftID
		triggerResp, err = authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		if err != nil {
			t.Fatalf("Failed to trigger deadlift progression: %v", err)
		}
		var deadliftTrigger TriggerResponse
		json.NewDecoder(triggerResp.Body).Decode(&deadliftTrigger)
		triggerResp.Body.Close()

		if deadliftTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected deadlift progression to apply")
		}
		if len(deadliftTrigger.Data.Results) > 0 && deadliftTrigger.Data.Results[0].Result != nil {
			if deadliftTrigger.Data.Results[0].Result.Delta != 10.0 {
				t.Errorf("Expected deadlift delta +10 (override), got %f", deadliftTrigger.Data.Results[0].Result.Delta)
			}
		}

		// Trigger progression for bench (upper body, default +5lb)
		triggerBody.LiftID = benchID
		triggerResp, err = authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		if err != nil {
			t.Fatalf("Failed to trigger bench progression: %v", err)
		}
		var benchTrigger TriggerResponse
		json.NewDecoder(triggerResp.Body).Decode(&benchTrigger)
		triggerResp.Body.Close()

		if benchTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected bench progression to apply")
		}
		if len(benchTrigger.Data.Results) > 0 && benchTrigger.Data.Results[0].Result != nil {
			if benchTrigger.Data.Results[0].Result.Delta != 5.0 {
				t.Errorf("Expected bench delta +5 (default), got %f", benchTrigger.Data.Results[0].Result.Delta)
			}
		}

		// Trigger progression for OHP (upper body, default +5lb)
		triggerBody.LiftID = ohpID
		triggerResp, err = authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		if err != nil {
			t.Fatalf("Failed to trigger OHP progression: %v", err)
		}
		var ohpTrigger TriggerResponse
		json.NewDecoder(triggerResp.Body).Decode(&ohpTrigger)
		triggerResp.Body.Close()

		if ohpTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected OHP progression to apply")
		}
		if len(ohpTrigger.Data.Results) > 0 && ohpTrigger.Data.Results[0].Result != nil {
			if ohpTrigger.Data.Results[0].Result.Delta != 5.0 {
				t.Errorf("Expected OHP delta +5 (default), got %f", ohpTrigger.Data.Results[0].Result.Delta)
			}
		}
	})

	// Advance to cycle 2
	for range 4 {
		advanceUserState(t, ts, userID)
	}

	// =============================================================================
	// VALIDATION PHASE: Cycle 2 uses new training maxes
	// =============================================================================
	t.Run("Cycle 2 Week 1 uses increased training maxes", func(t *testing.T) {
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

		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1 (new cycle), got week %d", workout.Data.WeekNumber)
		}

		if workout.Data.CycleIteration != 2 {
			t.Errorf("Expected cycle iteration 2, got %d", workout.Data.CycleIteration)
		}

		var mainWork *WorkoutExerciseData
		for i := range workout.Data.Exercises {
			ex := &workout.Data.Exercises[i]
			if len(ex.Sets) == 3 {
				mainWork = ex
				break
			}
		}

		if mainWork == nil {
			t.Fatal("Missing main work exercise (3 sets)")
		}

		// New squat TM should be 310 (300 + 10)
		newSquatTM := squatTM + 10.0
		week1Percentages := []float64{0.65, 0.75, 0.85}
		for i, set := range mainWork.Sets {
			expectedWeight := newSquatTM * week1Percentages[i]
			if !closeEnough(set.Weight, expectedWeight) {
				t.Errorf("Cycle 2 main work set %d: expected weight ~%.1f (%.0f%% of new TM %.1f), got %.1f",
					i+1, expectedWeight, week1Percentages[i]*100, newSquatTM, set.Weight)
			}
		}

		// BBB should also use new TM
		var bbbWork *WorkoutExerciseData
		for i := range workout.Data.Exercises {
			ex := &workout.Data.Exercises[i]
			if len(ex.Sets) == 5 {
				bbbWork = ex
				break
			}
		}

		if bbbWork != nil {
			expectedBBBWeight := newSquatTM * 0.50
			for i, set := range bbbWork.Sets {
				if !closeEnough(set.Weight, expectedBBBWeight) {
					t.Errorf("Cycle 2 BBB set %d: expected weight ~%.1f (50%% of new TM %.1f), got %.1f",
						i+1, expectedBBBWeight, newSquatTM, set.Weight)
				}
			}
		}
	})
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// create531RampPrescription creates a RAMP prescription for 5/3/1 main work.
func create531RampPrescription(t *testing.T, ts *testutil.TestServer, liftID string, steps []RampStep, order int) string {
	t.Helper()

	stepsJSON, err := json.Marshal(steps)
	if err != nil {
		t.Fatalf("Failed to marshal ramp steps: %v", err)
	}

	// Use 60% work set threshold so all 3 main work sets are classified as work sets
	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": 100.0},
		"setScheme": {"type": "RAMP", "steps": %s, "workSetThreshold": 60.0},
		"order": %d
	}`, liftID, string(stepsJSON), order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create 5/3/1 ramp prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create 5/3/1 ramp prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// createDay creates a day and returns its ID.
func createDay(t *testing.T, ts *testutil.TestServer, name, slug string) string {
	t.Helper()
	body := fmt.Sprintf(`{"name": "%s", "slug": "%s"}`, name, slug)
	resp, _ := adminPost(ts.URL("/days"), body)
	var envelope DayResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	resp.Body.Close()
	return envelope.Data.ID
}

// createWeek creates a week in a cycle and returns its ID.
func createWeek(t *testing.T, ts *testutil.TestServer, weekNumber int, cycleID string) string {
	t.Helper()
	body := fmt.Sprintf(`{"weekNumber": %d, "cycleId": "%s"}`, weekNumber, cycleID)
	resp, _ := adminPost(ts.URL("/weeks"), body)
	var envelope WeekResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	resp.Body.Close()
	return envelope.Data.ID
}

// linkProgressionToProgramWithOverride links a progression to a program with an override increment.
func linkProgressionToProgramWithOverride(t *testing.T, ts *testutil.TestServer, programID, progressionID, liftID string, priority int, overrideIncrement float64) {
	t.Helper()
	body := fmt.Sprintf(`{"progressionId": "%s", "liftId": "%s", "priority": %d, "enabled": true, "overrideIncrement": %f}`,
		progressionID, liftID, priority, overrideIncrement)
	resp, err := adminPost(ts.URL("/programs/"+programID+"/progressions"), body)
	if err != nil {
		t.Fatalf("Failed to link progression to program with override: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to link progression to program with override, status %d: %s", resp.StatusCode, bodyBytes)
	}
}
