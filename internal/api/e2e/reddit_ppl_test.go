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
// REDDIT PPL 6-DAY E2E TEST
// =============================================================================

// TestRedditPPL6DayProgram validates the complete Reddit PPL 6-day program
// configuration and execution through the API.
//
// Reddit PPL characteristics:
// - 6-day split: Pull/Push/Legs x 2 per week (Pull A -> Push A -> Legs A -> Pull B -> Push B -> Legs B)
// - Alternating Primary Lifts: Different primary compounds on A vs B days
// - AMRAP prescriptions: Deadlift 1x5+, Bench 4x5+1x5+, Squat 2x5+1x5+, Rows 4x5+1x5+, OHP 4x5+1x5+
// - Linear Progression: +2.5lb upper body, +5lb lower body per session
func TestRedditPPL6DayProgram(t *testing.T) {
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

	// Create additional lifts (OHP and Rows are not seeded)
	ohpSlug := "overhead-press-" + testID
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

	rowsSlug := "barbell-rows-" + testID
	rowsBody := fmt.Sprintf(`{"name": "Barbell Rows", "slug": "%s", "isCompetitionLift": false}`, rowsSlug)
	rowsResp, err := adminPost(ts.URL("/lifts"), rowsBody)
	if err != nil {
		t.Fatalf("Failed to create Rows lift: %v", err)
	}
	if rowsResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(rowsResp.Body)
		rowsResp.Body.Close()
		t.Fatalf("Failed to create Rows lift, status %d: %s", rowsResp.StatusCode, body)
	}
	var rowsEnvelope LiftResponse
	json.NewDecoder(rowsResp.Body).Decode(&rowsEnvelope)
	rowsResp.Body.Close()
	rowsID := rowsEnvelope.Data.ID

	// Reddit PPL training maxes
	squatTM := 225.0    // Squat training max
	benchTM := 185.0    // Bench training max
	deadliftTM := 275.0 // Deadlift training max
	ohpTM := 115.0      // OHP training max
	rowsTM := 155.0     // Rows training max

	// Create training maxes for the user
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatTM)
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchTM)
	createLiftMax(t, ts, userID, deadliftID, "TRAINING_MAX", deadliftTM)
	createLiftMax(t, ts, userID, ohpID, "TRAINING_MAX", ohpTM)
	createLiftMax(t, ts, userID, rowsID, "TRAINING_MAX", rowsTM)

	// =============================================================================
	// Create Prescriptions
	// =============================================================================

	// Deadlift: 1x5+ (pure AMRAP)
	deadliftPrescID := createAMRAPPrescription(t, ts, deadliftID, 1, 5, 100.0, 0)

	// Barbell Rows (Pull B): 4x5, 1x5+ (GREYSKULL scheme)
	rowsPrescID := createGreyskullPrescription(t, ts, rowsID, 4, 5, 1, 5, 100.0, 0)

	// Bench Press (Push A): 4x5, 1x5+ (GREYSKULL scheme)
	benchPrescID := createGreyskullPrescription(t, ts, benchID, 4, 5, 1, 5, 100.0, 0)

	// OHP (Push B): 4x5, 1x5+ (GREYSKULL scheme)
	ohpPrescID := createGreyskullPrescription(t, ts, ohpID, 4, 5, 1, 5, 100.0, 0)

	// Squat (Legs A & B): 2x5, 1x5+ (GREYSKULL scheme)
	squatPrescID := createGreyskullPrescription(t, ts, squatID, 2, 5, 1, 5, 100.0, 0)

	// =============================================================================
	// Create Days
	// =============================================================================

	// Pull Day A: Deadlifts 1x5+
	pullASlug := "pull-a-" + testID
	pullABody := fmt.Sprintf(`{"name": "Pull Day A", "slug": "%s"}`, pullASlug)
	pullAResp, _ := adminPost(ts.URL("/days"), pullABody)
	var pullAEnvelope DayResponse
	json.NewDecoder(pullAResp.Body).Decode(&pullAEnvelope)
	pullAResp.Body.Close()
	pullAID := pullAEnvelope.Data.ID
	addPrescToDay(t, ts, pullAID, deadliftPrescID)

	// Push Day A: Bench 4x5, 1x5+
	pushASlug := "push-a-" + testID
	pushABody := fmt.Sprintf(`{"name": "Push Day A", "slug": "%s"}`, pushASlug)
	pushAResp, _ := adminPost(ts.URL("/days"), pushABody)
	var pushAEnvelope DayResponse
	json.NewDecoder(pushAResp.Body).Decode(&pushAEnvelope)
	pushAResp.Body.Close()
	pushAID := pushAEnvelope.Data.ID
	addPrescToDay(t, ts, pushAID, benchPrescID)

	// Legs Day A: Squat 2x5, 1x5+
	legsASlug := "legs-a-" + testID
	legsABody := fmt.Sprintf(`{"name": "Legs Day A", "slug": "%s"}`, legsASlug)
	legsAResp, _ := adminPost(ts.URL("/days"), legsABody)
	var legsAEnvelope DayResponse
	json.NewDecoder(legsAResp.Body).Decode(&legsAEnvelope)
	legsAResp.Body.Close()
	legsAID := legsAEnvelope.Data.ID
	addPrescToDay(t, ts, legsAID, squatPrescID)

	// Pull Day B: Barbell Rows 4x5, 1x5+
	pullBSlug := "pull-b-" + testID
	pullBBody := fmt.Sprintf(`{"name": "Pull Day B", "slug": "%s"}`, pullBSlug)
	pullBResp, _ := adminPost(ts.URL("/days"), pullBBody)
	var pullBEnvelope DayResponse
	json.NewDecoder(pullBResp.Body).Decode(&pullBEnvelope)
	pullBResp.Body.Close()
	pullBID := pullBEnvelope.Data.ID
	addPrescToDay(t, ts, pullBID, rowsPrescID)

	// Push Day B: OHP 4x5, 1x5+
	pushBSlug := "push-b-" + testID
	pushBBody := fmt.Sprintf(`{"name": "Push Day B", "slug": "%s"}`, pushBSlug)
	pushBResp, _ := adminPost(ts.URL("/days"), pushBBody)
	var pushBEnvelope DayResponse
	json.NewDecoder(pushBResp.Body).Decode(&pushBEnvelope)
	pushBResp.Body.Close()
	pushBID := pushBEnvelope.Data.ID
	addPrescToDay(t, ts, pushBID, ohpPrescID)

	// Legs Day B: Squat 2x5, 1x5+ (separate prescription for independent tracking)
	squatPrescBID := createGreyskullPrescription(t, ts, squatID, 2, 5, 1, 5, 100.0, 0)
	legsBSlug := "legs-b-" + testID
	legsBBody := fmt.Sprintf(`{"name": "Legs Day B", "slug": "%s"}`, legsBSlug)
	legsBResp, _ := adminPost(ts.URL("/days"), legsBBody)
	var legsBEnvelope DayResponse
	json.NewDecoder(legsBResp.Body).Decode(&legsBEnvelope)
	legsBResp.Body.Close()
	legsBID := legsBEnvelope.Data.ID
	addPrescToDay(t, ts, legsBID, squatPrescBID)

	// =============================================================================
	// Create 1-week cycle with 6 training days
	// =============================================================================
	cycleName := "Reddit PPL Cycle " + testID
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

	// Add days to week: Pull A -> Push A -> Legs A -> Pull B -> Push B -> Legs B
	addDayToWeek(t, ts, weekID, pullAID, "MONDAY")    // Day 1: Pull A
	addDayToWeek(t, ts, weekID, pushAID, "TUESDAY")   // Day 2: Push A
	addDayToWeek(t, ts, weekID, legsAID, "WEDNESDAY") // Day 3: Legs A
	addDayToWeek(t, ts, weekID, pullBID, "THURSDAY")  // Day 4: Pull B
	addDayToWeek(t, ts, weekID, pushBID, "FRIDAY")    // Day 5: Push B
	addDayToWeek(t, ts, weekID, legsBID, "SATURDAY")  // Day 6: Legs B

	// =============================================================================
	// Create Program
	// =============================================================================
	programSlug := "reddit-ppl-6-day-" + testID
	programBody := fmt.Sprintf(`{"name": "Reddit PPL 6-Day", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// =============================================================================
	// Create Linear Progressions
	// =============================================================================

	// Upper body progression (+2.5lb)
	upperProgBody := `{"name": "PPL Upper +2.5lb", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 2.5, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
	upperProgResp, _ := adminPost(ts.URL("/progressions"), upperProgBody)
	var upperProgEnvelope ProgressionResponse
	json.NewDecoder(upperProgResp.Body).Decode(&upperProgEnvelope)
	upperProgResp.Body.Close()
	upperProgID := upperProgEnvelope.Data.ID

	// Lower body progression (+5lb)
	lowerProgBody := `{"name": "PPL Lower +5lb", "type": "LINEAR_PROGRESSION", "parameters": {"increment": 5.0, "maxType": "TRAINING_MAX", "triggerType": "AFTER_SESSION"}}`
	lowerProgResp, _ := adminPost(ts.URL("/progressions"), lowerProgBody)
	var lowerProgEnvelope ProgressionResponse
	json.NewDecoder(lowerProgResp.Body).Decode(&lowerProgEnvelope)
	lowerProgResp.Body.Close()
	lowerProgID := lowerProgEnvelope.Data.ID

	// Link progressions to program
	// Lower body lifts get +5lb
	linkProgressionToProgram(t, ts, programID, lowerProgID, squatID, 1)
	linkProgressionToProgram(t, ts, programID, lowerProgID, deadliftID, 2)
	// Upper body lifts get +2.5lb
	linkProgressionToProgram(t, ts, programID, upperProgID, benchID, 3)
	linkProgressionToProgram(t, ts, programID, upperProgID, ohpID, 4)
	linkProgressionToProgram(t, ts, programID, upperProgID, rowsID, 5)

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
	// EXECUTION PHASE: Day 1 - Pull A (Deadlifts 1x5+)
	// =============================================================================
	t.Run("Day 1 Pull A generates deadlift AMRAP set", func(t *testing.T) {
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

		// Verify Pull A structure
		if workout.Data.DaySlug != pullASlug {
			t.Errorf("Expected Pull A slug '%s', got '%s'", pullASlug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) != 1 {
			t.Fatalf("Expected 1 exercise on Pull A, got %d", len(workout.Data.Exercises))
		}

		// Verify deadlift AMRAP
		deadlift := workout.Data.Exercises[0]
		if deadlift.Lift.ID != deadliftID {
			t.Errorf("Expected deadlift lift, got %s", deadlift.Lift.ID)
		}

		// Deadlift: 1x5+ (1 AMRAP set)
		if len(deadlift.Sets) != 1 {
			t.Errorf("Deadlift: expected 1 AMRAP set, got %d sets", len(deadlift.Sets))
		}

		if len(deadlift.Sets) > 0 {
			if deadlift.Sets[0].Weight != deadliftTM {
				t.Errorf("Deadlift: expected weight %f, got %f", deadliftTM, deadlift.Sets[0].Weight)
			}
			if deadlift.Sets[0].TargetReps != 5 {
				t.Errorf("Deadlift: expected target reps 5, got %d", deadlift.Sets[0].TargetReps)
			}
		}
	})

	// Trigger progression for deadlift and advance to Push A
	triggerBody := ManualTriggerRequest{
		ProgressionID: lowerProgID,
		LiftID:        deadliftID,
		Force:         true,
	}
	triggerResp, _ := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
	triggerResp.Body.Close()
	advanceUserState(t, ts, userID)

	// =============================================================================
	// EXECUTION PHASE: Day 2 - Push A (Bench 4x5, 1x5+)
	// =============================================================================
	t.Run("Day 2 Push A generates bench 4x5+1x5+ sets", func(t *testing.T) {
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

		// Verify Push A structure
		if workout.Data.DaySlug != pushASlug {
			t.Errorf("Expected Push A slug '%s', got '%s'", pushASlug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) != 1 {
			t.Fatalf("Expected 1 exercise on Push A, got %d", len(workout.Data.Exercises))
		}

		// Verify bench press sets
		bench := workout.Data.Exercises[0]
		if bench.Lift.ID != benchID {
			t.Errorf("Expected bench lift, got %s", bench.Lift.ID)
		}

		// Bench: 4x5, 1x5+ (5 total sets)
		if len(bench.Sets) != 5 {
			t.Errorf("Bench: expected 5 sets (4x5 + 1x5+), got %d sets", len(bench.Sets))
		}

		// All sets should be at bench training max
		for i, set := range bench.Sets {
			if set.Weight != benchTM {
				t.Errorf("Bench set %d: expected weight %f, got %f", i+1, benchTM, set.Weight)
			}
			if set.TargetReps != 5 {
				t.Errorf("Bench set %d: expected target reps 5, got %d", i+1, set.TargetReps)
			}
		}
	})

	// Trigger progression for bench and advance to Legs A
	triggerBody.ProgressionID = upperProgID
	triggerBody.LiftID = benchID
	triggerResp, _ = authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
	triggerResp.Body.Close()
	advanceUserState(t, ts, userID)

	// =============================================================================
	// EXECUTION PHASE: Day 3 - Legs A (Squat 2x5, 1x5+)
	// =============================================================================
	t.Run("Day 3 Legs A generates squat 2x5+1x5+ sets", func(t *testing.T) {
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

		// Verify Legs A structure
		if workout.Data.DaySlug != legsASlug {
			t.Errorf("Expected Legs A slug '%s', got '%s'", legsASlug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) != 1 {
			t.Fatalf("Expected 1 exercise on Legs A, got %d", len(workout.Data.Exercises))
		}

		// Verify squat sets
		squat := workout.Data.Exercises[0]
		if squat.Lift.ID != squatID {
			t.Errorf("Expected squat lift, got %s", squat.Lift.ID)
		}

		// Squat: 2x5, 1x5+ (3 total sets)
		if len(squat.Sets) != 3 {
			t.Errorf("Squat: expected 3 sets (2x5 + 1x5+), got %d sets", len(squat.Sets))
		}

		// All sets should be at squat training max
		for i, set := range squat.Sets {
			if set.Weight != squatTM {
				t.Errorf("Squat set %d: expected weight %f, got %f", i+1, squatTM, set.Weight)
			}
			if set.TargetReps != 5 {
				t.Errorf("Squat set %d: expected target reps 5, got %d", i+1, set.TargetReps)
			}
		}
	})

	// Trigger progression for squat and advance to Pull B
	triggerBody.ProgressionID = lowerProgID
	triggerBody.LiftID = squatID
	triggerResp, _ = authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
	triggerResp.Body.Close()
	advanceUserState(t, ts, userID)

	// =============================================================================
	// EXECUTION PHASE: Day 4 - Pull B (Barbell Rows 4x5, 1x5+)
	// =============================================================================
	t.Run("Day 4 Pull B generates rows 4x5+1x5+ sets", func(t *testing.T) {
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

		// Verify Pull B structure
		if workout.Data.DaySlug != pullBSlug {
			t.Errorf("Expected Pull B slug '%s', got '%s'", pullBSlug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) != 1 {
			t.Fatalf("Expected 1 exercise on Pull B, got %d", len(workout.Data.Exercises))
		}

		// Verify rows sets
		rows := workout.Data.Exercises[0]
		if rows.Lift.ID != rowsID {
			t.Errorf("Expected rows lift, got %s", rows.Lift.ID)
		}

		// Rows: 4x5, 1x5+ (5 total sets)
		if len(rows.Sets) != 5 {
			t.Errorf("Rows: expected 5 sets (4x5 + 1x5+), got %d sets", len(rows.Sets))
		}

		// All sets should be at rows training max
		for i, set := range rows.Sets {
			if set.Weight != rowsTM {
				t.Errorf("Rows set %d: expected weight %f, got %f", i+1, rowsTM, set.Weight)
			}
			if set.TargetReps != 5 {
				t.Errorf("Rows set %d: expected target reps 5, got %d", i+1, set.TargetReps)
			}
		}
	})

	// Trigger progression for rows and advance to Push B
	triggerBody.ProgressionID = upperProgID
	triggerBody.LiftID = rowsID
	triggerResp, _ = authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
	triggerResp.Body.Close()
	advanceUserState(t, ts, userID)

	// =============================================================================
	// EXECUTION PHASE: Day 5 - Push B (OHP 4x5, 1x5+)
	// =============================================================================
	t.Run("Day 5 Push B generates OHP 4x5+1x5+ sets", func(t *testing.T) {
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

		// Verify Push B structure
		if workout.Data.DaySlug != pushBSlug {
			t.Errorf("Expected Push B slug '%s', got '%s'", pushBSlug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) != 1 {
			t.Fatalf("Expected 1 exercise on Push B, got %d", len(workout.Data.Exercises))
		}

		// Verify OHP sets
		ohp := workout.Data.Exercises[0]
		if ohp.Lift.ID != ohpID {
			t.Errorf("Expected OHP lift, got %s", ohp.Lift.ID)
		}

		// OHP: 4x5, 1x5+ (5 total sets)
		if len(ohp.Sets) != 5 {
			t.Errorf("OHP: expected 5 sets (4x5 + 1x5+), got %d sets", len(ohp.Sets))
		}

		// All sets should be at OHP training max
		for i, set := range ohp.Sets {
			if set.Weight != ohpTM {
				t.Errorf("OHP set %d: expected weight %f, got %f", i+1, ohpTM, set.Weight)
			}
			if set.TargetReps != 5 {
				t.Errorf("OHP set %d: expected target reps 5, got %d", i+1, set.TargetReps)
			}
		}
	})

	// Trigger progression for OHP and advance to Legs B
	triggerBody.LiftID = ohpID
	triggerResp, _ = authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
	triggerResp.Body.Close()
	advanceUserState(t, ts, userID)

	// =============================================================================
	// EXECUTION PHASE: Day 6 - Legs B (Squat 2x5, 1x5+)
	// =============================================================================
	t.Run("Day 6 Legs B shows updated squat weight after progression", func(t *testing.T) {
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

		// Verify Legs B structure
		if workout.Data.DaySlug != legsBSlug {
			t.Errorf("Expected Legs B slug '%s', got '%s'", legsBSlug, workout.Data.DaySlug)
		}

		if len(workout.Data.Exercises) != 1 {
			t.Fatalf("Expected 1 exercise on Legs B, got %d", len(workout.Data.Exercises))
		}

		// Verify squat sets with progressed weight
		squat := workout.Data.Exercises[0]
		if squat.Lift.ID != squatID {
			t.Errorf("Expected squat lift, got %s", squat.Lift.ID)
		}

		// Squat: 2x5, 1x5+ (3 total sets)
		if len(squat.Sets) != 3 {
			t.Errorf("Squat: expected 3 sets (2x5 + 1x5+), got %d sets", len(squat.Sets))
		}

		// Squat should now be at 230 (225 + 5 from Legs A progression)
		expectedSquat := squatTM + 5.0
		for i, set := range squat.Sets {
			if set.Weight != expectedSquat {
				t.Errorf("Squat set %d: expected weight %f (progressed from Legs A), got %f", i+1, expectedSquat, set.Weight)
			}
		}
	})

	// =============================================================================
	// VALIDATION PHASE: Full cycle complete - verify all accumulated progressions
	// =============================================================================
	t.Run("All progressions applied correctly after full 6-day cycle", func(t *testing.T) {
		// Trigger squat progression for Legs B
		triggerBody.ProgressionID = lowerProgID
		triggerBody.LiftID = squatID
		triggerResp, _ = authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
		triggerResp.Body.Close()

		// Advance to next week's Pull A
		advanceUserState(t, ts, userID)

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

		// Should be back to Pull A
		if workout.Data.DaySlug != pullASlug {
			t.Errorf("Expected Pull A slug '%s', got '%s'", pullASlug, workout.Data.DaySlug)
		}

		// Verify deadlift weight increased (+5lb from Day 1)
		expectedDeadlift := deadliftTM + 5.0
		deadlift := workout.Data.Exercises[0]
		if deadlift.Sets[0].Weight != expectedDeadlift {
			t.Errorf("Deadlift: expected weight %f (progressed), got %f", expectedDeadlift, deadlift.Sets[0].Weight)
		}
	})
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// createGreyskullPrescription creates a prescription with GREYSKULL set scheme (NxR + 1xAMRAP).
func createGreyskullPrescription(t *testing.T, ts *testutil.TestServer, liftID string, fixedSets, fixedReps, amrapSets, minAmrapReps int, percentage float64, order int) string {
	t.Helper()

	body := fmt.Sprintf(`{
		"liftId": "%s",
		"loadStrategy": {"type": "PERCENT_OF", "referenceType": "TRAINING_MAX", "percentage": %.1f},
		"setScheme": {"type": "GREYSKULL", "fixedSets": %d, "fixedReps": %d, "amrapSets": %d, "minAmrapReps": %d},
		"order": %d
	}`, liftID, percentage, fixedSets, fixedReps, amrapSets, minAmrapReps, order)

	resp, err := adminPost(ts.URL("/prescriptions"), body)
	if err != nil {
		t.Fatalf("Failed to create GREYSKULL prescription: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to create GREYSKULL prescription, status %d: %s", resp.StatusCode, bodyBytes)
	}

	var envelope PrescriptionResponse
	json.NewDecoder(resp.Body).Decode(&envelope)
	return envelope.Data.ID
}
