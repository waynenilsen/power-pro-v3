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
// GZCL COMPENDIUM (VDIP) E2E TEST
// =============================================================================

// TestGZCLCompendiumVDIPProgram validates the complete GZCL Compendium VDIP
// (Volume-Dependent Intensity Progression) program configuration and execution.
//
// VDIP characteristics:
// - 5 days per week, ongoing (no fixed duration)
// - Three-tier system: T1, T2, T3
// - T1: 3 MRS (Max Rep Sets) at 85% TM
// - T2: 3 MRS at 65% TM
// - T3: 4 MRS each (fixed weight based on 10RM)
//
// VDIP Progression Rules (based on total reps achieved):
// - T1: 15+ reps = +10lb, 10-14 reps = +5lb, <10 reps = maintain
// - T2: 30+ reps = +10lb, 25-29 reps = +5lb, <25 reps = maintain
// - T3: 50+ reps = +5lb, <50 reps = maintain
//
// Sample Week Layout:
// - Day 1: Back Squat (T1), Stiff Leg Deadlift (T2), T3 accessories
// - Day 2: Bench Press (T1), Spoto Bench (T2), T3 accessories
// - Day 3: Deadlift (T1), Front Squat (T2), T3 accessories
// - Day 4: OHP (T1), Push Press (T2), T3 accessories
// - Day 5: Front Squat (T1), Back Squat (T2), T3 accessories
func TestGZCLCompendiumVDIPProgram(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Test-unique identifiers
	testID := uuid.New().String()[:8]
	userID := "workout-test-user" // Uses seeded test user

	// Seeded lift IDs
	squatID := "00000000-0000-0000-0000-000000000001"
	benchID := "00000000-0000-0000-0000-000000000002"
	deadliftID := "00000000-0000-0000-0000-000000000003"

	// Create additional lifts (OHP, Front Squat, Stiff Leg DL, etc. are not seeded)
	ohpSlug := "ohp-vdip-" + testID
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

	// Front Squat (for T1 Day 5 and T2 Day 3)
	frontSquatSlug := "front-squat-vdip-" + testID
	frontSquatBody := fmt.Sprintf(`{"name": "Front Squat", "slug": "%s", "isCompetitionLift": false}`, frontSquatSlug)
	frontSquatResp, _ := adminPost(ts.URL("/lifts"), frontSquatBody)
	var frontSquatEnvelope LiftResponse
	json.NewDecoder(frontSquatResp.Body).Decode(&frontSquatEnvelope)
	frontSquatResp.Body.Close()
	frontSquatID := frontSquatEnvelope.Data.ID

	// Stiff Leg Deadlift (T2 Day 1)
	sldlSlug := "sldl-vdip-" + testID
	sldlBody := fmt.Sprintf(`{"name": "Stiff Leg Deadlift", "slug": "%s", "isCompetitionLift": false}`, sldlSlug)
	sldlResp, _ := adminPost(ts.URL("/lifts"), sldlBody)
	var sldlEnvelope LiftResponse
	json.NewDecoder(sldlResp.Body).Decode(&sldlEnvelope)
	sldlResp.Body.Close()
	sldlID := sldlEnvelope.Data.ID

	// Spoto Bench (T2 Day 2)
	spotoSlug := "spoto-bench-vdip-" + testID
	spotoBody := fmt.Sprintf(`{"name": "Spoto Bench", "slug": "%s", "isCompetitionLift": false}`, spotoSlug)
	spotoResp, _ := adminPost(ts.URL("/lifts"), spotoBody)
	var spotoEnvelope LiftResponse
	json.NewDecoder(spotoResp.Body).Decode(&spotoEnvelope)
	spotoResp.Body.Close()
	spotoID := spotoEnvelope.Data.ID

	// Push Press (T2 Day 4)
	pushPressSlug := "push-press-vdip-" + testID
	pushPressBody := fmt.Sprintf(`{"name": "Push Press", "slug": "%s", "isCompetitionLift": false}`, pushPressSlug)
	pushPressResp, _ := adminPost(ts.URL("/lifts"), pushPressBody)
	var pushPressEnvelope LiftResponse
	json.NewDecoder(pushPressResp.Body).Decode(&pushPressEnvelope)
	pushPressResp.Body.Close()
	pushPressID := pushPressEnvelope.Data.ID

	// T3 accessories: Lunges, Pull-ups (as placeholder exercises)
	lungesSlug := "lunges-vdip-" + testID
	lungesBody := fmt.Sprintf(`{"name": "Lunges", "slug": "%s", "isCompetitionLift": false}`, lungesSlug)
	lungesResp, _ := adminPost(ts.URL("/lifts"), lungesBody)
	var lungesEnvelope LiftResponse
	json.NewDecoder(lungesResp.Body).Decode(&lungesEnvelope)
	lungesResp.Body.Close()
	lungesID := lungesEnvelope.Data.ID

	pullUpsSlug := "pull-ups-vdip-" + testID
	pullUpsBody := fmt.Sprintf(`{"name": "Pull-ups", "slug": "%s", "isCompetitionLift": false}`, pullUpsSlug)
	pullUpsResp, _ := adminPost(ts.URL("/lifts"), pullUpsBody)
	var pullUpsEnvelope LiftResponse
	json.NewDecoder(pullUpsResp.Body).Decode(&pullUpsEnvelope)
	pullUpsResp.Body.Close()
	pullUpsID := pullUpsEnvelope.Data.ID

	// VDIP training maxes (daily 2RM - approximately 85-90% of 1RM)
	squatTM := 315.0       // Squat training max
	benchTM := 225.0       // Bench training max
	deadliftTM := 365.0    // Deadlift training max
	ohpTM := 145.0         // OHP training max
	frontSquatTM := 250.0  // Front Squat training max
	sldlTM := 275.0        // Stiff Leg Deadlift training max
	spotoTM := 200.0       // Spoto Bench training max
	pushPressTM := 155.0   // Push Press training max
	lungesWeight := 100.0  // T3 - fixed weight based on 10RM
	pullUpsWeight := 25.0  // T3 - fixed weight (weighted pull-ups)

	// Create training maxes for the user
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatTM)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchTM)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftTM)
	createLiftMax(t, ts, userID, ohpID, "TRAINING_MAX", ohpTM)
	createLiftMax(t, ts, userID, frontSquatID, "TRAINING_MAX", frontSquatTM)
	createLiftMax(t, ts, userID, sldlID, "TRAINING_MAX", sldlTM)
	createLiftMax(t, ts, userID, spotoID, "TRAINING_MAX", spotoTM)
	createLiftMax(t, ts, userID, pushPressID, "TRAINING_MAX", pushPressTM)
	createLiftMax(t, ts, userID, lungesID, "TRAINING_MAX", lungesWeight)
	createLiftMax(t, ts, userID, pullUpsID, "TRAINING_MAX", pullUpsWeight)

	// =============================================================================
	// Create Prescriptions
	// =============================================================================

	// T1 Prescriptions - 3 MRS at 85% TM
	// Using AMRAP set scheme to simulate MRS (Max Rep Sets)
	t1SquatPrescID := createAMRAPPrescription(t, ts, squatID, 3, 5, 85.0, 0)
	t1BenchPrescID := createAMRAPPrescription(t, ts, benchID, 3, 5, 85.0, 0)
	t1DeadliftPrescID := createAMRAPPrescription(t, ts, deadliftID, 3, 5, 85.0, 0)
	t1OhpPrescID := createAMRAPPrescription(t, ts, ohpID, 3, 5, 85.0, 0)
	t1FrontSquatPrescID := createAMRAPPrescription(t, ts, frontSquatID, 3, 5, 85.0, 0)

	// T2 Prescriptions - 3 MRS at 65% TM
	t2SldlPrescID := createAMRAPPrescription(t, ts, sldlID, 3, 8, 65.0, 1)
	t2SpotoPrescID := createAMRAPPrescription(t, ts, spotoID, 3, 8, 65.0, 1)
	t2FrontSquatPrescID := createAMRAPPrescription(t, ts, frontSquatID, 3, 8, 65.0, 1)
	t2PushPressPrescID := createAMRAPPrescription(t, ts, pushPressID, 3, 8, 65.0, 1)
	t2BackSquatPrescID := createAMRAPPrescription(t, ts, squatID, 3, 8, 65.0, 1)

	// T3 Prescriptions - 4 MRS at 100% of fixed weight (10RM weight)
	t3LungesPrescID := createAMRAPPrescription(t, ts, lungesID, 4, 10, 100.0, 2)
	t3PullUpsPrescID := createAMRAPPrescription(t, ts, pullUpsID, 4, 10, 100.0, 3)

	// =============================================================================
	// Create Days - 5 days per week VDIP structure
	// =============================================================================

	// Day 1: Back Squat (T1), Stiff Leg Deadlift (T2), Lunges + Pull-ups (T3)
	day1Slug := "day1-vdip-" + testID
	day1Body := fmt.Sprintf(`{"name": "Day 1 - Squat", "slug": "%s"}`, day1Slug)
	day1Resp, _ := adminPost(ts.URL("/days"), day1Body)
	var day1Envelope DayResponse
	json.NewDecoder(day1Resp.Body).Decode(&day1Envelope)
	day1Resp.Body.Close()
	day1ID := day1Envelope.Data.ID

	addPrescToDay(t, ts, day1ID, t1SquatPrescID)
	addPrescToDay(t, ts, day1ID, t2SldlPrescID)
	addPrescToDay(t, ts, day1ID, t3LungesPrescID)
	addPrescToDay(t, ts, day1ID, t3PullUpsPrescID)

	// Day 2: Bench Press (T1), Spoto Bench (T2), T3 accessories
	day2Slug := "day2-vdip-" + testID
	day2Body := fmt.Sprintf(`{"name": "Day 2 - Bench", "slug": "%s"}`, day2Slug)
	day2Resp, _ := adminPost(ts.URL("/days"), day2Body)
	var day2Envelope DayResponse
	json.NewDecoder(day2Resp.Body).Decode(&day2Envelope)
	day2Resp.Body.Close()
	day2ID := day2Envelope.Data.ID

	addPrescToDay(t, ts, day2ID, t1BenchPrescID)
	addPrescToDay(t, ts, day2ID, t2SpotoPrescID)

	// Day 3: Deadlift (T1), Front Squat (T2), T3 accessories
	day3Slug := "day3-vdip-" + testID
	day3Body := fmt.Sprintf(`{"name": "Day 3 - Deadlift", "slug": "%s"}`, day3Slug)
	day3Resp, _ := adminPost(ts.URL("/days"), day3Body)
	var day3Envelope DayResponse
	json.NewDecoder(day3Resp.Body).Decode(&day3Envelope)
	day3Resp.Body.Close()
	day3ID := day3Envelope.Data.ID

	addPrescToDay(t, ts, day3ID, t1DeadliftPrescID)
	addPrescToDay(t, ts, day3ID, t2FrontSquatPrescID)

	// Day 4: OHP (T1), Push Press (T2), T3 accessories
	day4Slug := "day4-vdip-" + testID
	day4Body := fmt.Sprintf(`{"name": "Day 4 - OHP", "slug": "%s"}`, day4Slug)
	day4Resp, _ := adminPost(ts.URL("/days"), day4Body)
	var day4Envelope DayResponse
	json.NewDecoder(day4Resp.Body).Decode(&day4Envelope)
	day4Resp.Body.Close()
	day4ID := day4Envelope.Data.ID

	addPrescToDay(t, ts, day4ID, t1OhpPrescID)
	addPrescToDay(t, ts, day4ID, t2PushPressPrescID)

	// Day 5: Front Squat (T1), Back Squat (T2), T3 accessories
	day5Slug := "day5-vdip-" + testID
	day5Body := fmt.Sprintf(`{"name": "Day 5 - Front Squat", "slug": "%s"}`, day5Slug)
	day5Resp, _ := adminPost(ts.URL("/days"), day5Body)
	var day5Envelope DayResponse
	json.NewDecoder(day5Resp.Body).Decode(&day5Envelope)
	day5Resp.Body.Close()
	day5ID := day5Envelope.Data.ID

	addPrescToDay(t, ts, day5ID, t1FrontSquatPrescID)
	addPrescToDay(t, ts, day5ID, t2BackSquatPrescID)

	// =============================================================================
	// Create 1-week cycle with 5 training days
	// =============================================================================
	cycleName := "VDIP Cycle " + testID
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

	// Add 5 days to week (Mon-Fri pattern)
	addDayToWeek(t, ts, weekID, day1ID, "MONDAY")
	addDayToWeek(t, ts, weekID, day2ID, "TUESDAY")
	addDayToWeek(t, ts, weekID, day3ID, "WEDNESDAY")
	addDayToWeek(t, ts, weekID, day4ID, "THURSDAY")
	addDayToWeek(t, ts, weekID, day5ID, "FRIDAY")

	// =============================================================================
	// Create Program
	// =============================================================================
	programSlug := "gzcl-vdip-" + testID
	programBody := fmt.Sprintf(`{"name": "GZCL Compendium VDIP", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// Create VDIP Progressions
	// =============================================================================

	// T1 VDIP Progression: Based on total reps in 3 MRS
	// 15+ reps = +10lb, 10-14 reps = +5lb, <10 reps = maintain
	t1ProgBody := `{"name": "VDIP T1 +10/+5lb", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 10.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
	t1ProgResp, _ := adminPost(ts.URL("/progressions"), t1ProgBody)
	var t1ProgEnvelope ProgressionResponse
	json.NewDecoder(t1ProgResp.Body).Decode(&t1ProgEnvelope)
	t1ProgResp.Body.Close()
	t1ProgID := t1ProgEnvelope.Data.ID

	// T2 VDIP Progression: 30+ reps = +10lb, 25-29 reps = +5lb, <25 reps = maintain
	t2ProgBody := `{"name": "VDIP T2 +10/+5lb", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 10.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
	t2ProgResp, _ := adminPost(ts.URL("/progressions"), t2ProgBody)
	var t2ProgEnvelope ProgressionResponse
	json.NewDecoder(t2ProgResp.Body).Decode(&t2ProgEnvelope)
	t2ProgResp.Body.Close()
	t2ProgID := t2ProgEnvelope.Data.ID

	// T3 VDIP Progression: 50+ reps = +5lb, <50 reps = maintain
	t3ProgBody := `{"name": "VDIP T3 +5lb", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
	t3ProgResp, _ := adminPost(ts.URL("/progressions"), t3ProgBody)
	var t3ProgEnvelope ProgressionResponse
	json.NewDecoder(t3ProgResp.Body).Decode(&t3ProgEnvelope)
	t3ProgResp.Body.Close()
	t3ProgID := t3ProgEnvelope.Data.ID

	// Link progressions to program
	// T1 lifts
	linkProgressionToProgram(t, ts, programID, t1ProgID, squatID, 1)
	linkProgressionToProgram(t, ts, programID, t1ProgID, benchID, 2)
	linkProgressionToProgram(t, ts, programID, t1ProgID, deadliftID, 3)
	linkProgressionToProgram(t, ts, programID, t1ProgID, ohpID, 4)
	linkProgressionToProgram(t, ts, programID, t1ProgID, frontSquatID, 5)

	// T2 lifts
	linkProgressionToProgram(t, ts, programID, t2ProgID, sldlID, 6)
	linkProgressionToProgram(t, ts, programID, t2ProgID, spotoID, 7)
	linkProgressionToProgram(t, ts, programID, t2ProgID, pushPressID, 8)

	// T3 lifts
	linkProgressionToProgram(t, ts, programID, t3ProgID, lungesID, 9)
	linkProgressionToProgram(t, ts, programID, t3ProgID, pullUpsID, 10)

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
	// EXECUTION PHASE: Day 1 - Back Squat (T1), SLDL (T2), T3 accessories
	// =============================================================================
	t.Run("Day 1 generates T1 squat MRS at 85% TM", func(t *testing.T) {
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

		// Verify Day 1 structure
		if workout.Data.DaySlug != day1Slug {
			t.Errorf("Expected Day 1 slug '%s', got '%s'", day1Slug, workout.Data.DaySlug)
		}

		// Should have 4 exercises: T1 Squat, T2 SLDL, T3 Lunges, T3 Pull-ups
		if len(workout.Data.Exercises) != 4 {
			t.Fatalf("Expected 4 exercises on Day 1, got %d", len(workout.Data.Exercises))
		}

		// Find T1 squat (first exercise)
		squat := workout.Data.Exercises[0]
		if squat.Lift.ID != squatID {
			t.Errorf("Expected squat lift first, got %s", squat.Lift.ID)
		}

		// T1: 3 MRS at 85% TM
		if len(squat.Sets) != 3 {
			t.Errorf("T1 Squat: expected 3 MRS sets, got %d", len(squat.Sets))
		}

		// Verify weight is 85% of TM (315 * 0.85 = 267.75)
		expectedT1Weight := squatTM * 0.85
		for i, set := range squat.Sets {
			if !withinTolerance(set.Weight, expectedT1Weight, 5.0) {
				t.Errorf("T1 Squat set %d: expected weight ~%.1f, got %.1f", i+1, expectedT1Weight, set.Weight)
			}
		}

		// Verify T2 SLDL (second exercise)
		sldl := workout.Data.Exercises[1]
		if sldl.Lift.ID != sldlID {
			t.Errorf("Expected SLDL lift second, got %s", sldl.Lift.ID)
		}

		// T2: 3 MRS at 65% TM
		if len(sldl.Sets) != 3 {
			t.Errorf("T2 SLDL: expected 3 MRS sets, got %d", len(sldl.Sets))
		}

		expectedT2Weight := sldlTM * 0.65
		for i, set := range sldl.Sets {
			if !withinTolerance(set.Weight, expectedT2Weight, 5.0) {
				t.Errorf("T2 SLDL set %d: expected weight ~%.1f, got %.1f", i+1, expectedT2Weight, set.Weight)
			}
		}

		// Verify T3 accessories (4 MRS each)
		lunges := workout.Data.Exercises[2]
		if lunges.Lift.ID != lungesID {
			t.Errorf("Expected lunges lift third, got %s", lunges.Lift.ID)
		}
		if len(lunges.Sets) != 4 {
			t.Errorf("T3 Lunges: expected 4 MRS sets, got %d", len(lunges.Sets))
		}

		pullUps := workout.Data.Exercises[3]
		if pullUps.Lift.ID != pullUpsID {
			t.Errorf("Expected pull-ups lift fourth, got %s", pullUps.Lift.ID)
		}
		if len(pullUps.Sets) != 4 {
			t.Errorf("T3 Pull-ups: expected 4 MRS sets, got %d", len(pullUps.Sets))
		}
	})

	// Trigger progression for Day 1 squat and complete workout
	triggerProgressionForLift(t, ts, userID, t1ProgID, squatID)
	completeVDIPWorkoutDay(t, ts, userID)

	// =============================================================================
	// EXECUTION PHASE: Day 2 - Bench Press (T1), Spoto Bench (T2)
	// =============================================================================
	t.Run("Day 2 generates T1 bench MRS at 85% TM", func(t *testing.T) {
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

		// Verify Day 2 structure
		if workout.Data.DaySlug != day2Slug {
			t.Errorf("Expected Day 2 slug '%s', got '%s'", day2Slug, workout.Data.DaySlug)
		}

		// Should have 2 exercises: T1 Bench, T2 Spoto
		if len(workout.Data.Exercises) != 2 {
			t.Fatalf("Expected 2 exercises on Day 2, got %d", len(workout.Data.Exercises))
		}

		// T1 Bench: 3 MRS at 85% TM
		bench := workout.Data.Exercises[0]
		if bench.Lift.ID != benchID {
			t.Errorf("Expected bench lift first, got %s", bench.Lift.ID)
		}
		if len(bench.Sets) != 3 {
			t.Errorf("T1 Bench: expected 3 MRS sets, got %d", len(bench.Sets))
		}

		expectedBenchWeight := benchTM * 0.85
		for i, set := range bench.Sets {
			if !withinTolerance(set.Weight, expectedBenchWeight, 5.0) {
				t.Errorf("T1 Bench set %d: expected weight ~%.1f, got %.1f", i+1, expectedBenchWeight, set.Weight)
			}
		}

		// T2 Spoto: 3 MRS at 65% TM
		spoto := workout.Data.Exercises[1]
		if spoto.Lift.ID != spotoID {
			t.Errorf("Expected Spoto bench lift second, got %s", spoto.Lift.ID)
		}
		if len(spoto.Sets) != 3 {
			t.Errorf("T2 Spoto: expected 3 MRS sets, got %d", len(spoto.Sets))
		}

		expectedSpotoWeight := spotoTM * 0.65
		for i, set := range spoto.Sets {
			if !withinTolerance(set.Weight, expectedSpotoWeight, 5.0) {
				t.Errorf("T2 Spoto set %d: expected weight ~%.1f, got %.1f", i+1, expectedSpotoWeight, set.Weight)
			}
		}
	})

	// Complete remaining days using explicit session lifecycle
	completeVDIPWorkoutDay(t, ts, userID) // Day 2 -> Day 3
	completeVDIPWorkoutDay(t, ts, userID) // Day 3 -> Day 4
	completeVDIPWorkoutDay(t, ts, userID) // Day 4 -> Day 5

	// =============================================================================
	// EXECUTION PHASE: Day 5 - Front Squat (T1), Back Squat (T2)
	// =============================================================================
	t.Run("Day 5 generates T1 front squat and T2 back squat", func(t *testing.T) {
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

		// Verify Day 5 structure
		if workout.Data.DaySlug != day5Slug {
			t.Errorf("Expected Day 5 slug '%s', got '%s'", day5Slug, workout.Data.DaySlug)
		}

		// Should have 2 exercises: T1 Front Squat, T2 Back Squat
		if len(workout.Data.Exercises) != 2 {
			t.Fatalf("Expected 2 exercises on Day 5, got %d", len(workout.Data.Exercises))
		}

		// T1 Front Squat: 3 MRS at 85% TM
		frontSquat := workout.Data.Exercises[0]
		if frontSquat.Lift.ID != frontSquatID {
			t.Errorf("Expected front squat lift first, got %s", frontSquat.Lift.ID)
		}
		if len(frontSquat.Sets) != 3 {
			t.Errorf("T1 Front Squat: expected 3 MRS sets, got %d", len(frontSquat.Sets))
		}

		expectedFrontSquatWeight := frontSquatTM * 0.85
		for i, set := range frontSquat.Sets {
			if !withinTolerance(set.Weight, expectedFrontSquatWeight, 5.0) {
				t.Errorf("T1 Front Squat set %d: expected weight ~%.1f, got %.1f", i+1, expectedFrontSquatWeight, set.Weight)
			}
		}

		// T2 Back Squat: 3 MRS at 65% TM
		// Note: Squat TM was increased by +10 from Day 1 progression
		backSquat := workout.Data.Exercises[1]
		if backSquat.Lift.ID != squatID {
			t.Errorf("Expected back squat lift second, got %s", backSquat.Lift.ID)
		}
		if len(backSquat.Sets) != 3 {
			t.Errorf("T2 Back Squat: expected 3 MRS sets, got %d", len(backSquat.Sets))
		}

		// Back Squat TM should be updated (315 + 10 = 325) from Day 1 T1 progression
		expectedBackSquatTM := squatTM + 10.0
		expectedBackSquatWeight := expectedBackSquatTM * 0.65
		for i, set := range backSquat.Sets {
			if !withinTolerance(set.Weight, expectedBackSquatWeight, 5.0) {
				t.Errorf("T2 Back Squat set %d: expected weight ~%.1f (65%% of updated TM), got %.1f", i+1, expectedBackSquatWeight, set.Weight)
			}
		}
	})

	// =============================================================================
	// VDIP PROGRESSION TEST: Verify progression applies correctly
	// =============================================================================
	t.Run("VDIP T1 progression applies +10lb on 15+ total reps", func(t *testing.T) {
		// Trigger T1 progression for bench (simulating 15+ total reps achieved)
		benchTrigger := triggerProgressionForLift(t, ts, userID, t1ProgID, benchID)

		if benchTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected bench progression to apply, got TotalApplied=%d", benchTrigger.Data.TotalApplied)
		}

		if len(benchTrigger.Data.Results) > 0 && benchTrigger.Data.Results[0].Result != nil {
			if benchTrigger.Data.Results[0].Result.Delta != 10.0 {
				t.Errorf("Expected bench delta +10, got %f", benchTrigger.Data.Results[0].Result.Delta)
			}
			expectedNewBench := benchTM + 10.0 // 225 + 10 = 235
			if benchTrigger.Data.Results[0].Result.NewValue != expectedNewBench {
				t.Errorf("Expected bench new value %f, got %f", expectedNewBench, benchTrigger.Data.Results[0].Result.NewValue)
			}
		}
	})

	t.Run("VDIP T3 progression applies +5lb on 50+ total reps", func(t *testing.T) {
		// Trigger T3 progression for lunges (simulating 50+ total reps achieved)
		lungesTrigger := triggerProgressionForLift(t, ts, userID, t3ProgID, lungesID)

		if lungesTrigger.Data.TotalApplied != 1 {
			t.Errorf("Expected lunges progression to apply, got TotalApplied=%d", lungesTrigger.Data.TotalApplied)
		}

		if len(lungesTrigger.Data.Results) > 0 && lungesTrigger.Data.Results[0].Result != nil {
			if lungesTrigger.Data.Results[0].Result.Delta != 5.0 {
				t.Errorf("Expected lunges delta +5, got %f", lungesTrigger.Data.Results[0].Result.Delta)
			}
			expectedNewLunges := lungesWeight + 5.0 // 100 + 5 = 105
			if lungesTrigger.Data.Results[0].Result.NewValue != expectedNewLunges {
				t.Errorf("Expected lunges new value %f, got %f", expectedNewLunges, lungesTrigger.Data.Results[0].Result.NewValue)
			}
		}
	})

	// =============================================================================
	// WEEK 2: Verify new week starts with updated training maxes
	// =============================================================================
	completeVDIPWorkoutDay(t, ts, userID) // Day 5 -> Day 1 (Week 2)

	t.Run("Week 2 Day 1 shows updated training maxes", func(t *testing.T) {
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

		// Should be back to Day 1
		if workout.Data.DaySlug != day1Slug {
			t.Errorf("Expected Day 1 slug '%s', got '%s'", day1Slug, workout.Data.DaySlug)
		}

		// Verify squat uses updated TM (315 + 10 = 325)
		squat := workout.Data.Exercises[0]
		expectedSquatTM := squatTM + 10.0
		expectedSquatWeight := expectedSquatTM * 0.85 // 325 * 0.85 = 276.25

		if !withinTolerance(squat.Sets[0].Weight, expectedSquatWeight, 5.0) {
			t.Errorf("Week 2 T1 Squat: expected weight ~%.1f (85%% of updated TM), got %.1f", expectedSquatWeight, squat.Sets[0].Weight)
		}
	})
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// completeVDIPWorkoutDay completes a GZCL VDIP workout day using explicit state machine flow.
func completeVDIPWorkoutDay(t *testing.T, ts *testutil.TestServer, userID string) {
	t.Helper()

	sessionID := startWorkoutSession(t, ts, userID)

	workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
	var workout WorkoutResponse
	json.NewDecoder(workoutResp.Body).Decode(&workout)
	workoutResp.Body.Close()

	for _, ex := range workout.Data.Exercises {
		for _, set := range ex.Sets {
			logVDIPSet(t, ts, userID, sessionID, ex.PrescriptionID, ex.Lift.ID, set.SetNumber, set.Weight, set.TargetReps, set.TargetReps, false)
		}
	}

	finishWorkoutSession(t, ts, sessionID, userID)
	advanceUserState(t, ts, userID)
}

// logVDIPSet logs a single set for GZCL VDIP workout.
func logVDIPSet(t *testing.T, ts *testutil.TestServer, userID, sessionID, prescriptionID, liftID string, setNumber int, weight float64, targetReps, repsPerformed int, isAmrap bool) {
	t.Helper()

	loggedSetBody := fmt.Sprintf(`{
		"sets": [{
			"prescriptionId": "%s",
			"liftId": "%s",
			"setNumber": %d,
			"weight": %.1f,
			"targetReps": %d,
			"repsPerformed": %d,
			"isAmrap": %t
		}]
	}`, prescriptionID, liftID, setNumber, weight, targetReps, repsPerformed, isAmrap)

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
