// Package e2e provides end-to-end tests for complete program workflows.
// This file tests GZCLP T2 stage-based progression WITHOUT AMRAP sets.
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

// TestGZCLPT2StageProgression validates the GZCLP T2 stage progression system.
// T2 lifts use stage progression WITHOUT AMRAP:
//
//	3x10 -> 3x8 -> 3x6 with no deload on reset
//
// Key behaviors tested:
//  1. On failure (not hitting total volume), advance to next stage
//  2. On last stage failure with ResetOnExhaustion=true but DeloadOnReset=false,
//     reset to stage 0 without weight change (requires manual weight reduction)
//  3. Success resets the failure counter
func TestGZCLPT2StageProgression(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user" // Uses seeded test user
	_ = testID                     // Keep for unique slugs

	// Seeded bench lift
	benchID := "00000000-0000-0000-0000-000000000002"
	benchMax := 100.0

	// Create training max
	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchMax)

	// Create GZCLP T2 Stage Progression (3x10 -> 3x8 -> 3x6)
	// Note: No AMRAP, no deload on reset
	stageProgParams := StageProgressionParams{
		Stages: []StageParam{
			{Name: "3x10", Sets: 3, Reps: 10, IsAMRAP: false, MinVolume: 30},
			{Name: "3x8", Sets: 3, Reps: 8, IsAMRAP: false, MinVolume: 24},
			{Name: "3x6", Sets: 3, Reps: 6, IsAMRAP: false, MinVolume: 18},
		},
		CurrentStage:      0,
		ResetOnExhaustion: true,
		DeloadOnReset:     false, // T2 doesn't auto-deload, requires manual intervention
		DeloadPercent:     0,
		MaxType:           "TRAINING_MAX",
	}

	stageProgParamsJSON, _ := json.Marshal(stageProgParams)
	stageProgBody := fmt.Sprintf(`{"name": "GZCLP T2 Bench", "type": "STAGE_PROGRESSION", "parameters": %s}`, stageProgParamsJSON)
	stageProgResp, err := adminPost(ts.URL("/progressions"), stageProgBody)
	if err != nil {
		t.Fatalf("Failed to create stage progression: %v", err)
	}
	if stageProgResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(stageProgResp.Body)
		stageProgResp.Body.Close()
		t.Fatalf("Failed to create stage progression, status %d: %s", stageProgResp.StatusCode, body)
	}
	var stageProgEnvelope ProgressionResponse
	json.NewDecoder(stageProgResp.Body).Decode(&stageProgEnvelope)
	stageProgResp.Body.Close()
	stageProgID := stageProgEnvelope.Data.ID

	// Create prescription for bench 3x10 at 100% TM (no AMRAP)
	benchPrescID := createPrescription(t, ts, benchID, 3, 10, 100.0, 0)

	// Create day
	daySlug := "gzclp-t2-day-" + testID
	dayBody := fmt.Sprintf(`{"name": "GZCLP T2 Day", "slug": "%s"}`, daySlug)
	dayResp, _ := adminPost(ts.URL("/days"), dayBody)
	var dayEnvelope DayResponse
	json.NewDecoder(dayResp.Body).Decode(&dayEnvelope)
	dayResp.Body.Close()
	dayID := dayEnvelope.Data.ID
	addPrescToDay(t, ts, dayID, benchPrescID)

	// Create cycle/week
	cycleName := "GZCLP T2 Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 1}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	weekBody := fmt.Sprintf(`{"weekNumber": 1, "cycleId": "%s"}`, cycleID)
	weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
	var weekEnvelope WeekResponse
	json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
	weekResp.Body.Close()
	weekID := weekEnvelope.Data.ID
	addDayToWeek(t, ts, weekID, dayID, "MONDAY")

	// Create program
	programSlug := "gzclp-t2-" + testID
	programBody := fmt.Sprintf(`{"name": "GZCLP T2", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	linkProgressionToProgram(t, ts, programID, stageProgID, benchID, 1)

	// Enroll user
	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Test 1: Verify initial workout shows 3x10
	t.Run("initial workout shows 3x10 at 100 lbs", func(t *testing.T) {
		workoutResp, _ := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		if len(workout.Data.Exercises) != 1 {
			t.Fatalf("Expected 1 exercise, got %d", len(workout.Data.Exercises))
		}

		ex := workout.Data.Exercises[0]
		if len(ex.Sets) != 3 {
			t.Errorf("Expected 3 sets, got %d", len(ex.Sets))
		}
		for i, set := range ex.Sets {
			if set.Weight != benchMax {
				t.Errorf("Set %d: expected weight %.1f, got %.1f", i+1, benchMax, set.Weight)
			}
			if set.TargetReps != 10 {
				t.Errorf("Set %d: expected 10 reps, got %d", i+1, set.TargetReps)
			}
		}
	})

	// Test 2: Stage 0 (3x10) SUCCESS - failure counter should reset
	t.Run("success at stage 0 keeps counter at 0", func(t *testing.T) {
		sessionID := uuid.New().String()
		// Hit all reps: 30 total >= 30 minimum
		logSets(t, ts, userID, sessionID, benchPrescID, benchID, []setLog{
			{weight: benchMax, targetReps: 10, repsPerformed: 10, isAMRAP: false},
			{weight: benchMax, targetReps: 10, repsPerformed: 10, isAMRAP: false},
			{weight: benchMax, targetReps: 10, repsPerformed: 10, isAMRAP: false},
		})

		count := getFailureCount(t, ts, userID, benchID, stageProgID)
		if count != 0 {
			t.Errorf("Expected failure count 0 after success, got %d", count)
		}
	})

	advanceUserState(t, ts, userID)

	// Test 3: Stage 0 (3x10) FAIL - advance to stage 1
	t.Run("failure at stage 0 advances to stage 1", func(t *testing.T) {
		sessionID := uuid.New().String()
		// Only 27 reps < 30 minimum
		logSets(t, ts, userID, sessionID, benchPrescID, benchID, []setLog{
			{weight: benchMax, targetReps: 10, repsPerformed: 10, isAMRAP: false},
			{weight: benchMax, targetReps: 10, repsPerformed: 10, isAMRAP: false},
			{weight: benchMax, targetReps: 10, repsPerformed: 7, isAMRAP: false}, // Failed last set
		})

		count := getFailureCount(t, ts, userID, benchID, stageProgID)
		if count != 1 {
			t.Errorf("Expected failure count 1, got %d", count)
		}

		// Trigger progression - should advance to stage 1 (3x8)
		result := triggerProgressionWithResult(t, ts, userID, stageProgID, benchID)
		if result.Data.TotalApplied != 1 {
			t.Errorf("Expected progression to apply, got TotalApplied=%d", result.Data.TotalApplied)
		}

		// Weight should not change
		if len(result.Data.Results) > 0 && result.Data.Results[0].Result != nil {
			if result.Data.Results[0].Result.Delta != 0 {
				t.Errorf("Expected delta=0 on stage advance, got %.1f", result.Data.Results[0].Result.Delta)
			}
		}
	})

	advanceUserState(t, ts, userID)

	// Test 4: Stage 1 (3x8) FAIL - advance to stage 2
	t.Run("failure at stage 1 advances to stage 2", func(t *testing.T) {
		sessionID := uuid.New().String()
		// Only 22 reps < 24 minimum (3x8)
		logSets(t, ts, userID, sessionID, benchPrescID, benchID, []setLog{
			{weight: benchMax, targetReps: 8, repsPerformed: 8, isAMRAP: false},
			{weight: benchMax, targetReps: 8, repsPerformed: 8, isAMRAP: false},
			{weight: benchMax, targetReps: 8, repsPerformed: 6, isAMRAP: false}, // Failed
		})

		count := getFailureCount(t, ts, userID, benchID, stageProgID)
		if count != 1 {
			t.Errorf("Expected failure count 1 (reset after previous trigger), got %d", count)
		}

		result := triggerProgressionWithResult(t, ts, userID, stageProgID, benchID)
		if result.Data.TotalApplied != 1 {
			t.Errorf("Expected progression to apply, got TotalApplied=%d", result.Data.TotalApplied)
		}
	})

	advanceUserState(t, ts, userID)

	// Test 5: Stage 2 (3x6) FAIL - reset to stage 0 WITHOUT deload
	t.Run("failure at stage 2 resets to stage 0 without deload", func(t *testing.T) {
		sessionID := uuid.New().String()
		// Only 16 reps < 18 minimum (3x6)
		logSets(t, ts, userID, sessionID, benchPrescID, benchID, []setLog{
			{weight: benchMax, targetReps: 6, repsPerformed: 6, isAMRAP: false},
			{weight: benchMax, targetReps: 6, repsPerformed: 6, isAMRAP: false},
			{weight: benchMax, targetReps: 6, repsPerformed: 4, isAMRAP: false}, // Failed
		})

		result := triggerProgressionWithResult(t, ts, userID, stageProgID, benchID)
		if result.Data.TotalApplied != 1 {
			t.Errorf("Expected progression to apply (reset), got TotalApplied=%d", result.Data.TotalApplied)
		}

		// Weight should NOT change because DeloadOnReset=false
		if len(result.Data.Results) > 0 && result.Data.Results[0].Result != nil {
			res := result.Data.Results[0].Result
			if res.Delta != 0 {
				t.Errorf("Expected delta=0 (no deload for T2), got %.1f", res.Delta)
			}
			if res.NewValue != benchMax {
				t.Errorf("Expected weight to remain %.1f, got %.1f", benchMax, res.NewValue)
			}
		}
	})
}

// TestGZCLPT2SuccessAfterFailure tests that success resets the failure counter.
func TestGZCLPT2SuccessAfterFailure(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user" // Uses seeded test user
	_ = testID                     // Keep for unique slugs
	benchID := "00000000-0000-0000-0000-000000000002"
	benchMax := 100.0

	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchMax)

	// Create progression
	stageProgParams := StageProgressionParams{
		Stages: []StageParam{
			{Name: "3x10", Sets: 3, Reps: 10, IsAMRAP: false, MinVolume: 30},
			{Name: "3x8", Sets: 3, Reps: 8, IsAMRAP: false, MinVolume: 24},
			{Name: "3x6", Sets: 3, Reps: 6, IsAMRAP: false, MinVolume: 18},
		},
		CurrentStage:      0,
		ResetOnExhaustion: true,
		DeloadOnReset:     false,
		MaxType:           "TRAINING_MAX",
	}

	stageProgParamsJSON, _ := json.Marshal(stageProgParams)
	stageProgBody := fmt.Sprintf(`{"name": "GZCLP T2 Success Test", "type": "STAGE_PROGRESSION", "parameters": %s}`, stageProgParamsJSON)
	stageProgResp, _ := adminPost(ts.URL("/progressions"), stageProgBody)
	var stageProgEnvelope ProgressionResponse
	json.NewDecoder(stageProgResp.Body).Decode(&stageProgEnvelope)
	stageProgResp.Body.Close()
	stageProgID := stageProgEnvelope.Data.ID

	// Create program structure (minimal)
	benchPrescID := createPrescription(t, ts, benchID, 3, 10, 100.0, 0)
	daySlug := "gzclp-t2-succ-day-" + testID
	dayBody := fmt.Sprintf(`{"name": "Day", "slug": "%s"}`, daySlug)
	dayResp, _ := adminPost(ts.URL("/days"), dayBody)
	var dayEnvelope DayResponse
	json.NewDecoder(dayResp.Body).Decode(&dayEnvelope)
	dayResp.Body.Close()
	dayID := dayEnvelope.Data.ID
	addPrescToDay(t, ts, dayID, benchPrescID)

	cycleName := "Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 1}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	weekBody := fmt.Sprintf(`{"weekNumber": 1, "cycleId": "%s"}`, cycleID)
	weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
	var weekEnvelope WeekResponse
	json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
	weekResp.Body.Close()
	weekID := weekEnvelope.Data.ID
	addDayToWeek(t, ts, weekID, dayID, "MONDAY")

	programSlug := "gzclp-succ-" + testID
	programBody := fmt.Sprintf(`{"name": "GZCLP Success", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	linkProgressionToProgram(t, ts, programID, stageProgID, benchID, 1)

	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Log a failure first (only one set fails)
	t.Run("failure increments counter", func(t *testing.T) {
		sessionID := uuid.New().String()
		logSets(t, ts, userID, sessionID, benchPrescID, benchID, []setLog{
			{weight: benchMax, targetReps: 10, repsPerformed: 10, isAMRAP: false}, // Success
			{weight: benchMax, targetReps: 10, repsPerformed: 10, isAMRAP: false}, // Success
			{weight: benchMax, targetReps: 10, repsPerformed: 8, isAMRAP: false},  // Fail - last set fails
		})

		count := getFailureCount(t, ts, userID, benchID, stageProgID)
		if count != 1 {
			t.Errorf("Expected failure count 1, got %d", count)
		}
	})

	advanceUserState(t, ts, userID)

	// Log a success - counter should reset
	t.Run("success resets counter to 0", func(t *testing.T) {
		sessionID := uuid.New().String()
		logSets(t, ts, userID, sessionID, benchPrescID, benchID, []setLog{
			{weight: benchMax, targetReps: 10, repsPerformed: 10, isAMRAP: false}, // Success
			{weight: benchMax, targetReps: 10, repsPerformed: 10, isAMRAP: false},
			{weight: benchMax, targetReps: 10, repsPerformed: 10, isAMRAP: false},
		})

		count := getFailureCount(t, ts, userID, benchID, stageProgID)
		if count != 0 {
			t.Errorf("Expected failure count 0 after success, got %d", count)
		}
	})
}

// TestGZCLPT2ManualInterventionNeeded tests the case where ResetOnExhaustion=false.
func TestGZCLPT2ManualInterventionNeeded(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user" // Uses seeded test user
	_ = testID                     // Keep for unique slugs
	benchID := "00000000-0000-0000-0000-000000000002"
	benchMax := 100.0

	createLiftMax(t, ts, userID, benchID, "TRAINING_MAX", benchMax)

	// Create progression with ResetOnExhaustion=false
	stageProgParams := StageProgressionParams{
		Stages: []StageParam{
			{Name: "3x10", Sets: 3, Reps: 10, IsAMRAP: false, MinVolume: 30},
			{Name: "3x8", Sets: 3, Reps: 8, IsAMRAP: false, MinVolume: 24},
		},
		CurrentStage:      1, // Start at last stage
		ResetOnExhaustion: false,
		DeloadOnReset:     false,
		MaxType:           "TRAINING_MAX",
	}

	stageProgParamsJSON, _ := json.Marshal(stageProgParams)
	stageProgBody := fmt.Sprintf(`{"name": "GZCLP T2 Manual", "type": "STAGE_PROGRESSION", "parameters": %s}`, stageProgParamsJSON)
	stageProgResp, _ := adminPost(ts.URL("/progressions"), stageProgBody)
	var stageProgEnvelope ProgressionResponse
	json.NewDecoder(stageProgResp.Body).Decode(&stageProgEnvelope)
	stageProgResp.Body.Close()
	stageProgID := stageProgEnvelope.Data.ID

	// Minimal program setup
	benchPrescID := createPrescription(t, ts, benchID, 3, 8, 100.0, 0)
	daySlug := "gzclp-manual-day-" + testID
	dayBody := fmt.Sprintf(`{"name": "Day", "slug": "%s"}`, daySlug)
	dayResp, _ := adminPost(ts.URL("/days"), dayBody)
	var dayEnvelope DayResponse
	json.NewDecoder(dayResp.Body).Decode(&dayEnvelope)
	dayResp.Body.Close()
	dayID := dayEnvelope.Data.ID
	addPrescToDay(t, ts, dayID, benchPrescID)

	cycleName := "Manual Cycle " + testID
	cycleBody := fmt.Sprintf(`{"name": "%s", "lengthWeeks": 1}`, cycleName)
	cycleResp, _ := adminPost(ts.URL("/cycles"), cycleBody)
	var cycleEnvelope CycleResponse
	json.NewDecoder(cycleResp.Body).Decode(&cycleEnvelope)
	cycleResp.Body.Close()
	cycleID := cycleEnvelope.Data.ID

	weekBody := fmt.Sprintf(`{"weekNumber": 1, "cycleId": "%s"}`, cycleID)
	weekResp, _ := adminPost(ts.URL("/weeks"), weekBody)
	var weekEnvelope WeekResponse
	json.NewDecoder(weekResp.Body).Decode(&weekEnvelope)
	weekResp.Body.Close()
	weekID := weekEnvelope.Data.ID
	addDayToWeek(t, ts, weekID, dayID, "MONDAY")

	programSlug := "gzclp-manual-" + testID
	programBody := fmt.Sprintf(`{"name": "GZCLP Manual", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	linkProgressionToProgram(t, ts, programID, stageProgID, benchID, 1)

	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Fail at last stage with ResetOnExhaustion=false
	t.Run("failure at last stage with no reset returns not applied", func(t *testing.T) {
		sessionID := uuid.New().String()
		logSets(t, ts, userID, sessionID, benchPrescID, benchID, []setLog{
			{weight: benchMax, targetReps: 8, repsPerformed: 6, isAMRAP: false}, // Fail
			{weight: benchMax, targetReps: 8, repsPerformed: 6, isAMRAP: false},
			{weight: benchMax, targetReps: 8, repsPerformed: 6, isAMRAP: false},
		})

		result := triggerProgressionWithResult(t, ts, userID, stageProgID, benchID)

		// Progression should NOT be applied because ResetOnExhaustion=false
		if result.Data.TotalApplied != 0 {
			t.Errorf("Expected no progression (manual intervention required), got TotalApplied=%d", result.Data.TotalApplied)
		}

		// Check for skip reason
		if len(result.Data.Results) > 0 {
			if result.Data.Results[0].Applied {
				t.Error("Expected progression to not be applied")
			}
		}
	})
}
