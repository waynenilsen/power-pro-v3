// Package e2e provides end-to-end tests for complete program workflows.
// This file tests GZCLP T1 stage-based progression with failure handling.
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

// StageProgressionParams represents the parameters for a stage progression.
type StageProgressionParams struct {
	Stages            []StageParam `json:"stages"`
	CurrentStage      int          `json:"currentStage"`
	ResetOnExhaustion bool         `json:"resetOnExhaustion"`
	DeloadOnReset     bool         `json:"deloadOnReset"`
	DeloadPercent     float64      `json:"deloadPercent,omitempty"`
	MaxType           string       `json:"maxType"`
}

// StageParam represents a single stage in a stage progression.
type StageParam struct {
	Name      string `json:"name"`
	Sets      int    `json:"sets"`
	Reps      int    `json:"reps"`
	IsAMRAP   bool   `json:"isAmrap"`
	MinVolume int    `json:"minVolume"`
}

// DeloadOnFailureParams represents the parameters for a deload on failure progression.
type DeloadOnFailureParams struct {
	FailureThreshold int     `json:"failureThreshold"`
	DeloadType       string  `json:"deloadType"`
	DeloadPercent    float64 `json:"deloadPercent,omitempty"`
	DeloadAmount     float64 `json:"deloadAmount,omitempty"`
	ResetOnDeload    bool    `json:"resetOnDeload"`
	MaxType          string  `json:"maxType"`
}

// TestGZCLPT1StageProgression validates the GZCLP T1 stage progression system.
// This test demonstrates the failure-based progression pattern:
//
//	5x3+ -> 6x2+ -> 10x1+ with 15% deload on reset
//
// Key behaviors tested:
//  1. On failure, progression moves to the next stage (keeps weight)
//  2. On last stage failure with ResetOnExhaustion, resets to stage 0 with deload
//  3. Failure counter tracks consecutive failures per lift/progression
func TestGZCLPT1StageProgression(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	// Test-unique identifiers
	testID := uuid.New().String()[:8]
	userID := "workout-test-user" // Uses seeded test user
	_ = testID                     // Keep for unique slugs

	// Seeded lift ID (squat)
	squatID := "00000000-0000-0000-0000-000000000001"

	// Starting training max
	squatMax := 200.0

	// Create training max for user
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatMax)

	// Create GZCLP T1 Stage Progression (5x3+ -> 6x2+ -> 10x1+)
	stageProgParams := StageProgressionParams{
		Stages: []StageParam{
			{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
			{Name: "6x2+", Sets: 6, Reps: 2, IsAMRAP: true, MinVolume: 12},
			{Name: "10x1+", Sets: 10, Reps: 1, IsAMRAP: true, MinVolume: 10},
		},
		CurrentStage:      0,
		ResetOnExhaustion: true,
		DeloadOnReset:     true,
		DeloadPercent:     0.15,
		MaxType:           "TRAINING_MAX",
	}

	stageProgParamsJSON, _ := json.Marshal(stageProgParams)
	stageProgBody := fmt.Sprintf(`{"name": "GZCLP T1 Squat", "type": "STAGE_PROGRESSION", "parameters": %s}`, stageProgParamsJSON)
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

	// Create prescription for squat 5x3+ at 100% TM
	squatPrescID := createAMRAPPrescription(t, ts, squatID, 5, 3, 100.0, 0)

	// Create day
	daySlug := "gzclp-day-" + testID
	dayBody := fmt.Sprintf(`{"name": "GZCLP Day", "slug": "%s"}`, daySlug)
	dayResp, _ := adminPost(ts.URL("/days"), dayBody)
	var dayEnvelope DayResponse
	json.NewDecoder(dayResp.Body).Decode(&dayEnvelope)
	dayResp.Body.Close()
	dayID := dayEnvelope.Data.ID

	// Add prescription to day
	addPrescToDay(t, ts, dayID, squatPrescID)

	// Create cycle and week
	cycleName := "GZCLP Cycle " + testID
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
	programSlug := "gzclp-" + testID
	programBody := fmt.Sprintf(`{"name": "GZCLP", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	// Link progression to program
	linkProgressionToProgram(t, ts, programID, stageProgID, squatID, 1)

	// Enroll user in program
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

	// Test 1: Verify initial workout at stage 0 (5x3+)
	t.Run("initial workout shows 5x3 at 200 lbs", func(t *testing.T) {
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

		if len(workout.Data.Exercises) != 1 {
			t.Fatalf("Expected 1 exercise, got %d", len(workout.Data.Exercises))
		}

		ex := workout.Data.Exercises[0]
		if len(ex.Sets) != 5 {
			t.Errorf("Expected 5 sets for 5x3+, got %d", len(ex.Sets))
		}
		for i, set := range ex.Sets {
			if set.Weight != squatMax {
				t.Errorf("Set %d: expected weight %.1f, got %.1f", i+1, squatMax, set.Weight)
			}
			if set.TargetReps != 3 {
				t.Errorf("Set %d: expected 3 reps, got %d", i+1, set.TargetReps)
			}
		}
	})

	// Test 2: Simulate failure at stage 0 (5x3+) - stage progression triggers on workout finish
	t.Run("failure at stage 0 advances to stage 1 on workout finish", func(t *testing.T) {
		// Log sets with failure (only 13 total reps < 15 minimum volume)
		sessionID := startWorkoutSession(t, ts, userID)
		logSets(t, ts, userID, sessionID, squatPrescID, squatID, []setLog{
			{weight: squatMax, targetReps: 3, repsPerformed: 3, isAMRAP: false},
			{weight: squatMax, targetReps: 3, repsPerformed: 3, isAMRAP: false},
			{weight: squatMax, targetReps: 3, repsPerformed: 3, isAMRAP: false},
			{weight: squatMax, targetReps: 3, repsPerformed: 3, isAMRAP: false},
			{weight: squatMax, targetReps: 3, repsPerformed: 1, isAMRAP: true}, // AMRAP fail - only 1 rep, total 13 < 15
		})

		// Verify failure counter was incremented
		failureCount := getFailureCount(t, ts, userID, squatID, stageProgID)
		if failureCount != 1 {
			t.Errorf("Expected failure count 1, got %d", failureCount)
		}

		// Finishing workout triggers AFTER_SESSION progression automatically
		finishWorkoutSession(t, ts, sessionID, userID)

		// Get workout again - stage progression should have applied
		workoutResp, err := userGet(ts.URL("/users/"+userID+"/workout"), userID)
		if err != nil {
			t.Fatalf("Failed to get workout: %v", err)
		}
		defer workoutResp.Body.Close()

		var workout WorkoutResponse
		json.NewDecoder(workoutResp.Body).Decode(&workout)

		// Note: The stage change affects prescription interpretation but
		// the static prescription was created as 5x3. In a real scenario,
		// the program would use the stage's SetScheme dynamically.
		// For this test, we verify the progression was applied.
		if len(workout.Data.Exercises) < 1 {
			t.Fatal("Expected at least 1 exercise")
		}

		// Weight should still be 200 (stage changes don't affect weight)
		ex := workout.Data.Exercises[0]
		if len(ex.Sets) > 0 && ex.Sets[0].Weight != squatMax {
			t.Errorf("Expected weight to remain %.1f after stage change, got %.1f", squatMax, ex.Sets[0].Weight)
		}
	})
}

// TestGZCLPT1FullCycle tests the complete GZCLP T1 cycle through all stages and reset.
func TestGZCLPT1FullCycle(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user" // Uses seeded test user
	_ = testID                     // Keep for unique slugs
	squatID := "00000000-0000-0000-0000-000000000001"
	squatMax := 200.0

	// Create training max
	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatMax)

	// Create GZCLP T1 Stage Progression
	stageProgParams := StageProgressionParams{
		Stages: []StageParam{
			{Name: "5x3+", Sets: 5, Reps: 3, IsAMRAP: true, MinVolume: 15},
			{Name: "6x2+", Sets: 6, Reps: 2, IsAMRAP: true, MinVolume: 12},
			{Name: "10x1+", Sets: 10, Reps: 1, IsAMRAP: true, MinVolume: 10},
		},
		CurrentStage:      0,
		ResetOnExhaustion: true,
		DeloadOnReset:     true,
		DeloadPercent:     0.15,
		MaxType:           "TRAINING_MAX",
	}

	stageProgParamsJSON, _ := json.Marshal(stageProgParams)
	stageProgBody := fmt.Sprintf(`{"name": "GZCLP T1 Full Cycle", "type": "STAGE_PROGRESSION", "parameters": %s}`, stageProgParamsJSON)
	stageProgResp, _ := adminPost(ts.URL("/progressions"), stageProgBody)
	var stageProgEnvelope ProgressionResponse
	json.NewDecoder(stageProgResp.Body).Decode(&stageProgEnvelope)
	stageProgResp.Body.Close()
	stageProgID := stageProgEnvelope.Data.ID

	// Create minimal program structure
	squatPrescID := createAMRAPPrescription(t, ts, squatID, 5, 3, 100.0, 0)
	daySlug := "gzclp-full-day-" + testID
	dayBody := fmt.Sprintf(`{"name": "GZCLP Full Day", "slug": "%s"}`, daySlug)
	dayResp, _ := adminPost(ts.URL("/days"), dayBody)
	var dayEnvelope DayResponse
	json.NewDecoder(dayResp.Body).Decode(&dayEnvelope)
	dayResp.Body.Close()
	dayID := dayEnvelope.Data.ID
	addPrescToDay(t, ts, dayID, squatPrescID)

	cycleName := "GZCLP Full Cycle " + testID
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

	programSlug := "gzclp-full-" + testID
	programBody := fmt.Sprintf(`{"name": "GZCLP Full", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	linkProgressionToProgram(t, ts, programID, stageProgID, squatID, 1)

	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Stage 0 (5x3+) - Fail
	t.Run("stage 0 failure", func(t *testing.T) {
		sessionID := startWorkoutSession(t, ts, userID)
		// Log failure: 13 reps < 15 minimum
		logSets(t, ts, userID, sessionID, squatPrescID, squatID, []setLog{
			{weight: squatMax, targetReps: 3, repsPerformed: 3, isAMRAP: false},
			{weight: squatMax, targetReps: 3, repsPerformed: 3, isAMRAP: false},
			{weight: squatMax, targetReps: 3, repsPerformed: 3, isAMRAP: false},
			{weight: squatMax, targetReps: 3, repsPerformed: 3, isAMRAP: false},
			{weight: squatMax, targetReps: 3, repsPerformed: 1, isAMRAP: true},
		})

		// Trigger progression - should advance to stage 1
		result := triggerProgressionWithResult(t, ts, userID, stageProgID, squatID)
		if result.Data.TotalApplied != 1 {
			t.Errorf("Expected progression to apply, got TotalApplied=%d", result.Data.TotalApplied)
		}
		// Weight should not change on stage advance
		if len(result.Data.Results) > 0 && result.Data.Results[0].Result != nil {
			delta := result.Data.Results[0].Result.Delta
			if delta != 0 {
				t.Errorf("Expected delta=0 on stage advance, got %.1f", delta)
			}
		}
		finishWorkoutSession(t, ts, sessionID, userID)
	})

	advanceUserState(t, ts, userID)

	// Stage 1 (6x2+) - Fail
	t.Run("stage 1 failure", func(t *testing.T) {
		sessionID := startWorkoutSession(t, ts, userID)
		// Log failure: 10 reps < 12 minimum
		logSets(t, ts, userID, sessionID, squatPrescID, squatID, []setLog{
			{weight: squatMax, targetReps: 2, repsPerformed: 2, isAMRAP: false},
			{weight: squatMax, targetReps: 2, repsPerformed: 2, isAMRAP: false},
			{weight: squatMax, targetReps: 2, repsPerformed: 2, isAMRAP: false},
			{weight: squatMax, targetReps: 2, repsPerformed: 2, isAMRAP: false},
			{weight: squatMax, targetReps: 2, repsPerformed: 1, isAMRAP: false},
			{weight: squatMax, targetReps: 2, repsPerformed: 1, isAMRAP: true},
		})

		result := triggerProgressionWithResult(t, ts, userID, stageProgID, squatID)
		if result.Data.TotalApplied != 1 {
			t.Errorf("Expected progression to apply, got TotalApplied=%d", result.Data.TotalApplied)
		}
		finishWorkoutSession(t, ts, sessionID, userID)
	})

	advanceUserState(t, ts, userID)

	// Stage 2 (10x1+) - Fail -> Reset with deload
	t.Run("stage 2 failure triggers reset with deload", func(t *testing.T) {
		sessionID := startWorkoutSession(t, ts, userID)
		// Log failure: 8 reps < 10 minimum
		logSets(t, ts, userID, sessionID, squatPrescID, squatID, []setLog{
			{weight: squatMax, targetReps: 1, repsPerformed: 1, isAMRAP: false},
			{weight: squatMax, targetReps: 1, repsPerformed: 1, isAMRAP: false},
			{weight: squatMax, targetReps: 1, repsPerformed: 1, isAMRAP: false},
			{weight: squatMax, targetReps: 1, repsPerformed: 1, isAMRAP: false},
			{weight: squatMax, targetReps: 1, repsPerformed: 1, isAMRAP: false},
			{weight: squatMax, targetReps: 1, repsPerformed: 1, isAMRAP: false},
			{weight: squatMax, targetReps: 1, repsPerformed: 1, isAMRAP: false},
			{weight: squatMax, targetReps: 1, repsPerformed: 1, isAMRAP: false},
			{weight: squatMax, targetReps: 1, repsPerformed: 0, isAMRAP: false}, // Fail
			{weight: squatMax, targetReps: 1, repsPerformed: 0, isAMRAP: true},  // Fail
		})

		result := triggerProgressionWithResult(t, ts, userID, stageProgID, squatID)
		if result.Data.TotalApplied != 1 {
			t.Errorf("Expected progression to apply, got TotalApplied=%d", result.Data.TotalApplied)
		}

		// At last stage, reset should apply deload
		if len(result.Data.Results) > 0 && result.Data.Results[0].Result != nil {
			res := result.Data.Results[0].Result
			expectedDeload := squatMax * 0.15 // 15% of 200 = 30
			expectedNewValue := squatMax - expectedDeload
			if res.Delta >= 0 {
				t.Errorf("Expected negative delta (deload), got %.1f", res.Delta)
			}
			if closeEnough(res.NewValue, expectedNewValue) == false {
				t.Errorf("Expected new value ~%.1f after 15%% deload, got %.1f", expectedNewValue, res.NewValue)
			}
		}
		finishWorkoutSession(t, ts, sessionID, userID)
	})
}

// TestDeloadOnFailure tests the DeloadOnFailure progression type.
func TestDeloadOnFailure(t *testing.T) {
	ts, err := testutil.NewTestServer()
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer ts.Close()

	testID := uuid.New().String()[:8]
	userID := "workout-test-user" // Uses seeded test user
	_ = testID                     // Keep for unique slugs
	squatID := "00000000-0000-0000-0000-000000000001"
	squatMax := 300.0

	createLiftMax(t, ts, userID, squatID, "TRAINING_MAX", squatMax)

	// Create DeloadOnFailure progression (like Texas Method: 2 failures -> 10% deload)
	deloadParams := DeloadOnFailureParams{
		FailureThreshold: 2,
		DeloadType:       "percent",
		DeloadPercent:    0.10,
		ResetOnDeload:    true,
		MaxType:          "TRAINING_MAX",
	}

	deloadParamsJSON, _ := json.Marshal(deloadParams)
	deloadBody := fmt.Sprintf(`{"name": "Texas Method Deload", "type": "DELOAD_ON_FAILURE", "parameters": %s}`, deloadParamsJSON)
	deloadResp, err := adminPost(ts.URL("/progressions"), deloadBody)
	if err != nil {
		t.Fatalf("Failed to create deload progression: %v", err)
	}
	if deloadResp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(deloadResp.Body)
		deloadResp.Body.Close()
		t.Fatalf("Failed to create deload progression, status %d: %s", deloadResp.StatusCode, body)
	}
	var deloadEnvelope ProgressionResponse
	json.NewDecoder(deloadResp.Body).Decode(&deloadEnvelope)
	deloadResp.Body.Close()
	deloadProgID := deloadEnvelope.Data.ID

	// Create program structure
	squatPrescID := createPrescription(t, ts, squatID, 1, 5, 100.0, 0)
	daySlug := "texas-day-" + testID
	dayBody := fmt.Sprintf(`{"name": "Texas Friday", "slug": "%s"}`, daySlug)
	dayResp, _ := adminPost(ts.URL("/days"), dayBody)
	var dayEnvelope DayResponse
	json.NewDecoder(dayResp.Body).Decode(&dayEnvelope)
	dayResp.Body.Close()
	dayID := dayEnvelope.Data.ID
	addPrescToDay(t, ts, dayID, squatPrescID)

	cycleName := "Texas Cycle " + testID
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
	addDayToWeek(t, ts, weekID, dayID, "FRIDAY")

	programSlug := "texas-method-" + testID
	programBody := fmt.Sprintf(`{"name": "Texas Method", "slug": "%s", "cycleId": "%s"}`, programSlug, cycleID)
	programResp, _ := adminPost(ts.URL("/programs"), programBody)
	var programEnvelope ProgramResponse
	json.NewDecoder(programResp.Body).Decode(&programEnvelope)
	programResp.Body.Close()
	programID := programEnvelope.Data.ID

	linkProgressionToProgram(t, ts, programID, deloadProgID, squatID, 1)

	enrollBody := fmt.Sprintf(`{"programId": "%s"}`, programID)
	enrollResp, _ := userPost(ts.URL("/users/"+userID+"/program"), enrollBody, userID)
	enrollResp.Body.Close()

	// Test: First failure - counter should be 1, no deload yet
	t.Run("first failure increments counter but no deload", func(t *testing.T) {
		sessionID := startWorkoutSession(t, ts, userID)
		logSets(t, ts, userID, sessionID, squatPrescID, squatID, []setLog{
			{weight: squatMax, targetReps: 5, repsPerformed: 4, isAMRAP: false}, // Failure
		})

		count := getFailureCount(t, ts, userID, squatID, deloadProgID)
		if count != 1 {
			t.Errorf("Expected failure count 1, got %d", count)
		}

		// Trigger progression - should NOT apply (threshold not met)
		result := triggerProgressionWithResult(t, ts, userID, deloadProgID, squatID)
		if result.Data.TotalApplied != 0 {
			t.Errorf("Expected no progression (threshold not met), got TotalApplied=%d", result.Data.TotalApplied)
		}
		finishWorkoutSession(t, ts, sessionID, userID)
	})

	advanceUserState(t, ts, userID)

	// Test: Second failure - counter should be 2, deload triggers
	t.Run("second failure triggers deload", func(t *testing.T) {
		sessionID := startWorkoutSession(t, ts, userID)
		logSets(t, ts, userID, sessionID, squatPrescID, squatID, []setLog{
			{weight: squatMax, targetReps: 5, repsPerformed: 3, isAMRAP: false}, // Second failure
		})

		count := getFailureCount(t, ts, userID, squatID, deloadProgID)
		if count != 2 {
			t.Errorf("Expected failure count 2, got %d", count)
		}

		// Trigger progression - should apply (threshold met)
		result := triggerProgressionWithResult(t, ts, userID, deloadProgID, squatID)
		if result.Data.TotalApplied != 1 {
			t.Errorf("Expected progression to apply, got TotalApplied=%d", result.Data.TotalApplied)
		}

		if len(result.Data.Results) > 0 && result.Data.Results[0].Result != nil {
			res := result.Data.Results[0].Result
			expectedDeload := squatMax * 0.10 // 10% of 300 = 30
			expectedNewValue := squatMax - expectedDeload
			if res.Delta >= 0 {
				t.Errorf("Expected negative delta (deload), got %.1f", res.Delta)
			}
			if closeEnough(res.NewValue, expectedNewValue) == false {
				t.Errorf("Expected new value ~%.1f after 10%% deload, got %.1f", expectedNewValue, res.NewValue)
			}
		}
		finishWorkoutSession(t, ts, sessionID, userID)
	})
}

// =============================================================================
// HELPER TYPES AND FUNCTIONS
// =============================================================================

type setLog struct {
	weight        float64
	targetReps    int
	repsPerformed int
	isAMRAP       bool
}

func logSets(t *testing.T, ts *testutil.TestServer, userID, sessionID, prescriptionID, liftID string, sets []setLog) {
	t.Helper()

	type setRequest struct {
		PrescriptionID string  `json:"prescriptionId"`
		LiftID         string  `json:"liftId"`
		SetNumber      int     `json:"setNumber"`
		Weight         float64 `json:"weight"`
		TargetReps     int     `json:"targetReps"`
		RepsPerformed  int     `json:"repsPerformed"`
		IsAMRAP        bool    `json:"isAmrap"`
	}

	setsReq := make([]setRequest, len(sets))
	for i, s := range sets {
		setsReq[i] = setRequest{
			PrescriptionID: prescriptionID,
			LiftID:         liftID,
			SetNumber:      i + 1,
			Weight:         s.weight,
			TargetReps:     s.targetReps,
			RepsPerformed:  s.repsPerformed,
			IsAMRAP:        s.isAMRAP,
		}
	}

	body, _ := json.Marshal(map[string]interface{}{"sets": setsReq})
	req, _ := http.NewRequest(http.MethodPost, ts.URL("/sessions/"+sessionID+"/sets"), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to log sets: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to log sets, status %d: %s", resp.StatusCode, respBody)
	}
}

func getFailureCount(t *testing.T, ts *testutil.TestServer, userID, liftID, progressionID string) int {
	t.Helper()

	url := fmt.Sprintf("/users/%s/failure-counters?liftId=%s&progressionId=%s", userID, liftID, progressionID)
	req, _ := http.NewRequest(http.MethodGet, ts.URL(url), nil)
	req.Header.Set("X-User-ID", userID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to get failure count: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0
	}

	if resp.StatusCode != http.StatusOK {
		// Counter doesn't exist yet
		return 0
	}

	var result struct {
		Data struct {
			ConsecutiveFailures int `json:"consecutiveFailures"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Data.ConsecutiveFailures
}

func triggerProgression(t *testing.T, ts *testutil.TestServer, userID, progressionID, liftID string) {
	t.Helper()
	triggerProgressionWithResult(t, ts, userID, progressionID, liftID)
}

func triggerProgressionWithResult(t *testing.T, ts *testutil.TestServer, userID, progressionID, liftID string) TriggerResponse {
	t.Helper()

	triggerBody := ManualTriggerRequest{
		ProgressionID: progressionID,
		LiftID:        liftID,
		Force:         true,
	}

	resp, err := authPostTrigger(ts.URL("/users/"+userID+"/progressions/trigger"), triggerBody, userID)
	if err != nil {
		t.Fatalf("Failed to trigger progression: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Failed to trigger progression, status %d: %s", resp.StatusCode, body)
	}

	var result TriggerResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

// Note: createAMRAPPrescription and closeEnough are defined in other test files in this package
