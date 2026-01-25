// Package e2e provides end-to-end tests for complete program workflows.
// This file provides E2E tests for Phase 5 rotation programs: CAP3, Inverted Juggernaut, GreySkull LP.
package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/waynenilsen/power-pro-v3/internal/testutil"
)

// =============================================================================
// GREYSKULL LP E2E TEST
// =============================================================================

// TestGreySkullLPProgram validates the complete GreySkull LP program through the API.
//
// GreySkull LP characteristics:
// - A/B Rotation: Week parity determines day variant (Week 1: A,B,A; Week 2: B,A,B)
// - Main Lifts: 2x5 + 1x5+ (AMRAP) using GreySkull set scheme
// - AMRAP-Based Progression:
//   - <5 reps: 10% deload
//   - 5-9 reps: +2.5lb (upper) / +5lb (lower)
//   - 10+ reps: double progression (+5lb upper / +10lb lower)
func TestGreySkullLPProgram(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user" // Uses pre-seeded test user

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	// Create OHP lift
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

	// GreySkull training maxes
	benchTM := 135.0
	ohpTM := 95.0
	squatTM := 185.0
	deadliftTM := 225.0

	// Create training maxes
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchTM)
	createLiftMax(t, ts, userID, ohpID, "TRAINING_MAX", ohpTM)
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatTM)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftTM)

	// =============================================================================
	// Create GreySkull set scheme prescriptions (2x5 + 1x5+)
	// =============================================================================

	// Variant A Day: Bench, Row, Squat
	benchPrescID := createGreySkullPrescription(t, ts, benchID, 2, 5, 1, 5, 100.0, 0)
	squatPrescID := createGreySkullPrescription(t, ts, squatID, 2, 5, 1, 5, 100.0, 2)

	// Variant B Day: OHP, Chinups, Deadlift
	ohpPrescID := createGreySkullPrescription(t, ts, ohpID, 2, 5, 1, 5, 100.0, 0)
	deadliftPrescID := createGreySkullPrescription(t, ts, deadliftID, 2, 5, 1, 5, 100.0, 2)

	// =============================================================================
	// Create Days (Variant A and Variant B)
	// =============================================================================

	// Variant A Day (Bench day)
	variantASlug := "variant-a-" + testID
	variantABody := fmt.Sprintf(`{"name": "Variant A", "slug": "%s"}`, variantASlug)
	variantAResp, _ := adminPost(ts.URL("/days"), variantABody)
	var variantAEnvelope DayResponse
	json.NewDecoder(variantAResp.Body).Decode(&variantAEnvelope)
	variantAResp.Body.Close()
	variantADayID := variantAEnvelope.Data.ID

	addPrescToDay(t, ts, variantADayID, benchPrescID)
	addPrescToDay(t, ts, variantADayID, squatPrescID)

	// Variant B Day (OHP day)
	variantBSlug := "variant-b-" + testID
	variantBBody := fmt.Sprintf(`{"name": "Variant B", "slug": "%s"}`, variantBSlug)
	variantBResp, _ := adminPost(ts.URL("/days"), variantBBody)
	var variantBEnvelope DayResponse
	json.NewDecoder(variantBResp.Body).Decode(&variantBEnvelope)
	variantBResp.Body.Close()
	variantBDayID := variantBEnvelope.Data.ID

	addPrescToDay(t, ts, variantBDayID, ohpPrescID)
	addPrescToDay(t, ts, variantBDayID, deadliftPrescID)

	// =============================================================================
	// Create 2-week cycle (A/B rotation pattern)
	// =============================================================================
	cycleName := "GreySkull LP Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 2}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	// Week 1: A, B, A pattern
	week1Body := fmt.Sprintf(`{"weekNumber": 1, "cycleId": "%s"}`, cycleID)
	week1Resp, _ := adminPost(ts.URL("/weeks"), week1Body)
	var week1Envelope WeekResponse
	json.NewDecoder(week1Resp.Body).Decode(&week1Envelope)
	week1Resp.Body.Close()
	week1ID := week1Envelope.Data.ID

	addDayToWeek(t, ts, week1ID, variantADayID, "MONDAY")
	addDayToWeek(t, ts, week1ID, variantBDayID, "WEDNESDAY")
	addDayToWeek(t, ts, week1ID, variantADayID, "FRIDAY")

	// Week 2: B, A, B pattern
	week2Body := fmt.Sprintf(`{"weekNumber": 2, "cycleId": "%s"}`, cycleID)
	week2Resp, _ := adminPost(ts.URL("/weeks"), week2Body)
	var week2Envelope WeekResponse
	json.NewDecoder(week2Resp.Body).Decode(&week2Envelope)
	week2Resp.Body.Close()
	week2ID := week2Envelope.Data.ID

	addDayToWeek(t, ts, week2ID, variantBDayID, "MONDAY")
	addDayToWeek(t, ts, week2ID, variantADayID, "WEDNESDAY")
	addDayToWeek(t, ts, week2ID, variantBDayID, "FRIDAY")

	// =============================================================================
	// Create Program
	// =============================================================================
	programSlug := "greyskull-lp-" + testID
	programBody := fmt.Sprintf(`{"name": "GreySkull LP", "slug": "%s", "cycleId": "%s"}`,
		programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// Create GreySkull Progressions
	// =============================================================================

	// Upper body progression (+2.5lb standard, +5lb double)
	upperProgBody := `{"name": "GreySkull Upper", "type": "GREYSKULL_PROGRESSION", "parameters": {"increment": 2.5, "minReps": 5, "doubleThreshold": 10, "deloadPercent": 0.10, "maxType": "TRAINING_MAX"}}`
	upperProgResp, _ := adminPost(ts.URL("/progressions"), upperProgBody)
	var upperProgEnvelope ProgressionResponse
	json.NewDecoder(upperProgResp.Body).Decode(&upperProgEnvelope)
	upperProgResp.Body.Close()
	upperProgID := upperProgEnvelope.Data.ID

	// Lower body progression (+5lb standard, +10lb double)
	lowerProgBody := `{"name": "GreySkull Lower", "type": "GREYSKULL_PROGRESSION", "parameters": {"increment": 5.0, "minReps": 5, "doubleThreshold": 10, "deloadPercent": 0.10, "maxType": "TRAINING_MAX"}}`
	lowerProgResp, _ := adminPost(ts.URL("/progressions"), lowerProgBody)
	var lowerProgEnvelope ProgressionResponse
	json.NewDecoder(lowerProgResp.Body).Decode(&lowerProgEnvelope)
	lowerProgResp.Body.Close()
	lowerProgID := lowerProgEnvelope.Data.ID

	// Link progressions to program
	linkProgressionToProgram(t, ts, programID, upperProgID, benchID, 1)
	linkProgressionToProgram(t, ts, programID, upperProgID, ohpID, 2)
	linkProgressionToProgram(t, ts, programID, lowerProgID, squatID, 3)
	linkProgressionToProgram(t, ts, programID, lowerProgID, deadliftID, 4)

	// =============================================================================
	// Enroll user
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
	// Test: Week 1 Day 1 (Variant A) generates correct workout
	// =============================================================================
	t.Run("Week 1 Day 1 - Variant A workout", func(t *testing.T) {
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

		// Verify Week 1
		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		// Verify Variant A day
		if workout.Data.DaySlug != variantASlug {
			t.Errorf("Expected day '%s', got '%s'", variantASlug, workout.Data.DaySlug)
		}

		// Verify exercises present (Bench, Squat on Variant A)
		if len(workout.Data.Exercises) < 2 {
			t.Fatalf("Expected at least 2 exercises, got %d", len(workout.Data.Exercises))
		}

		// Verify GreySkull set structure (3 sets total: 2 fixed + 1 AMRAP)
		for _, ex := range workout.Data.Exercises {
			if len(ex.Sets) != 3 {
				t.Errorf("Exercise %s: expected 3 sets (2+1 AMRAP), got %d", ex.Lift.Name, len(ex.Sets))
			}
		}
	})

	// Complete Week 1 Day 1 and advance to Week 1 Day 2 (Variant B)
	completeRotationWorkoutDay(t, ts, userID)

	t.Run("Week 1 Day 2 - Variant B workout", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Should still be Week 1
		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		// Verify Variant B day
		if workout.Data.DaySlug != variantBSlug {
			t.Errorf("Expected day '%s', got '%s'", variantBSlug, workout.Data.DaySlug)
		}
	})

	// Complete remaining Week 1 workouts and advance to Week 2
	completeRotationWorkoutDay(t, ts, userID) // Week 1 Day 3 (Variant A)
	completeRotationWorkoutDay(t, ts, userID) // Week 2 Day 1 (Variant B)

	t.Run("Week 2 Day 1 - Variant B workout (pattern flip)", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Should be Week 2
		if workout.Data.WeekNumber != 2 {
			t.Errorf("Expected week 2, got %d", workout.Data.WeekNumber)
		}

		// Week 2 Day 1 should be Variant B (pattern flipped)
		if workout.Data.DaySlug != variantBSlug {
			t.Errorf("Expected Variant B on Week 2 Day 1, got '%s'", workout.Data.DaySlug)
		}
	})
}

// =============================================================================
// HELPER FUNCTIONS (Phase 5 specific)
// =============================================================================

// createGreySkullPrescription creates a prescription with GreySkull set scheme (fixed + AMRAP).
func createGreySkullPrescription(t *testing.T, ts *testutil.TestServer, liftID string, fixedSets, fixedReps, amrapSets, minAmrapReps int, percentage float64, order int) string {
	t.Helper()

	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": %.1f},
		"setScheme": {"type": "GREYSKULL", "fixedSets": %d, "fixedReps": %d, "amrapSets": %d, "minAmrapReps": %d},
		"order": %d
	}`, liftID, percentage, fixedSets, fixedReps, amrapSets, minAmrapReps, order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create GreySkull prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create GreySkull prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}

// Note: createLiftMax, addPrescToDay, addDayToWeek, linkProgressionToProgram, and advanceUserState
// are defined in starting_strength_test.go and shared across all E2E tests in this package.

// =============================================================================
// CAP3 E2E TEST
// =============================================================================

// TestCAP3Program validates the complete CAP3 program through the API.
//
// CAP3 characteristics:
// - 3-week rotation cycle: Each lift gets AMRAP focus once per cycle
// - Week 1: Deadlift AMRAP focus, Squat/Bench volume
// - Week 2: Squat AMRAP focus, Deadlift/Bench medium
// - Week 3: Bench AMRAP focus, Deadlift/Squat volume
// - High Intensity: 79.5% x 6, 83.5% x 4, 88.5% x 2+ (AMRAP)
func TestCAP3Program(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "bill-starr-test-user" // Uses pre-seeded test user

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	// CAP3 training maxes
	squatTM := 315.0
	benchTM := 225.0
	deadliftTM := 405.0

	// Create training maxes
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatTM)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchTM)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftTM)

	// =============================================================================
	// Create Weekly Lookups for CAP3 intensity levels
	// =============================================================================

	// High Intensity AMRAP lookup: 79.5% x 6, 83.5% x 4, 88.5% x 2+
	highIntensityLookup := `{
		"name": "CAP3 High Intensity",
		"entries": [
			{"weekNumber": 1, "percentages": [79.5, 83.5, 88.5], "reps": [6, 4, 2]}
		]
	}`
	hiResp, err := adminPost(ts.URL("/weekly-lookups"), highIntensityLookup)
	if err != nil {
		t.Fatalf("Failed to create high intensity lookup: %v", err)
	}
	var hiEnvelope struct {
		Data struct{ ID string `json:"id"` } `json:"data"`
	}
	json.NewDecoder(hiResp.Body).Decode(&hiEnvelope)
	hiResp.Body.Close()
	_ = hiEnvelope.Data.ID // Used by program configuration

	// =============================================================================
	// Create Rotation Lookup for CAP3 (3-position rotation)
	// =============================================================================
	rotationBody := `{
		"name": "CAP3 Rotation",
		"entries": [
			{"position": 0, "liftIdentifier": "deadlift", "description": "Deadlift Focus Week"},
			{"position": 1, "liftIdentifier": "squat", "description": "Squat Focus Week"},
			{"position": 2, "liftIdentifier": "bench", "description": "Bench Focus Week"}
		]
	}`
	rotResp, err := adminPost(ts.URL("/rotation-lookups"), rotationBody)
	if err != nil {
		t.Fatalf("Failed to create rotation lookup: %v", err)
	}
	var rotEnvelope struct {
		Data struct{ ID string `json:"id"` } `json:"data"`
	}
	json.NewDecoder(rotResp.Body).Decode(&rotEnvelope)
	rotResp.Body.Close()
	_ = rotEnvelope.Data.ID // Used by program configuration

	// =============================================================================
	// Create prescriptions (AMRAP sets for high intensity)
	// =============================================================================
	deadliftPrescID := createAMRAPPrescription(t, ts, deadliftID, 1, 2, 88.5, 0)
	squatPrescID := createAMRAPPrescription(t, ts, squatID, 1, 2, 88.5, 1)
	benchPrescID := createAMRAPPrescription(t, ts, benchID, 1, 2, 88.5, 2)

	// =============================================================================
	// Create Days for each lift focus
	// =============================================================================

	// Deadlift Focus Day
	dlDaySlug := "deadlift-focus-" + testID
	dlDayBody := fmt.Sprintf(`{"name": "Deadlift Focus", "slug": "%s"}`, dlDaySlug)
	dlDayResp, _ := adminPost(ts.URL("/days"), dlDayBody)
	var dlDayEnvelope DayResponse
	json.NewDecoder(dlDayResp.Body).Decode(&dlDayEnvelope)
	dlDayResp.Body.Close()
	dlDayID := dlDayEnvelope.Data.ID
	addPrescToDay(t, ts, dlDayID, deadliftPrescID)

	// Squat Focus Day
	sqDaySlug := "squat-focus-" + testID
	sqDayBody := fmt.Sprintf(`{"name": "Squat Focus", "slug": "%s"}`, sqDaySlug)
	sqDayResp, _ := adminPost(ts.URL("/days"), sqDayBody)
	var sqDayEnvelope DayResponse
	json.NewDecoder(sqDayResp.Body).Decode(&sqDayEnvelope)
	sqDayResp.Body.Close()
	sqDayID := sqDayEnvelope.Data.ID
	addPrescToDay(t, ts, sqDayID, squatPrescID)

	// Bench Focus Day
	bnDaySlug := "bench-focus-" + testID
	bnDayBody := fmt.Sprintf(`{"name": "Bench Focus", "slug": "%s"}`, bnDaySlug)
	bnDayResp, _ := adminPost(ts.URL("/days"), bnDayBody)
	var bnDayEnvelope DayResponse
	json.NewDecoder(bnDayResp.Body).Decode(&bnDayEnvelope)
	bnDayResp.Body.Close()
	bnDayID := bnDayEnvelope.Data.ID
	addPrescToDay(t, ts, bnDayID, benchPrescID)

	// =============================================================================
	// Create 3-week cycle
	// =============================================================================
	cycleName := "CAP3 Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 3}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	// Week 1: Deadlift focus
	week1Body := fmt.Sprintf(`{"weekNumber": 1, "cycleId": "%s"}`, cycleID)
	week1Resp, _ := adminPost(ts.URL("/weeks"), week1Body)
	var week1Envelope WeekResponse
	json.NewDecoder(week1Resp.Body).Decode(&week1Envelope)
	week1Resp.Body.Close()
	addDayToWeek(t, ts, week1Envelope.Data.ID, dlDayID, "MONDAY")

	// Week 2: Squat focus
	week2Body := fmt.Sprintf(`{"weekNumber": 2, "cycleId": "%s"}`, cycleID)
	week2Resp, _ := adminPost(ts.URL("/weeks"), week2Body)
	var week2Envelope WeekResponse
	json.NewDecoder(week2Resp.Body).Decode(&week2Envelope)
	week2Resp.Body.Close()
	addDayToWeek(t, ts, week2Envelope.Data.ID, sqDayID, "MONDAY")

	// Week 3: Bench focus
	week3Body := fmt.Sprintf(`{"weekNumber": 3, "cycleId": "%s"}`, cycleID)
	week3Resp, _ := adminPost(ts.URL("/weeks"), week3Body)
	var week3Envelope WeekResponse
	json.NewDecoder(week3Resp.Body).Decode(&week3Envelope)
	week3Resp.Body.Close()
	addDayToWeek(t, ts, week3Envelope.Data.ID, bnDayID, "MONDAY")

	// =============================================================================
	// Create Program
	// =============================================================================
	programSlug := "cap3-" + testID
	programBody := fmt.Sprintf(`{"name": "nSuns CAP3", "slug": "%s", "cycleId": "%s"}`,
		programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// Enroll user
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
	// Test: Week 1 - Deadlift focus
	// =============================================================================
	t.Run("Week 1 - Deadlift AMRAP focus", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		if workout.Data.DaySlug != dlDaySlug {
			t.Errorf("Expected deadlift focus day, got '%s'", workout.Data.DaySlug)
		}

		// Verify deadlift is present
		if len(workout.Data.Exercises) < 1 {
			t.Fatalf("Expected at least 1 exercise")
		}

		// Verify weight is approximately 88.5% of TM
		dl := workout.Data.Exercises[0]
		expectedWeight := deadliftTM * 0.885 // 405 * 0.885 = 358.425
		if !withinTolerance(dl.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected weight ~%.1f (88.5%% of TM), got %.1f", expectedWeight, dl.Sets[0].Weight)
		}
	})

	// Complete Week 1 and advance to Week 2
	completeRotationWorkoutDay(t, ts, userID)

	t.Run("Week 2 - Squat AMRAP focus", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		if workout.Data.WeekNumber != 2 {
			t.Errorf("Expected week 2, got %d", workout.Data.WeekNumber)
		}

		if workout.Data.DaySlug != sqDaySlug {
			t.Errorf("Expected squat focus day, got '%s'", workout.Data.DaySlug)
		}
	})

	// Complete Week 2 and advance to Week 3
	completeRotationWorkoutDay(t, ts, userID)

	t.Run("Week 3 - Bench AMRAP focus", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		if workout.Data.WeekNumber != 3 {
			t.Errorf("Expected week 3, got %d", workout.Data.WeekNumber)
		}

		if workout.Data.DaySlug != bnDaySlug {
			t.Errorf("Expected bench focus day, got '%s'", workout.Data.DaySlug)
		}
	})

	// Complete Week 3 and advance to new cycle (Week 1 again)
	completeRotationWorkoutDay(t, ts, userID)

	t.Run("Rotation cycles back to Week 1", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Should be back to Week 1
		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected rotation to cycle back to week 1, got %d", workout.Data.WeekNumber)
		}

		// Should be Deadlift focus again
		if workout.Data.DaySlug != dlDaySlug {
			t.Errorf("Expected deadlift focus after cycle, got '%s'", workout.Data.DaySlug)
		}
	})
}

// =============================================================================
// INVERTED JUGGERNAUT E2E TEST
// =============================================================================

// TestInvertedJuggernautProgram validates the complete Inverted Juggernaut program through the API.
//
// Inverted Juggernaut characteristics:
// - 16-week cycle with 4 waves (10s, 8s, 5s, 3s)
// - Each wave has 4 phases: Accumulation, Intensification, Realization, Deload
// - Volume sets vary by wave (9/7/5/6 sets)
// - 5/3/1 overlay: 65/75/85, 70/80/90, 75/85/95, 40/50/60
// - TM progression based on rep standard performance
func TestInvertedJuggernautProgram(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "wendler-531-test-user" // Uses pre-seeded test user

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"

	// Juggernaut training max
	squatTM := 315.0

	// Create training max
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatTM)

	// =============================================================================
	// Create Weekly Lookup for 5/3/1 overlay (4-week repeating pattern)
	// =============================================================================
	weeklyLookupBody := `{
		"name": "Inverted Juggernaut 5/3/1",
		"entries": [
			{"weekNumber": 1, "percentages": [65.0, 75.0, 85.0, 75.0, 65.0], "reps": [5, 5, -5, 5, -5]},
			{"weekNumber": 2, "percentages": [70.0, 80.0, 90.0, 80.0, 70.0], "reps": [3, 3, -3, 3, -3]},
			{"weekNumber": 3, "percentages": [75.0, 85.0, 95.0, 85.0, 75.0], "reps": [5, 3, -1, 3, -5]},
			{"weekNumber": 4, "percentages": [40.0, 50.0, 60.0], "reps": [5, 5, 5]}
		]
	}`
	wlResp, err := adminPost(ts.URL("/weekly-lookups"), weeklyLookupBody)
	if err != nil {
		t.Fatalf("Failed to create weekly lookup: %v", err)
	}
	var wlEnvelope struct {
		Data struct{ ID string `json:"id"` } `json:"data"`
	}
	json.NewDecoder(wlResp.Body).Decode(&wlEnvelope)
	wlResp.Body.Close()
	weeklyLookupID := wlEnvelope.Data.ID

	// =============================================================================
	// Create prescription for squat
	// =============================================================================
	squatPrescID := createAMRAPPrescription(t, ts, squatID, 1, 5, 85.0, 0)

	// =============================================================================
	// Create Day
	// =============================================================================
	squatDaySlug := "squat-day-" + testID
	squatDayBody := fmt.Sprintf(`{"name": "Squat Day", "slug": "%s"}`, squatDaySlug)
	squatDayResp, _ := adminPost(ts.URL("/days"), squatDayBody)
	var squatDayEnvelope DayResponse
	json.NewDecoder(squatDayResp.Body).Decode(&squatDayEnvelope)
	squatDayResp.Body.Close()
	squatDayID := squatDayEnvelope.Data.ID
	addPrescToDay(t, ts, squatDayID, squatPrescID)

	// =============================================================================
	// Create 16-week cycle (simplified: just testing key weeks)
	// =============================================================================
	cycleName := "Juggernaut Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 16}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	// Create all 16 weeks
	for w := 1; w <= 16; w++ {
		weekBody := fmt.Sprintf(`{"weekNumber": %d, "cycleId": "%s"}`, w, cycleID)
		weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
		var weekEnvelope WeekResponse
		json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
		weekResp.Body.Close()
		addDayToWeek(t, ts, weekEnvelope.Data.ID, squatDayID, "MONDAY")
	}

	// =============================================================================
	// Create Program with weekly lookup
	// =============================================================================
	programSlug := "juggernaut-" + testID
	programBody := fmt.Sprintf(`{"name": "Inverted Juggernaut", "slug": "%s", "cycleId": "%s", "weeklyLookupId": "%s"}`,
		programSlug, cycleID, weeklyLookupID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// Enroll user
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
	// Test: Week 1 - 10s Wave Accumulation (85% on set 3)
	// =============================================================================
	t.Run("Week 1 - 10s Wave Accumulation", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		if workout.Data.WeekNumber != 1 {
			t.Errorf("Expected week 1, got %d", workout.Data.WeekNumber)
		}

		// Verify we have exercises
		if len(workout.Data.Exercises) < 1 {
			t.Fatalf("Expected at least 1 exercise")
		}

		// Verify weight is approximately 85% of TM (from weekly lookup)
		squat := workout.Data.Exercises[0]
		expectedWeight := squatTM * 0.85 // 315 * 0.85 = 267.75
		if !withinTolerance(squat.Sets[0].Weight, expectedWeight, 5.0) {
			t.Errorf("Expected weight ~%.1f (85%% of TM), got %.1f", expectedWeight, squat.Sets[0].Weight)
		}
	})

	// Complete workouts to advance to Week 4 (Deload)
	for i := 0; i < 3; i++ {
		completeRotationWorkoutDay(t, ts, userID)
	}

	t.Run("Week 4 - 10s Wave Deload", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		if workout.Data.WeekNumber != 4 {
			t.Errorf("Expected week 4 (deload), got %d", workout.Data.WeekNumber)
		}
	})

	// Complete Week 4 and advance to Week 5 (8s Wave start)
	completeRotationWorkoutDay(t, ts, userID)

	t.Run("Week 5 - 8s Wave Accumulation", func(t *testing.T) {
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Should be Week 5 (first week of 8s wave)
		if workout.Data.WeekNumber != 5 {
			t.Errorf("Expected week 5 (8s wave start), got %d", workout.Data.WeekNumber)
		}
	})
}

// =============================================================================
// SHARED HELPER FUNCTIONS
// =============================================================================

// completeRotationWorkoutDay completes a rotation-based workout day using explicit state machine flow.
// This function starts a session, logs all sets, finishes the session, and advances to the next day.
func completeRotationWorkoutDay(t *testing.T, ts *testutil.TestServer, userID string) {
	t.Helper()

	sessionID := startWorkoutSession(t, ts, userID)

	// Get the workout to find prescription IDs
	workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
	var workout WorkoutResponse
	json.NewDecoder(workoutResp.Body).Decode(&workout)
	workoutResp.Body.Close()

	// Log sets for each exercise
	for _, ex := range workout.Data.Exercises {
		for _, set := range ex.Sets {
			logRotationSet(t, ts, userID, sessionID, ex.PrescriptionID, ex.Lift.ID, set.SetNumber, set.Weight, set.TargetReps, set.TargetReps)
		}
	}

	finishWorkoutSession(t, ts, sessionID, userID)

	// Advance to next day
	advanceUserState(t, ts, userID)
}

// logRotationSet logs a single set for rotation-based program workouts.
func logRotationSet(t *testing.T, ts *testutil.TestServer, userID, sessionID, prescriptionID, liftID string, setNumber int, weight float64, targetReps, repsPerformed int) {
	t.Helper()

	type setRequest struct {
		PrescriptionID string  `json:"prescriptionId"`
		LiftID         string  `json:"liftId"`
		SetNumber      int     `json:"setNumber"`
		Weight         float64 `json:"weight"`
		TargetReps     int     `json:"targetReps"`
		RepsPerformed  int     `json:"repsPerformed"`
	}

	setsReq := []setRequest{{
		PrescriptionID: prescriptionID,
		LiftID:         liftID,
		SetNumber:      setNumber,
		Weight:         weight,
		TargetReps:     targetReps,
		RepsPerformed:  repsPerformed,
	}}

	body, _ := json.Marshal(map[string]interface{}{"sets": setsReq})
	req, _ := http.NewRequest(http.MethodPost, ts.URL("/sessions/"+sessionID+"/sets"), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to log set: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to log set, status %d: %s", resp.StatusCode, respBody)
	}
}
